package service_test

import (
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=kickbase sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available, skipping test")
	}

	// Clean tables
	db.Exec("DROP TABLE IF EXISTS goals, match_results, matches, players, teams, team_histories, player_histories, match_histories, match_result_histories, goal_histories CASCADE")

	err = db.AutoMigrate(
		&domain.Team{},
		&domain.TeamHistory{},
		&domain.Player{},
		&domain.PlayerHistory{},
		&domain.Match{},
		&domain.MatchHistory{},
		&domain.MatchResult{},
		&domain.MatchResultHistory{},
		&domain.Goal{},
		&domain.GoalHistory{},
	)
	require.NoError(t, err)

	return db
}

func setupTeamService(t *testing.T) (*service.TeamService, *gorm.DB) {
	db := setupTestDB(t)
	teamRepo := repository.NewTeamRepository(db)
	svc := service.NewTeamService(teamRepo)
	return svc, db
}

func TestTeamService_CreateTeam_Success(t *testing.T) {
	svc, _ := setupTeamService(t)

	team := &domain.Team{
		Name:        "Persija Jakarta",
		LogoURL:     "https://example.com/logo.png",
		FoundedYear: 1928,
		Address:     "Jl. Pintu Satu Senayan",
		City:        "Jakarta",
	}

	err := svc.CreateTeam(team)
	assert.NoError(t, err)
	assert.NotEmpty(t, team.ID)
	assert.Equal(t, 1, team.Version)
}

func TestTeamService_CreateTeam_DuplicateName(t *testing.T) {
	svc, _ := setupTeamService(t)

	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	err := svc.CreateTeam(team1)
	require.NoError(t, err)

	team2 := &domain.Team{Name: "Persija", City: "Bandung", FoundedYear: 1933, Address: "Addr2"}
	err = svc.CreateTeam(team2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	svc, _ := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	found, err := svc.GetTeam(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Persija", found.Name)
	assert.Equal(t, "Jakarta", found.City)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	svc, _ := setupTeamService(t)

	_, err := svc.GetTeam("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTeamService_ListTeams_Pagination(t *testing.T) {
	svc, _ := setupTeamService(t)

	// Create 5 teams
	for i := 0; i < 5; i++ {
		svc.CreateTeam(&domain.Team{
			Name:        "Team " + string(rune('A'+i)),
			City:        "City",
			FoundedYear: 2000,
			Address:     "Addr",
		})
	}

	// Page 1, limit 2
	teams, total, err := svc.ListTeams(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, teams, 2)

	// Page 2, limit 2
	teams, total, err = svc.ListTeams(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, teams, 2)

	// Page 3, limit 2
	teams, total, err = svc.ListTeams(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, teams, 1)
}

func TestTeamService_UpdateTeam_Success(t *testing.T) {
	svc, _ := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	team.City = "Jakarta Pusat"
	err := svc.UpdateTeam(team)
	assert.NoError(t, err)

	updated, _ := svc.GetTeam(team.ID)
	assert.Equal(t, "Jakarta Pusat", updated.City)
	assert.Equal(t, 2, updated.Version)
}

func TestTeamService_UpdateTeam_VersionMismatch(t *testing.T) {
	svc, _ := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	team.Version = 999
	err := svc.UpdateTeam(team)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version mismatch")
}

func TestTeamService_DeleteTeam_Success(t *testing.T) {
	svc, db := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	err := svc.DeleteTeam(team.ID)
	assert.NoError(t, err)

	// Should not appear in list
	teams, total, _ := svc.ListTeams(1, 10)
	assert.Equal(t, int64(0), total)
	assert.Len(t, teams, 0)

	// Verify history was recorded
	var historyCount int64
	db.Model(&domain.TeamHistory{}).Where("team_id = ?", team.ID).Count(&historyCount)
	assert.Greater(t, historyCount, int64(0))
}

func TestTeamService_DeleteTeam_WithPlayers(t *testing.T) {
	svc, db := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	// Add a player
	player := &domain.Player{
		TeamID:       team.ID,
		Name:         "Player 1",
		Height:       175,
		Weight:       70,
		Position:     "CF",
		JerseyNumber: 10,
	}
	db.Create(player)

	err := svc.DeleteTeam(team.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "active players")
}

func TestTeamService_DeleteTeam_NotFound(t *testing.T) {
	svc, _ := setupTeamService(t)

	err := svc.DeleteTeam("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTeamService_GetTeamHistory(t *testing.T) {
	svc, _ := setupTeamService(t)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	// Update to create more history
	team.City = "Jakarta Pusat"
	svc.UpdateTeam(team)

	history, err := svc.GetTeamHistory(team.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2) // At least create + update
}
