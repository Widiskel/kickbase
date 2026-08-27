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
	assert.Contains(t, err.Error(), "already exists")
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	expectedTeam := &domain.Team{
		ID:   "test-id",
		Name: "Persija Jakarta",
		City: "Jakarta",
	}

	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return expectedTeam, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team, err := svc.GetTeam("test-id")
	assert.NoError(t, err)
	assert.Equal(t, expectedTeam, team)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewTeamService(mockRepo)

	_, err := svc.GetTeam("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTeamService_DeleteTeam_Success(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
		CountPlayersFunc: func(teamID string) (int64, error) {
			return 0, nil // No players
		},
		CreateHistoryFunc: func(history *domain.TeamHistory) error {
			return nil
		},
		DeleteFunc: func(id string) error {
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
			return &domain.Team{ID: id, Name: "Persija"}, nil
		},
		CountPlayersFunc: func(teamID string) (int64, error) {
			return 5, nil // Has players
		},
	}

	svc := service.NewTeamService(mockRepo)

	err := svc.DeleteTeam("test-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "active players")
}

func TestTeamService_UpdateTeam_VersionMismatch(t *testing.T) {
	mockRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Version: 5}, nil
		},
	}

	svc := service.NewTeamService(mockRepo)

	team := &domain.Team{
		ID:      "test-id",
		Version: 3, // Wrong version
		Name:    "Updated",
	}

	err := svc.UpdateTeam(team)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version mismatch")
}
