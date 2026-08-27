package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
)

type PlayerService struct {
	playerRepo *repository.PlayerRepository
	teamRepo   *repository.TeamRepository
}

func NewPlayerService(playerRepo *repository.PlayerRepository, teamRepo *repository.TeamRepository) *PlayerService {
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

func isValidPosition(position string) bool {
	validPositions := map[string]bool{
		"CF": true, "SS": true, "LWF": true, "RWF": true,
		"AMF": true, "CMF": true, "DMF": true, "LMF": true, "RMF": true,
		"CB": true, "LB": true, "RB": true,
		"GK": true,
	}
	return validPositions[position]
}
