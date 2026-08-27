package service_test

import (
	"strings"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupResultService(t *testing.T) (*service.ResultService, *gorm.DB) {
	db := setupTestDB(t)
	resultRepo := repository.NewResultRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	svc := service.NewResultService(resultRepo, matchRepo, goalRepo, playerRepo, db)
	return svc, db
}

func createTestMatchForResult(t *testing.T, db *gorm.DB) (*domain.Team, *domain.Team, *domain.Match) {
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)

	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		Status:     "scheduled",
		Version:    1,
	}
	matchRepo.Create(match)

	return team1, team2, match
}

func TestResultService_CreateResult_Success(t *testing.T) {
	svc, db := setupResultService(t)
	team1, team2, match := createTestMatchForResult(t, db)

	// Create players
	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Player 1", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)
	player2 := &domain.Player{TeamID: team2.ID, Name: "Player 2", Height: 180, Weight: 75, Position: "CMF", JerseyNumber: 8, Version: 1}
	playerRepo.Create(player2)

	input := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
		Goals: []service.GoalInput{
			{PlayerID: player1.ID, GoalTime: "23'"},
			{PlayerID: player1.ID, GoalTime: "45'"},
			{PlayerID: player2.ID, GoalTime: "67'"},
			{PlayerID: player1.ID, GoalTime: "89'"},
		},
	}

	result, err := svc.CreateResult(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.HomeScore)
	assert.Equal(t, 1, result.AwayScore)

	// Verify match status updated
	matchRepo := repository.NewMatchRepository(db)
	updatedMatch, _ := matchRepo.FindByID(match.ID)
	assert.Equal(t, "completed", updatedMatch.Status)
}

func TestResultService_CreateResult_MatchNotFound(t *testing.T) {
	svc, _ := setupResultService(t)

	input := service.CreateResultInput{
		MatchID:   "00000000-0000-0000-0000-000000000000",
		HomeScore: 3,
		AwayScore: 1,
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "match not found")
}

func TestResultService_CreateResult_MatchNotScheduled(t *testing.T) {
	svc, db := setupResultService(t)
	_, _, match := createTestMatchForResult(t, db)

	// Update match status to completed
	matchRepo := repository.NewMatchRepository(db)
	match.Status = "completed"
	matchRepo.Update(match)

	input := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in scheduled status")
}

func TestResultService_CreateResult_DuplicateResult(t *testing.T) {
	svc, db := setupResultService(t)
	team1, _, match := createTestMatchForResult(t, db)

	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Player 1", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)

	// Create first result
	input := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
		Goals:     []service.GoalInput{{PlayerID: player1.ID, GoalTime: "23'"}},
	}
	svc.CreateResult(input)

	// Try to create duplicate - should fail because match is now completed
	input2 := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 2,
		AwayScore: 0,
	}

	_, err := svc.CreateResult(input2)
	assert.Error(t, err)
	// Either "already exists" or "not in scheduled status" is acceptable
	assert.True(t,
		strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "not in scheduled status"),
		"Expected error about duplicate or status, got: %s", err.Error(),
	)
}

func TestResultService_CreateResult_PlayerFromWrongTeam(t *testing.T) {
	svc, db := setupResultService(t)
	team1, _, match := createTestMatchForResult(t, db)

	// Create a player from team1
	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Player 1", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)

	// Create a third team and player
	teamRepo := repository.NewTeamRepository(db)
	team3 := &domain.Team{Name: "Arema", City: "Malang", FoundedYear: 1987, Address: "Addr"}
	teamRepo.Create(team3)
	player3 := &domain.Player{TeamID: team3.ID, Name: "Player 3", Height: 170, Weight: 65, Position: "LB", JerseyNumber: 3, Version: 1}
	playerRepo.Create(player3)

	input := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
		Goals:     []service.GoalInput{{PlayerID: player3.ID, GoalTime: "23'"}}, // Wrong team
	}

	_, err := svc.CreateResult(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to either team")
}

func TestResultService_GetResult_Success(t *testing.T) {
	svc, db := setupResultService(t)
	team1, _, match := createTestMatchForResult(t, db)

	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Player 1", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)

	input := service.CreateResultInput{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
		Goals:     []service.GoalInput{{PlayerID: player1.ID, GoalTime: "23'"}},
	}
	svc.CreateResult(input)

	result, goals, err := svc.GetResult(match.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, goals, 1)
	assert.Equal(t, "23'", goals[0].GoalTime)
}

func TestResultService_GetResult_NotFound(t *testing.T) {
	svc, _ := setupResultService(t)

	_, _, err := svc.GetResult("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
