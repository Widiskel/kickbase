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

func TestReportService_GetMatchReport_Scheduled(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				MatchDate:  "2026-09-01",
				MatchTime:  "19:00:00",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "scheduled",
			}, nil
		},
	}
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			if id == "team-1" {
				return &domain.Team{ID: id, Name: "Persija"}, nil
			}
			return &domain.Team{ID: id, Name: "Persib"}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return nil, nil // Not completed yet
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, mockResultRepo, nil, mockTeamRepo, nil)

	report, err := svc.GetMatchReport("match-1")
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Persija", report.HomeTeam)
	assert.Equal(t, "Persib", report.AwayTeam)
	assert.Equal(t, "scheduled", report.Status)
	assert.Nil(t, report.HomeScore)
	assert.Nil(t, report.AwayScore)
}

func TestReportService_GetMatchReport_CompletedHomeWin(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				MatchDate:  "2026-09-01",
				MatchTime:  "19:00:00",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "completed",
			}, nil
		},
	}
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			if id == "team-1" {
				return &domain.Team{ID: id, Name: "Persija"}, nil
			}
			return &domain.Team{ID: id, Name: "Persib"}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return &domain.MatchResult{
				ID:        "result-1",
				MatchID:   matchID,
				HomeScore: 2,
				AwayScore: 0,
			}, nil
		},
	}
	mockGoalRepo := &mocks.MockGoalRepository{
		ListByMatchResultFunc: func(resultID string) ([]domain.Goal, error) {
			return []domain.Goal{
				{PlayerID: "player-1", GoalTime: "20:00"},
				{PlayerID: "player-1", GoalTime: "80:00"},
			}, nil
		},
	}
	mockPlayerRepo := &mocks.MockPlayerRepository{
		FindByIDFunc: func(id string) (*domain.Player, error) {
			return &domain.Player{ID: id, Name: "Bambang Pamungkas"}, nil
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, mockResultRepo, mockGoalRepo, mockTeamRepo, mockPlayerRepo)

	report, err := svc.GetMatchReport("match-1")
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Tim Home Menang", report.Status)
	assert.Equal(t, 2, *report.HomeScore)
	assert.Equal(t, 0, *report.AwayScore)
	assert.Len(t, report.TopScorers, 1)
	assert.Equal(t, "Bambang Pamungkas", report.TopScorers[0].PlayerName)
	assert.Equal(t, 2, report.TopScorers[0].Goals)
}

func TestReportService_GetMatchReport_CompletedAwayWin(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				MatchDate:  "2026-09-01",
				MatchTime:  "19:00:00",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "completed",
			}, nil
		},
	}
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			if id == "team-1" {
				return &domain.Team{ID: id, Name: "Persija"}, nil
			}
			return &domain.Team{ID: id, Name: "Persib"}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return &domain.MatchResult{
				ID:        "result-1",
				MatchID:   matchID,
				HomeScore: 0,
				AwayScore: 1,
			}, nil
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, mockResultRepo, nil, mockTeamRepo, nil)

	report, err := svc.GetMatchReport("match-1")
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Tim Away Menang", report.Status)
}

func TestReportService_GetMatchReport_CompletedDraw(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{
				ID:         id,
				MatchDate:  "2026-09-01",
				MatchTime:  "19:00:00",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     "completed",
			}, nil
		},
	}
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			if id == "team-1" {
				return &domain.Team{ID: id, Name: "Persija"}, nil
			}
			return &domain.Team{ID: id, Name: "Persib"}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return &domain.MatchResult{
				ID:        "result-1",
				MatchID:   matchID,
				HomeScore: 2,
				AwayScore: 2,
			}, nil
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, mockResultRepo, nil, mockTeamRepo, nil)

	report, err := svc.GetMatchReport("match-1")
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Draw", report.Status)
}

func TestReportService_GetMatchReport_NotFound(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return nil, errors.New("not found")
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, nil, nil, nil, nil)

	_, err := svc.GetMatchReport("non-existent")
	assert.Error(t, err)
	assert.Equal(t, "match not found", err.Error())
}

func TestReportService_ListMatchReports(t *testing.T) {
	mockMatchRepo := &mocks.MockMatchRepository{
		ListFunc: func(opts interfaces.MatchFilterOptions) ([]domain.Match, int64, error) {
			return []domain.Match{
				{ID: "m-1", MatchDate: "2026-09-01", HomeTeamID: "t-1", AwayTeamID: "t-2", Status: "scheduled"},
			}, 1, nil
		},
		FindByIDFunc: func(id string) (*domain.Match, error) {
			return &domain.Match{ID: id, MatchDate: "2026-09-01", HomeTeamID: "t-1", AwayTeamID: "t-2", Status: "scheduled"}, nil
		},
	}
	mockTeamRepo := &mocks.MockTeamRepository{
		FindByIDFunc: func(id string) (*domain.Team, error) {
			return &domain.Team{ID: id, Name: "Team"}, nil
		},
	}
	mockResultRepo := &mocks.MockResultRepository{
		FindByMatchIDFunc: func(matchID string) (*domain.MatchResult, error) {
			return nil, nil
		},
	}

	svc := service.NewReportService(nil, mockMatchRepo, mockResultRepo, nil, mockTeamRepo, nil)

	reports, total, err := svc.ListMatchReports(interfaces.ReportFilterOptions{Page: 1, Limit: 10, Status: "scheduled"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, reports, 1)
}
