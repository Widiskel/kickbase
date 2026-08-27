package service

import (
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/repository"

	"gorm.io/gorm"
)

type ReportService struct {
	db         *gorm.DB
	matchRepo  *repository.MatchRepository
	resultRepo *repository.ResultRepository
	goalRepo   *repository.GoalRepository
	teamRepo   *repository.TeamRepository
	playerRepo *repository.PlayerRepository
}

func NewReportService(
	db *gorm.DB,
	matchRepo *repository.MatchRepository,
	resultRepo *repository.ResultRepository,
	goalRepo *repository.GoalRepository,
	teamRepo *repository.TeamRepository,
	playerRepo *repository.PlayerRepository,
) *ReportService {
	return &ReportService{
		db:         db,
		matchRepo:  matchRepo,
		resultRepo: resultRepo,
		goalRepo:   goalRepo,
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
	}
}

type MatchReport struct {
	MatchID          string      `json:"match_id"`
	MatchDate        string      `json:"match_date"`
	MatchTime        string      `json:"match_time"`
	HomeTeam         string      `json:"home_team"`
	AwayTeam         string      `json:"away_team"`
	HomeScore        *int        `json:"home_score,omitempty"`
	AwayScore        *int        `json:"away_score,omitempty"`
	Status           string      `json:"status"`
	TopScorers       []TopScorer `json:"top_scorers,omitempty"`
	CumulativeHomeWins int       `json:"cumulative_home_wins"`
	CumulativeAwayWins int       `json:"cumulative_away_wins"`
}

type TopScorer struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Goals      int    `json:"goals"`
}

func (s *ReportService) GetMatchReport(matchID string) (*MatchReport, error) {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, errors.New("match not found")
	}

	homeTeam, _ := s.teamRepo.FindByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.FindByID(match.AwayTeamID)

	report := &MatchReport{
		MatchID:   match.ID,
		MatchDate: match.MatchDate,
		MatchTime: match.MatchTime,
		HomeTeam:  homeTeam.Name,
		AwayTeam:  awayTeam.Name,
		Status:    match.Status,
	}

	// Get result if exists
	result, err := s.resultRepo.FindByMatchID(matchID)
	if err == nil {
		report.HomeScore = &result.HomeScore
		report.AwayScore = &result.AwayScore

		// Determine winner
		if result.HomeScore > result.AwayScore {
			report.Status = "Tim Home Menang"
		} else if result.AwayScore > result.HomeScore {
			report.Status = "Tim Away Menang"
		} else {
			report.Status = "Draw"
		}

		// Get top scorers
		goals, _ := s.goalRepo.ListByMatchResult(result.ID)
		scorerMap := make(map[string]int)
		for _, goal := range goals {
			scorerMap[goal.PlayerID]++
		}

		maxGoals := 0
		for _, count := range scorerMap {
			if count > maxGoals {
				maxGoals = count
			}
		}

		for playerID, count := range scorerMap {
			if count == maxGoals {
				player, _ := s.playerRepo.FindByID(playerID)
				report.TopScorers = append(report.TopScorers, TopScorer{
					PlayerID:   playerID,
					PlayerName: player.Name,
					Goals:      count,
				})
			}
		}
	}

	// Calculate cumulative wins
	report.CumulativeHomeWins = s.countTeamWins(match.HomeTeamID, match.MatchDate)
	report.CumulativeAwayWins = s.countTeamWins(match.AwayTeamID, match.MatchDate)

	return report, nil
}

func (s *ReportService) ListMatchReports(page, limit int) ([]MatchReport, int64, error) {
	matches, total, err := s.matchRepo.List(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list matches: %w", err)
	}

	var reports []MatchReport
	for _, match := range matches {
		report, err := s.GetMatchReport(match.ID)
		if err != nil {
			continue
		}
		reports = append(reports, *report)
	}

	return reports, total, nil
}

func (s *ReportService) countTeamWins(teamID string, beforeDate string) int {
	var count int64

	// Count as home team wins
	s.db.Model(&domain.Match{}).
		Joins("JOIN match_results ON match_results.match_id = matches.id").
		Where("matches.home_team_id = ? AND match_results.home_score > match_results.away_score AND matches.match_date <= ?", teamID, beforeDate).
		Count(&count)

	// Count as away team wins
	var awayCount int64
	s.db.Model(&domain.Match{}).
		Joins("JOIN match_results ON match_results.match_id = matches.id").
		Where("matches.away_team_id = ? AND match_results.away_score > match_results.home_score AND matches.match_date <= ?", teamID, beforeDate).
		Count(&awayCount)

	return int(count + awayCount)
}
