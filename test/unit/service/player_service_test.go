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

func TestPlayerService_GetPlayer_Success(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{ID: id, Name: "Bambang", Position: "CF"}, nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player, err := svc.GetPlayer("player-1")
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "Bambang", player.Name)
}

func TestPlayerService_ListPlayersByTeam(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		ListByTeamFunc: func(teamID string, page, limit int) ([]domain.Player, int64, error) {
			return []domain.Player{
				{ID: "p-1", Name: "Player 1"},
				{ID: "p-2", Name: "Player 2"},
			}, 2, nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	players, total, err := svc.ListPlayersByTeam("team-1", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, players, 2)
	assert.Equal(t, int64(2), total)
}

func TestPlayerService_UpdatePlayer_Success(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{
				ID:           id,
				TeamID:       "team-1",
				Name:         "Bambang",
				Height:       178,
				Weight:       72,
				Position:     "CF",
				JerseyNumber: 20,
				Version:      1,
			}, nil
		},
		CheckJerseyUniqueFunc: func(teamID string, jerseyNumber int, excludeID string) (bool, error) {
			return true, nil
		},
		UpdateFunc: func(player *domain.Player) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.PlayerHistory) error {
			return nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	player := &domain.Player{
		ID:           "player-1",
		TeamID:       "team-1",
		Name:         "Bambang Updated",
		Height:       178,
		Weight:       72,
		Position:     "CF",
		JerseyNumber: 20,
		Version:      1,
	}

	err := svc.UpdatePlayer(player)
	assert.NoError(t, err)
	assert.Equal(t, 2, player.Version)
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

func TestPlayerService_DeletePlayer_Success(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{ID: id, Name: "Bambang", Version: 1}, nil
		},
		CountGoalsFunc: func(playerID string) (int64, error) {
			return 0, nil
		},
		DeleteFunc: func(id string) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.PlayerHistory) error {
			return nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, mockTeamRepo)

	err := svc.DeletePlayer("player-1")
	assert.NoError(t, err)
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

func TestPlayerService_GetPlayerHistory(t *testing.T) {
	mockPlayerRepo := &mocks.MockPlayerRepository{
		GetHistoryFunc: func(playerID string) ([]domain.PlayerHistory, error) {
			return []domain.PlayerHistory{
				{PlayerID: playerID, Version: 1, Changes: `{"name":"Bambang"}`},
			}, nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, nil)

	histories, err := svc.GetPlayerHistory("player-1")
	assert.NoError(t, err)
	assert.Len(t, histories, 1)
}

func TestPlayerService_RevertPlayer_Success(t *testing.T) {
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDIncludingDeletedFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{
				ID:      id,
				Name:    "Bambang Current",
				TeamID:  "team-1",
				Version: 2,
			}, nil
		},
		GetHistoryByVersionFunc: func(playerID string, version int) (*domain.PlayerHistory, error) {
			return &domain.PlayerHistory{
				PlayerID: playerID,
				Version:  1,
				Changes:  `{"name":"Bambang V1","height":178,"weight":72,"position":"CF","jersey_number":20}`,
			}, nil
		},
		CheckJerseyUniqueFunc: func(teamID string, jerseyNumber int, excludeID string) (bool, error) {
			return true, nil
		},
		UpdateFunc: func(player *domain.Player) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.PlayerHistory) error {
			return nil
		},
	}

	svc := service.NewPlayerService(mockPlayerRepo, nil)

	err := svc.RevertPlayer("player-1", 1)
	assert.NoError(t, err)
}
