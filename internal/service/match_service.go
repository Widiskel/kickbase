package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

type MatchService struct {
	matchRepo interfaces.MatchRepository
	teamRepo  interfaces.TeamRepository
}

func NewMatchService(matchRepo interfaces.MatchRepository, teamRepo interfaces.TeamRepository) *MatchService {
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

func (s *MatchService) RevertMatch(matchID string, targetVersion int) error {
	// Get current match
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return errors.New("match not found")
	}

	// Get target version from history
	history, err := s.matchRepo.GetHistoryByVersion(matchID, targetVersion)
	if err != nil {
		return errors.New("target version not found in history")
	}

	// Parse the changes from history
	var changes map[string]interface{}
	if err := json.Unmarshal([]byte(history.Changes), &changes); err != nil {
		return fmt.Errorf("failed to parse history changes: %w", err)
	}

	// Apply changes to match
	if v, ok := changes["match_date"]; ok {
		match.MatchDate = v.(string)
	}
	if v, ok := changes["match_time"]; ok {
		match.MatchTime = v.(string)
	}
	if v, ok := changes["home_team_id"]; ok {
		match.HomeTeamID = v.(string)
	}
	if v, ok := changes["away_team_id"]; ok {
		match.AwayTeamID = v.(string)
	}
	if v, ok := changes["status"]; ok {
		match.Status = v.(string)
	}

	match.Version = match.Version + 1
	if err := s.matchRepo.Update(match); err != nil {
		return fmt.Errorf("failed to revert match: %w", err)
	}

	// Record revert in history
	revertChanges := map[string]interface{}{
		"reverted_to_version": targetVersion,
	}
	revertJSON, _ := json.Marshal(revertChanges)
	revertHistory := &domain.MatchHistory{
		MatchID: match.ID,
		Version: match.Version,
		Changes: string(revertJSON),
	}
	s.matchRepo.CreateHistory(revertHistory)

	return nil
}

func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"scheduled": {"completed", "cancelled", "deferred"},
		"deferred":  {"scheduled", "cancelled"},
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
