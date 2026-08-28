package service_test

import (
	"errors"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/service"
	"kickbase/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestPlayerService_CreatePlayer_Success(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		CheckJerseyUniqueFunc: func(teamID string, jerseyNumber int, excludeID string) (bool, error) {
			return true, nil
		},
		CreateFunc: func(player *domain.Player) error {
			player.ID = "player-1"
			return nil
		},
		CreateHistoryFunc: func(history *domain.PlayerHistory) error {
			return nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		TeamID:       "team-1",
		Name:         "Bambang Pamungkas",
		Height:       178,
		Weight:       72,
		Position:     "CF",
		JerseyNumber: 20,
	}

	err := svc.CreatePlayer(player)
	assert.NoError(t, err)
	assert.Equal(t, "player-1", player.ID)
	assert.Equal(t, 1, player.Version)
}

func TestPlayerService_CreatePlayer_TeamNotFound(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return nil, errors.New("team not found")
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		TeamID:       "non-existent-team",
		Name:         "Bambang Pamungkas",
		Height:       178,
		Weight:       72,
		Position:     "CF",
		JerseyNumber: 20,
	}

	err := svc.CreatePlayer(player)
	assert.Error(t, err)
	assert.Equal(t, "team not found", err.Error())
}

func TestPlayerService_CreatePlayer_DuplicateJersey(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		CheckJerseyUniqueFunc: func(teamID string, jerseyNumber int, excludeID string) (bool, error) {
			return false, nil // Not unique
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		TeamID:       "team-1",
		Name:         "Bambang Pamungkas",
		Height:       178,
		Weight:       72,
		Position:     "CF",
		JerseyNumber: 20,
	}

	err := svc.CreatePlayer(player)
	assert.Error(t, err)
	assert.Equal(t, "jersey number already exists in this team", err.Error())
}

func TestPlayerService_CreatePlayer_InvalidPosition(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		TeamID:       "team-1",
		Name:         "Bambang Pamungkas",
		Height:       178,
		Weight:       72,
		Position:     "INVALID_POS",
		JerseyNumber: 20,
	}

	err := svc.CreatePlayer(player)
	assert.Error(t, err)
	assert.Equal(t, "invalid position", err.Error())
}

func TestPlayerService_DeletePlayer_WithGoals(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{ID: id, Name: "Bambang", Version: 1}, nil
		},
		CountGoalsFunc: func(playerID string) (int64, error) {
			return 5, nil // Has scored 5 goals
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	err := svc.DeletePlayer("player-1")
	assert.Error(t, err)
	assert.Equal(t, "cannot delete player with goal records", err.Error())
}

func TestPlayerService_UpdatePlayer_VersionMismatch(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{ID: id, Name: "Bambang", Version: 2}, nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		ID:           "player-1",
		Name:         "Bambang Updated",
		Height:       178,
		Weight:       72,
		Position:     "CF",
		JerseyNumber: 20,
		Version:      1, // Expected 1 but DB is 2
	}

	err := svc.UpdatePlayer(player)
	assert.Error(t, err)
	assert.Equal(t, "version mismatch", err.Error())
}
