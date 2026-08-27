package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
)

type TeamService struct {
	teamRepo *repository.TeamRepository
}

func NewTeamService(teamRepo *repository.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(team *domain.Team) error {
	// Check duplicate name
	existing, _ := s.teamRepo.FindByName(team.Name)
	if existing != nil {
		return errors.New("team name already exists")
	}

	team.Version = 1
	if err := s.teamRepo.Create(team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// Create history record
	changes := map[string]interface{}{
		"name":         team.Name,
		"logo_url":     team.LogoURL,
		"founded_year": team.FoundedYear,
		"address":      team.Address,
		"city":         team.City,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.TeamHistory{
		TeamID:  team.ID,
		Version: 1,
		Changes: string(changesJSON),
	}
	s.teamRepo.CreateHistory(history)

	return nil
}

func (s *TeamService) GetTeam(id string) (*domain.Team, error) {
	return s.teamRepo.FindByID(id)
}

func (s *TeamService) ListTeams(page, limit int) ([]domain.Team, int64, error) {
	return s.teamRepo.List(page, limit)
}

func (s *TeamService) UpdateTeam(team *domain.Team) error {
	existing, err := s.teamRepo.FindByID(team.ID)
	if err != nil {
		return errors.New("team not found")
	}

	if team.Version != existing.Version {
		return errors.New("version mismatch")
	}

	// Track changes for history
	changes := map[string]interface{}{}
	if team.Name != existing.Name {
		changes["name"] = map[string]interface{}{"old": existing.Name, "new": team.Name}
	}
	if team.LogoURL != existing.LogoURL {
		changes["logo_url"] = map[string]interface{}{"old": existing.LogoURL, "new": team.LogoURL}
	}
	if team.FoundedYear != existing.FoundedYear {
		changes["founded_year"] = map[string]interface{}{"old": existing.FoundedYear, "new": team.FoundedYear}
	}
	if team.Address != existing.Address {
		changes["address"] = map[string]interface{}{"old": existing.Address, "new": team.Address}
	}
	if team.City != existing.City {
		changes["city"] = map[string]interface{}{"old": existing.City, "new": team.City}
	}

	team.Version = existing.Version + 1
	if err := s.teamRepo.Update(team); err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	// Create history record
	changesJSON, _ := json.Marshal(changes)
	history := &domain.TeamHistory{
		TeamID:  team.ID,
		Version: team.Version,
		Changes: string(changesJSON),
	}
	s.teamRepo.CreateHistory(history)

	return nil
}

func (s *TeamService) DeleteTeam(id string) error {
	team, err := s.teamRepo.FindByID(id)
	if err != nil {
		return errors.New("team not found")
	}

	// Check if team has active players
	playerCount, err := s.teamRepo.CountPlayers(id)
	if err != nil {
		return fmt.Errorf("failed to check players: %w", err)
	}
	if playerCount > 0 {
		return errors.New("cannot delete team with active players")
	}

	// Create history record before deletion
	changes := map[string]interface{}{
		"deleted": true,
	}
	changesJSON, _ := json.Marshal(changes)
	history := &domain.TeamHistory{
		TeamID:  team.ID,
		Version: team.Version + 1,
		Changes: string(changesJSON),
	}
	s.teamRepo.CreateHistory(history)

	return s.teamRepo.Delete(id)
}

func (s *TeamService) GetTeamHistory(teamID string) ([]domain.TeamHistory, error) {
	return s.teamRepo.GetHistory(teamID)
}
