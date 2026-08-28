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

func TestMatchService_CreateMatch_Success(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Team " + id}, nil
		},
	}
	mockMatchRepo := &mocks.MockMatchRepository{
		CreateFunc: func(match *domain.Match) error {
			match.ID = "match-1"
			return nil
		},
		CreateHistoryFunc: func(history *domain.MatchHistory) error {
			return nil
		},
	}

	svc := service.NewMatchService(mockMatchRepo, mockTeamRepo)

	match := &domain.Match{
		MatchDate:  "2026-09-01",
		MatchTime:  "19:00:00",
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
	}

	err := svc.CreateMatch(match)
	assert.NoError(t, err)
	assert.Equal(t, "match-1", match.ID)
	assert.Equal(t, "scheduled", match.Status)
	assert.Equal(t, 1, match.Version)
}

func TestMatchService_CreateMatch_SameTeam(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{}
	mockMatchRepo := &mocks.MockMatchRepository{}

	svc := service.NewMatchService(mockMatchRepo, mockTeamRepo)

	match := &domain.Match{
		MatchDate:  "2026-09-01",
		MatchTime:  "19:00:00",
		HomeTeamID: "team-1",
		AwayTeamID: "team-1",
	}

	err := svc.CreateMatch(match)
	assert.Error(t, err)
	assert.Equal(t, "home team and away team must be different", err.Error())
}

func TestMatchService_CreateMatch_HomeTeamNotFound(t *testing.T) {
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			if id == "team-1" {
				return nil, errors.New("team not found")
			}
			return &domain.Team{ID: id}, nil
		},
	}
	mockMatchRepo := &mocks.MockMatchRepository{}

	svc := service.NewMatchService(mockMatchRepo, mockTeamRepo)

	match := &domain.Match{
		MatchDate:  "2026-09-01",
		MatchTime:  "19:00:00",
		HomeTeamID: "team-1",
		AwayTeamID: "team-2",
	}

	err := svc.CreateMatch(match)
	assert.Error(t, err)
	assert.Equal(t, "home team not found", err.Error())
}

func TestMatchService_GetMatch(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{ID: id, Status: "scheduled"}, nil
		},
	}

	svc := service.NewMatchService(mockMatchRepo, nil)

	match, err := svc.GetMatch("m-1")
	assert.NoError(t, err)
	assert.NotNil(t, match)
	assert.Equal(t, "m-1", match.ID)
}

func TestMatchService_ListMatches(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		ListFunc: func(opts interfaces.MatchFilterOptions) ([]domain.Match, int64, error) {
			return []domain.Match{
				{ID: "m-1", Status: "scheduled"},
				{ID: "m-2", Status: "completed"},
			}, 2, nil
		},
	}

	svc := service.NewMatchService(mockMatchRepo, nil)

	matches, total, err := svc.ListMatches(interfaces.MatchFilterOptions{Page: 1, Limit: 10, Status: "scheduled"})
	assert.NoError(t, err)
	assert.Len(t, matches, 2)
	assert.Equal(t, int64(2), total)
}

func TestMatchService_UpdateMatchStatus_ValidTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus string
		targetStatus  string
		wantErr       bool
	}{
		{"scheduled to completed", "scheduled", "completed", false},
		{"scheduled to cancelled", "scheduled", "cancelled", false},
		{"scheduled to deferred", "scheduled", "deferred", false},
		{"deferred to scheduled", "deferred", "scheduled", false},
		{"completed to cancelled (invalid)", "completed", "cancelled", true},
		{"cancelled to completed (invalid)", "cancelled", "completed", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTeamRepo := &mocks.MockTeamRepository{}
			mockMatchRepo := &mocks.MockMatchRepository{
				FindByIDFunc: func(id string) (*domain.Match, error) {
					return &domain.Match{ID: id, Status: tt.initialStatus, Version: 1}, nil
				},
				UpdateFunc: func(m *domain.Match) error {
					return nil
				},
				CreateHistoryFunc: func(h *domain.MatchHistory) error {
					return nil
				},
			}

			svc := service.NewMatchService(mockMatchRepo, mockTeamRepo)
			err := svc.UpdateMatchStatus("match-1", tt.targetStatus)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMatchService_GetMatchHistory(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		GetHistoryFunc: func(matchID string) ([]domain.MatchHistory, error) {
			return []domain.MatchHistory{
				{MatchID: matchID, Version: 1, Changes: `{"status":"scheduled"}`},
			}, nil
		},
	}

	svc := service.NewMatchService(mockMatchRepo, nil)

	histories, err := svc.GetMatchHistory("m-1")
	assert.NoError(t, err)
	assert.Len(t, histories, 1)
}

func TestMatchService_RevertMatch(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{ID: id, Status: "cancelled", Version: 2}, nil
		},
		GetHistoryByVersionFunc: func(matchID string, version int) (*domain.MatchHistory, error) {
			return &domain.MatchHistory{
				MatchID: matchID,
				Version: 1,
				Changes: `{"status":"scheduled","match_date":"2026-09-01","match_time":"19:00:00"}`,
			}, nil
		},
		UpdateFunc: func(m *domain.Match) error {
			return nil
		},
		CreateHistoryFunc: func(h *domain.MatchHistory) error {
			return nil
		},
	}

	svc := service.NewMatchService(mockMatchRepo, nil)

	err := svc.RevertMatch("m-1", 1)
	assert.NoError(t, err)
}
