package service_test

import (
	"errors"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/service"
	"kickbase/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTeamService_CreateTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByNameFunc: func(name string) (*domain.Team, error) {
			return nil, nil // No duplicate
		},
		CreateFunc: func(team *domain.Team) error {
			team.ID = "test-id"
			return nil
		},
		CreateHistoryFunc: func(history *domain.TeamHistory) error {
			return nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team := &domain.Team{
		Name:        "Persija Jakarta",
		City:        "Jakarta",
		FoundedYear: 1928,
		Address:     "Jl. Pintu Satu Senayan",
	}

	err := svc.CreateTeam(team)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", team.ID)
	assert.Equal(t, 1, team.Version)
}

func TestTeamService_CreateTeam_DuplicateName(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByNameFunc: func(name string) (*domain.Team, error) {
			return &domain.Team{Name: name}, nil // Duplicate found
		},
	}

	svc := service.NewTeamService(mockRepo)

	team := &domain.Team{
		Name:        "Persija Jakarta",
		City:        "Jakarta",
		FoundedYear: 1928,
		Address:     "Jl. Pintu Satu Senayan",
	}

	err := svc.CreateTeam(team)
	assert.Error(t, err)
	assert.Equal(t, "team name already exists", err.Error())
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:          id,
				Name:        "Persija Jakarta",
				City:        "Jakarta",
				FoundedYear: 1928,
			}, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team, err := svc.GetTeam("test-id")
	assert.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, "Persija Jakarta", team.Name)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewTeamService(mockRepo)

	team, err := svc.GetTeam("non-existent")
	assert.Error(t, err)
	assert.Nil(t, team)
}

func TestTeamService_ListTeams(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		ListFunc: func(page, limit int) ([]domain.Team, int64, error) {
			return []domain.Team{
				{ID: "team-1", Name: "Persija"},
				{ID: "team-2", Name: "Persib"},
			}, 2, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	teams, total, err := svc.ListTeams(1, 10)
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, int64(2), total)
}

func TestTeamService_UpdateTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:          id,
				Name:        "Persija",
				City:        "Jakarta",
				FoundedYear: 1928,
				Version:     1,
			}, nil
		},
		UpdateFunc: func(team *domain.Team) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.TeamHistory) error {
			return nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team := &domain.Team{
		ID:          "team-1",
		Name:        "Persija Updated",
		City:        "Jakarta Pusat",
		FoundedYear: 1928,
		Version:     1,
	}

	err := svc.UpdateTeam(team)
	assert.NoError(t, err)
	assert.Equal(t, 2, team.Version)
}

func TestTeamService_UpdateTeam_VersionMismatch(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:      id,
				Version: 2,
			}, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team := &domain.Team{
		ID:      "test-id",
		Version: 1,
	}

	err := svc.UpdateTeam(team)
	assert.Error(t, err)
	assert.Equal(t, "version mismatch", err.Error())
}

func TestTeamService_DeleteTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:      id,
				Name:    "Persija",
				Version: 1,
			}, nil
		},
		CountPlayersFunc: func(teamID string) (int64, error) {
			return 0, nil // No players
		},
		DeleteFunc: func(id string) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.TeamHistory) error {
			return nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	err := svc.DeleteTeam("test-id")
	assert.NoError(t, err)
}

func TestTeamService_DeleteTeam_WithPlayers(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:      id,
				Version: 1,
			}, nil
		},
		CountPlayersFunc: func(teamID string) (int64, error) {
			return 5, nil // Has 5 players
		},
	}

	svc := service.NewTeamService(mockRepo)

	err := svc.DeleteTeam("test-id")
	assert.Error(t, err)
	assert.Equal(t, "cannot delete team with active players", err.Error())
}

func TestTeamService_GetTeamHistory(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		GetHistoryFunc: func(teamID string) ([]domain.TeamHistory, error) {
			return []domain.TeamHistory{
				{TeamID: teamID, Version: 1, Changes: `{"name":"Persija"}`},
			}, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	histories, err := svc.GetTeamHistory("team-1")
	assert.NoError(t, err)
	assert.Len(t, histories, 1)
}

func TestTeamService_RevertTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDIncludingDeletedFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{
				ID:      id,
				Name:    "Persija Current",
				Version: 2,
			}, nil
		},
		GetHistoryByVersionFunc: func(teamID string, version int) (*domain.TeamHistory, error) {
			return &domain.TeamHistory{
				TeamID:  teamID,
				Version: 1,
				Changes: `{"name":"Persija V1","city":"Jakarta","founded_year":1928,"address":"Senayan"}`,
			}, nil
		},
		UpdateFunc: func(team *domain.Team) error {
			return nil
		},
		CreateHistoryFunc: func(history *domain.TeamHistory) error {
			return nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	err := svc.RevertTeam("team-1", 1)
	assert.NoError(t, err)
}
