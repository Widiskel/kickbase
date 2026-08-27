package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
)

type MatchService struct {
	matchRepo *repository.MatchRepository
	teamRepo  *repository.TeamRepository
}

func NewMatchService(matchRepo *repository.MatchRepository, teamRepo *repository.TeamRepository) *MatchService {
	return &MatchService{
		matchRepo: matchRepo,
		teamRepo:  teamRepo,
	}
}

func (s *MatchService) CreateMatch(match *domain.Match) error {
	// Validate home != away
	if match.HomeTeamID == match.AwayTeamID {
		return errors.New("home team and away team must be different")
	}

	// Verify both teams exist
	_, err := s.teamRepo.FindByID(match.HomeTeamID)
	if err != nil {
		return errors.New("home team not found")
	}

	_, err = s.teamRepo.FindByID(match.AwayTeamID)
	if err != nil {
		return errors.New("away team not found")
	}

	match.Status = "scheduled"
	match.Version = 1
	if err := s.matchRepo.Create(match); err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	// Create history record
	changes := map[string]interface{}{
		"match_date":    match.MatchDate,
		"match_time":    match.MatchTime,
		"home_team_id":  match.HomeTeamID,
		"away_team_id":  match.AwayTeamID,
		"status":        match.Status,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.MatchHistory{
		MatchID: match.ID,
		Version: 1,
		Changes: string(changesJSON),
	}
	s.matchRepo.CreateHistory(history)

	return nil
}

func (s *MatchService) GetMatch(id string) (*domain.Match, error) {
	return s.matchRepo.FindByID(id)
}

func (s *MatchService) ListMatches(page, limit int) ([]domain.Match, int64, error) {
	return s.matchRepo.List(page, limit)
}

func (s *MatchService) UpdateMatchStatus(id string, status string) error {
	match, err := s.matchRepo.FindByID(id)
	if err != nil {
		return errors.New("match not found")
	}

	// Validate status transition
	if !isValidStatusTransition(match.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", match.Status, status)
	}

	oldStatus := match.Status
	match.Status = status
	match.Version = match.Version + 1
	if err := s.matchRepo.Update(match); err != nil {
		return fmt.Errorf("failed to update match status: %w", err)
	}

	// Create history record
	changes := map[string]interface{}{
		"status": map[string]interface{}{"old": oldStatus, "new": status},
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.MatchHistory{
		MatchID: match.ID,
		Version: match.Version,
		Changes: string(changesJSON),
	}
	s.matchRepo.CreateHistory(history)

	return nil
}

func (s *MatchService) GetMatchHistory(matchID string) ([]domain.MatchHistory, error) {
	return s.matchRepo.GetHistory(matchID)
}

func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"scheduled": {"completed", "cancelled"},
		"completed": {},
		"cancelled": {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
