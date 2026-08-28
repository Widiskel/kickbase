package service_test

import (
	"errors"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
	"kickbase/internal/service"
	"kickbase/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestResultService_CreateResult_Success(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "scheduled",
			}, nil
		},
		UpdateFunc: func(m *domain.Match) error {
			return nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return nil, nil // No result yet
		},
		CreateFunc: func(result *domain.MatchResult) error {
			result.ID = "result-1"
			return nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			if id == "player-1" {
				return &domain.Player{ID: id, TeamID: "team-1"}, nil
			}
			return &domain.Player{ID: id, TeamID: "team-2"}, nil
		},
	}
	mockGoalRepo := &mocks.MockGoalRepository{
		CreateFunc: func(g *domain.Goal) error {
			return nil
		},
	}

	svc := service.NewResultService(mockResultRepo, mockMatchRepo, mockGoalRepo, mockPlayerRepo, nil)

	input := interfaces.CreateResultInput{
		MatchID:   "match-1",
		HomeScore: 1,
		AwayScore: 1,
		Goals: []interfaces.GoalInput{
			{PlayerID: "player-1", GoalTime: "23:00"},
			{PlayerID: "player-2", GoalTime: "67:00"},
		},
	}

	result, err := svc.CreateResult(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.HomeScore)
	assert.Equal(t, 1, result.AwayScore)
}

func TestResultService_CreateResult_MatchNotScheduled(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:     id,
				Status: "completed", // Already completed
			}, nil
		},
	}

	svc := service.NewResultService(nil, mockMatchRepo, nil, nil, nil)

	input := interfaces.CreateResultInput{
		MatchID:   "match-1",
		HomeScore: 1,
		AwayScore: 0,
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Equal(t, "match is not in scheduled status", err.Error())
}

func TestResultService_CreateResult_PlayerWrongTeam(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "scheduled",
			}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return nil, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			// Player belongs to team-3 which is not in this match
			return &domain.Player{ID: id, TeamID: "team-3"}, nil
		},
	}

	svc := service.NewResultService(mockResultRepo, mockMatchRepo, nil, mockPlayerRepo, nil)

	input := interfaces.CreateResultInput{
		MatchID:   "match-1",
		HomeScore: 1,
		AwayScore: 0,
		Goals: []interfaces.GoalInput{
			{PlayerID: "player-outsider", GoalTime: "15:00"},
		},
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to either team")
}

func TestResultService_CreateResult_MatchNotFound(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewResultService(nil, mockMatchRepo, nil, nil, nil)

	input := interfaces.CreateResultInput{
		MatchID: "non-existent",
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Equal(t, "match not found", err.Error())
}

func TestResultService_GetResult_Success(t *testing.T) {
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return &domain.MatchResult{
				ID:        "res-1",
				MatchID:   matchID,
				HomeScore: 2,
				AwayScore: 1,
			}, nil
		},
	}
	mockGoalRepo := &mocks.MockGoalRepository{
		ListByMatchResultFunc: func(resultID string) ([]domain.Goal, error) {
			return []domain.Goal{
				{PlayerID: "p-1", GoalTime: "10:00"},
			}, nil
		},
	}

	svc := service.NewResultService(mockResultRepo, nil, mockGoalRepo, nil, nil)

	res, goals, err := svc.GetResult("m-1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, goals, 1)
}
