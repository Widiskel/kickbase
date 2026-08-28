package service_test

import (
	"errors"
	"testing"

	"kickbase/internal/domain"
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
