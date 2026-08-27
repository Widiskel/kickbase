package service_test

import (
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupPlayerService(t *testing.T) (*service.PlayerService, *gorm.DB) {
	db := setupTestDB(t)
	playerRepo := repository.NewPlayerRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	svc := service.NewPlayerService(playerRepo, teamRepo)
	return svc, db
}

func createTestTeam(t *testing.T, db *gorm.DB) *domain.Team {
	teamRepo := repository.NewTeamRepository(db)
	team := &domain.Team{
		Name:        "Persija Jakarta",
		City:        "Jakarta",
		FoundedYear: 1928,
		Address:     "Jl. Pintu Satu Senayan",
	}
	teamRepo.Create(team)
	return team
}

func TestPlayerService_CreatePlayer_Success(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		Playstyle:    stringPtr("goal_poacher"),
		JerseyNumber: 20,
	}

	err := svc.CreatePlayer(player)
	assert.NoError(t, err)
	assert.NotEmpty(t, player.ID)
	assert.Equal(t, 1, player.Version)
}

func TestPlayerService_CreatePlayer_TeamNotFound(t *testing.T) {
	svc, _ := setupPlayerService(t)

	player := &domain.Player{
		TeamID:       "00000000-0000-0000-0000-000000000000",
		Name:         "Player 1",
		Height:       175,
		Weight:       70,
		Position:     "CF",
		JerseyNumber: 10,
	}

	err := svc.CreatePlayer(player)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "team not found")
}

func TestPlayerService_CreatePlayer_DuplicateJersey(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player1 := &domain.Player{
		TeamID:       team.ID,
		Name:         "Player 1",
		Height:       175,
		Weight:       70,
		Position:     "CF",
		JerseyNumber: 10,
	}
	svc.CreatePlayer(player1)

	player2 := &domain.Player{
		TeamID:       team.ID,
		Name:         "Player 2",
		Height:       180,
		Weight:       75,
		Position:     "CMF",
		JerseyNumber: 10, // Same jersey number
	}

	err := svc.CreatePlayer(player2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jersey number already exists")
}

func TestPlayerService_CreatePlayer_SameJerseyDifferentTeam(t *testing.T) {
	svc, db := setupPlayerService(t)
	team1 := createTestTeam(t, db)

	teamRepo := repository.NewTeamRepository(db)
	team2 := &domain.Team{
		Name:        "Persib Bandung",
		City:        "Bandung",
		FoundedYear: 1933,
		Address:     "Jl. ABC",
	}
	teamRepo.Create(team2)

	player1 := &domain.Player{
		TeamID:       team1.ID,
		Name:         "Player 1",
		Height:       175,
		Weight:       70,
		Position:     "CF",
		JerseyNumber: 10,
	}
	svc.CreatePlayer(player1)

	player2 := &domain.Player{
		TeamID:       team2.ID,
		Name:         "Player 2",
		Height:       180,
		Weight:       75,
		Position:     "CMF",
		JerseyNumber: 10, // Same jersey, different team
	}

	err := svc.CreatePlayer(player2)
	assert.NoError(t, err) // Should succeed
}

func TestPlayerService_CreatePlayer_InvalidPosition(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Player 1",
		Height:       175,
		Weight:       70,
		Position:     "INVALID", // Invalid position
		JerseyNumber: 10,
	}

	err := svc.CreatePlayer(player)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid position")
}

func TestPlayerService_GetPlayer_Success(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	found, err := svc.GetPlayer(player.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Bambang Pamungkas", found.Name)
}

func TestPlayerService_GetPlayer_NotFound(t *testing.T) {
	svc, _ := setupPlayerService(t)

	_, err := svc.GetPlayer("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPlayerService_ListPlayersByTeam(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	// Create 3 players
	for i := 1; i <= 3; i++ {
		svc.CreatePlayer(&domain.Player{
			TeamID:       team.ID,
			Name:         "Player " + string(rune('0'+i)),
			Height:       175,
			Weight:       70,
			Position:     "CF",
			JerseyNumber: i,
		})
	}

	players, total, err := svc.ListPlayersByTeam(team.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, players, 3)
}

func TestPlayerService_UpdatePlayer_Success(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	player.Position = "SS"
	player.Name = "Bambang P. Updated"
	err := svc.UpdatePlayer(player)
	assert.NoError(t, err)

	updated, _ := svc.GetPlayer(player.ID)
	assert.Equal(t, "SS", updated.Position)
	assert.Equal(t, 2, updated.Version)
}

func TestPlayerService_UpdatePlayer_VersionMismatch(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	player.Version = 999
	err := svc.UpdatePlayer(player)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version mismatch")
}

func TestPlayerService_DeletePlayer_Success(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	err := svc.DeletePlayer(player.ID)
	assert.NoError(t, err)

	// Verify deleted
	players, total, _ := svc.ListPlayersByTeam(team.ID, 1, 10)
	assert.Equal(t, int64(0), total)
	assert.Len(t, players, 0)
}

func TestPlayerService_DeletePlayer_WithGoals(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	// Create a match and result first
	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team.ID,
		AwayTeamID: team.ID,
		Status:     "completed",
	}
	db.Create(match)

	result := &domain.MatchResult{
		MatchID:   match.ID,
		HomeScore: 3,
		AwayScore: 1,
	}
	db.Create(result)

	// Add a goal record
	goal := &domain.Goal{
		MatchResultID: result.ID,
		PlayerID:      player.ID,
		GoalTime:      "23'",
	}
	db.Create(goal)

	err := svc.DeletePlayer(player.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "goal records")
}

func TestPlayerService_DeletePlayer_NotFound(t *testing.T) {
	svc, _ := setupPlayerService(t)

	err := svc.DeletePlayer("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPlayerService_GetPlayerHistory(t *testing.T) {
	svc, db := setupPlayerService(t)
	team := createTestTeam(t, db)

	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Bambang Pamungkas",
		Height:       178.5,
		Weight:       72.0,
		Position:     "CF",
		JerseyNumber: 20,
	}
	svc.CreatePlayer(player)

	// Update to create more history
	player.Position = "SS"
	svc.UpdatePlayer(player)

	history, err := svc.GetPlayerHistory(player.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2)
}

func stringPtr(s string) *string {
	return &s
}
