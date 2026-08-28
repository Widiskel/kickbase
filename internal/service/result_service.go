package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"

	"gorm.io/gorm"
)

type ResultService struct {
	resultRepo interfaces.ResultRepository
	matchRepo  interfaces.MatchRepository
	goalRepo   interfaces.GoalRepository
	playerRepo interfaces.PlayerRepository
	db         *gorm.DB
}

func NewResultService(
	resultRepo interfaces.ResultRepository,
	matchRepo interfaces.MatchRepository,
	goalRepo interfaces.GoalRepository,
	playerRepo interfaces.PlayerRepository,
	db *gorm.DB,
) *ResultService {
	return &ResultService{
		resultRepo: resultRepo,
		matchRepo:  matchRepo,
		goalRepo:   goalRepo,
		playerRepo: playerRepo,
		db:         db,
	}
}

func (s *ResultService) CreateResult(input interfaces.CreateResultInput) (*domain.MatchResult, error) {
	// Verify match exists and is scheduled
	match, err := s.matchRepo.FindByID(input.MatchID)
	if err != nil {
		return nil, errors.New("match not found")
	}

	if match.Status != "scheduled" {
		return nil, errors.New("match is not in scheduled status")
	}

	if input.HomeScore < 0 || input.AwayScore < 0 {
		return nil, errors.New("score cannot be negative")
	}

	totalExpectedGoals := input.HomeScore + input.AwayScore
	if len(input.Goals) != totalExpectedGoals {
		return nil, fmt.Errorf("number of goal entries (%d) must match total score (%d)", len(input.Goals), totalExpectedGoals)
	}

	// Check no existing result
	existingResult, _ := s.resultRepo.FindByMatchID(input.MatchID)
	if existingResult != nil {
		return nil, errors.New("match result already exists")
	}

	// Validate goal scorers belong to participating teams
	for _, goal := range input.Goals {
		player, err := s.playerRepo.FindByID(goal.PlayerID)
		if err != nil {
			return nil, fmt.Errorf("player %s not found", goal.PlayerID)
		}
		if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
			return nil, fmt.Errorf("player %s does not belong to either team", goal.PlayerID)
		}
	}

	// Result model
	result := &domain.MatchResult{
		MatchID:   input.MatchID,
		HomeScore: input.HomeScore,
		AwayScore: input.AwayScore,
		Version:   1,
	}

	// If DB is nil (e.g. unit tests with mocks)
	if s.db == nil {
		if err := s.resultRepo.Create(result); err != nil {
			return nil, fmt.Errorf("failed to create result: %w", err)
		}
		for _, goalInput := range input.Goals {
			goal := &domain.Goal{
				MatchResultID: result.ID,
				PlayerID:      goalInput.PlayerID,
				GoalTime:      goalInput.GoalTime,
			}
			if s.goalRepo != nil {
				_ = s.goalRepo.Create(goal)
			}
		}
		match.Status = "completed"
		_ = s.matchRepo.Update(match)
		return result, nil
	}

	// Use transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(result).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create result: %w", err)
	}

	// Create goals
	for _, goalInput := range input.Goals {
		goal := &domain.Goal{
			MatchResultID: result.ID,
			PlayerID:      goalInput.PlayerID,
			GoalTime:      goalInput.GoalTime,
		}
		if err := tx.Create(goal).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create goal: %w", err)
		}
	}

	// Update match status to completed
	if err := tx.Model(&domain.Match{}).Where("id = ?", input.MatchID).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update match status: %w", err)
	}

	// Create history record for result
	changes := map[string]interface{}{
		"home_score": result.HomeScore,
		"away_score": result.AwayScore,
		"match_id":   result.MatchID,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.MatchResultHistory{
		MatchResultID: result.ID,
		Version:       1,
		Changes:       string(changesJSON),
	}
	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	// Create history record for match status change
	matchChanges := map[string]interface{}{
		"status": map[string]interface{}{"old": "scheduled", "new": "completed"},
	}
	matchChangesJSON, _ := json.Marshal(matchChanges)
	matchHistory := &domain.MatchHistory{
		MatchID: match.ID,
		Version: match.Version + 1,
		Changes: string(matchChangesJSON),
	}
	if err := tx.Create(matchHistory).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create match history: %w", err)
	}

	tx.Commit()
	return result, nil
}

func (s *ResultService) GetResult(matchID string) (*domain.MatchResult, []domain.Goal, error) {
	result, err := s.resultRepo.FindByMatchID(matchID)
	if err != nil {
		return nil, nil, errors.New("match result not found")
	}

	goals, err := s.goalRepo.ListByMatchResult(result.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get goals: %w", err)
	}

	return result, goals, nil
}
