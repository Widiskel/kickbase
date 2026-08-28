package service

import (
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"

	"gorm.io/gorm"
)

type ReportService struct {
	db         *gorm.DB
	matchRepo  interfaces.MatchRepository
	resultRepo interfaces.ResultRepository
	goalRepo   interfaces.GoalRepository
	teamRepo   interfaces.TeamRepository
	playerRepo interfaces.PlayerRepository
}

func NewReportService(
	db *gorm.DB,
	matchRepo interfaces.MatchRepository,
	resultRepo interfaces.ResultRepository,
	goalRepo interfaces.GoalRepository,
	teamRepo interfaces.TeamRepository,
	playerRepo interfaces.PlayerRepository,
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

func (s *ReportService) GetMatchReport(matchID string) (*interfaces.MatchReport, error) {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, errors.New("match not found")
	}

	homeTeam, _ := s.teamRepo.FindByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.FindByID(match.AwayTeamID)

	report := &interfaces.MatchReport{
		MatchID:    match.ID,
		MatchDate:  match.MatchDate,
		MatchTime:  match.MatchTime,
		HomeTeam:   homeTeam.Name,
		AwayTeam:   awayTeam.Name,
		Status:     match.Status,
		TopScorers: []interfaces.TopScorer{},
	}

	// Get result if exists
	result, err := s.resultRepo.FindByMatchID(matchID)
	if err == nil && result != nil {
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
		if s.goalRepo != nil {
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
					playerName := ""
					if s.playerRepo != nil {
						if player, _ := s.playerRepo.FindByID(playerID); player != nil {
							playerName = player.Name
						}
					}
					report.TopScorers = append(report.TopScorers, interfaces.TopScorer{
						PlayerID:   playerID,
						PlayerName: playerName,
						Goals:      count,
					})
				}
			}
		}
	}

	// Calculate cumulative wins
	report.CumulativeHomeWins = s.countTeamWins(match.HomeTeamID, match.MatchDate)
	report.CumulativeAwayWins = s.countTeamWins(match.AwayTeamID, match.MatchDate)

	return report, nil
}

func (s *ReportService) ListMatchReports(opts interfaces.ReportFilterOptions) ([]interfaces.MatchReport, int64, error) {
	matchOpts := interfaces.MatchFilterOptions{
		TeamID: opts.TeamID,
		SortBy: opts.SortBy,
		Order:  opts.Order,
		Page:   opts.Page,
		Limit:  opts.Limit,
	}

	matches, total, err := s.matchRepo.List(matchOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list matches: %w", err)
	}

	reports := make([]interfaces.MatchReport, 0, len(matches))
	for _, match := range matches {
		report, err := s.GetMatchReport(match.ID)
		if err != nil {
			continue
		}
		if opts.Status != "" && report.Status != opts.Status {
			continue
		}
		reports = append(reports, *report)
	}

	return reports, total, nil
}

func (s *ReportService) countTeamWins(teamID string, beforeDate string) int {
	if s.db == nil {
		return 0
	}

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
