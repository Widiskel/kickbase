package service_test

import (
	"os"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=kickbase sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available, skipping test")
	}

	// Clean tables
	db.Exec("DROP TABLE IF EXISTS goals, match_results, matches, players, teams, goals_histories, match_result_histories, match_histories, player_histories, team_histories CASCADE")

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
	assert.NoError(t, err)

	return db
}

func TestTeamService_CreateTeam(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

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
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	err := svc.CreateTeam(team1)
	assert.NoError(t, err)

	team2 := &domain.Team{Name: "Persija", City: "Bandung", FoundedYear: 1933, Address: "Addr2"}
	err = svc.CreateTeam(team2)
	assert.Error(t, err)
}

func TestTeamService_GetTeam(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	found, err := svc.GetTeam(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Persija", found.Name)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	_, err := svc.GetTeam("non-existent-id")
	assert.Error(t, err)
}

func TestTeamService_ListTeams(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	svc.CreateTeam(&domain.Team{Name: "Team A", City: "City A", FoundedYear: 2000, Address: "Addr"})
	svc.CreateTeam(&domain.Team{Name: "Team B", City: "City B", FoundedYear: 2001, Address: "Addr"})

	teams, total, err := svc.ListTeams(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, teams, 2)
}

func TestTeamService_UpdateTeam(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

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
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	team.Version = 999
	err := svc.UpdateTeam(team)
	assert.Error(t, err)
}

func TestTeamService_DeleteTeam(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

	err := svc.DeleteTeam(team.ID)
	assert.NoError(t, err)

	teams, total, _ := svc.ListTeams(1, 10)
	assert.Equal(t, int64(0), total)
	assert.Len(t, teams, 0)
}

func TestTeamService_DeleteTeam_WithPlayers(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewTeamService(db)

	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	svc.CreateTeam(team)

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
}
