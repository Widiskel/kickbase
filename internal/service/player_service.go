package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

type PlayerService struct {
	playerRepo interfaces.PlayerRepository
	teamRepo   interfaces.TeamRepository
}

func NewPlayerService(playerRepo interfaces.PlayerRepository, teamRepo interfaces.TeamRepository) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

func (s *PlayerService) CreatePlayer(player *domain.Player) error {
	// Verify team exists
	_, err := s.teamRepo.FindByID(player.TeamID)
	if err != nil {
		return errors.New("team not found")
	}

	// Check jersey uniqueness
	unique, err := s.playerRepo.CheckJerseyUnique(player.TeamID, player.JerseyNumber, "")
	if err != nil {
		return fmt.Errorf("failed to check jersey uniqueness: %w", err)
	}
	if !unique {
		return errors.New("jersey number already exists in this team")
	}

	// Validate position
	if !isValidPosition(player.Position) {
		return errors.New("invalid position")
	}

	player.Version = 1
	if err := s.playerRepo.Create(player); err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	// Create history record
	changes := map[string]interface{}{
		"name":          player.Name,
		"height":        player.Height,
		"weight":        player.Weight,
		"position":      player.Position,
		"playstyle":     player.Playstyle,
		"jersey_number": player.JerseyNumber,
		"team_id":       player.TeamID,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.PlayerHistory{
		PlayerID: player.ID,
		Version:  1,
		Changes:  string(changesJSON),
	}
	s.playerRepo.CreateHistory(history)

	return nil
}

func (s *PlayerService) GetPlayer(id string) (*domain.Player, error) {
	return s.playerRepo.FindByID(id)
}

func (s *PlayerService) ListPlayersByTeam(teamID string, page, limit int) ([]domain.Player, int64, error) {
	return s.playerRepo.ListByTeam(teamID, page, limit)
}

func (s *PlayerService) UpdatePlayer(player *domain.Player) error {
	existing, err := s.playerRepo.FindByID(player.ID)
	if err != nil {
		return errors.New("player not found")
	}

	if player.Version != existing.Version {
		return errors.New("version mismatch")
	}

	// Check jersey uniqueness if changed
	if player.JerseyNumber != existing.JerseyNumber {
		unique, err := s.playerRepo.CheckJerseyUnique(player.TeamID, player.JerseyNumber, player.ID)
		if err != nil {
			return fmt.Errorf("failed to check jersey uniqueness: %w", err)
		}
		if !unique {
			return errors.New("jersey number already exists in this team")
		}
	}

	// Validate position
	if !isValidPosition(player.Position) {
		return errors.New("invalid position")
	}

	player.Version = existing.Version + 1
	if err := s.playerRepo.Update(player); err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	// Create history record
	changes := map[string]interface{}{
		"name":          map[string]interface{}{"old": existing.Name, "new": player.Name},
		"height":        map[string]interface{}{"old": existing.Height, "new": player.Height},
		"weight":        map[string]interface{}{"old": existing.Weight, "new": player.Weight},
		"position":      map[string]interface{}{"old": existing.Position, "new": player.Position},
		"playstyle":     map[string]interface{}{"old": existing.Playstyle, "new": player.Playstyle},
		"jersey_number": map[string]interface{}{"old": existing.JerseyNumber, "new": player.JerseyNumber},
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.PlayerHistory{
		PlayerID: player.ID,
		Version:  player.Version,
		Changes:  string(changesJSON),
	}
	s.playerRepo.CreateHistory(history)

	return nil
}

func (s *PlayerService) DeletePlayer(id string) error {
	player, err := s.playerRepo.FindByID(id)
	if err != nil {
		return errors.New("player not found")
	}

	// Check if player has goals
	goalCount, err := s.playerRepo.CountGoals(id)
	if err != nil {
		return fmt.Errorf("failed to check goals: %w", err)
	}
	if goalCount > 0 {
		return errors.New("cannot delete player with goal records")
	}

	// Create history record before deletion
	changes := map[string]interface{}{
		"deleted": true,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.PlayerHistory{
		PlayerID: player.ID,
		Version:  player.Version + 1,
		Changes:  string(changesJSON),
	}
	s.playerRepo.CreateHistory(history)

	return s.playerRepo.Delete(id)
}

func (s *PlayerService) GetPlayerHistory(playerID string) ([]domain.PlayerHistory, error) {
	return s.playerRepo.GetHistory(playerID)
}

func (s *PlayerService) RevertPlayer(playerID string, targetVersion int) error {
	// Get current player
	player, err := s.playerRepo.FindByID(playerID)
	if err != nil {
		// Try to find in soft-deleted records
		player, err = s.playerRepo.FindByIDIncludingDeleted(playerID)
		if err != nil {
			return errors.New("player not found")
		}
	}

	// Get target version from history
	history, err := s.playerRepo.GetHistoryByVersion(playerID, targetVersion)
	if err != nil {
		return errors.New("target version not found in history")
	}

	// Parse the changes from history
	var changes map[string]interface{}
	if err := json.Unmarshal([]byte(history.Changes), &changes); err != nil {
		return fmt.Errorf("failed to parse history changes: %w", err)
	}

	// Apply changes to player
	if v, ok := changes["name"]; ok {
		player.Name = v.(string)
	}
	if v, ok := changes["height"]; ok {
		player.Height = v.(float64)
	}
	if v, ok := changes["weight"]; ok {
		player.Weight = v.(float64)
	}
	if v, ok := changes["position"]; ok {
		player.Position = v.(string)
	}
	if v, ok := changes["playstyle"]; ok {
		if v == nil {
			player.Playstyle = nil
		} else {
			s := v.(string)
			player.Playstyle = &s
		}
	}
	if v, ok := changes["jersey_number"]; ok {
		player.JerseyNumber = int(v.(float64))
	}

	// Restore if deleted
	if player.DeletedAt.Valid {
		player.DeletedAt.Time = time.Time{}
		player.DeletedAt.Valid = false
	}

	player.Version = player.Version + 1
	if err := s.playerRepo.Update(player); err != nil {
		return fmt.Errorf("failed to revert player: %w", err)
	}

	// Record revert in history
	revertChanges := map[string]interface{}{
		"reverted_to_version": targetVersion,
	}
	revertJSON, _ := json.Marshal(revertChanges)
	revertHistory := &domain.PlayerHistory{
		PlayerID: player.ID,
		Version:  player.Version,
		Changes:  string(revertJSON),
	}
	s.playerRepo.CreateHistory(revertHistory)

	return nil
}

func isValidPosition(position string) bool {
	validPositions := map[string]bool{
		"CF": true, "SS": true, "LWF": true, "RWF": true,
		"AMF": true, "CMF": true, "DMF": true, "LMF": true, "RMF": true,
		"CB": true, "LB": true, "RB": true,
		"GK": true,
	}
	return validPositions[position]
}
