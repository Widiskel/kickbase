package service_test

import (
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMatchService(t *testing.T) (*service.MatchService, *gorm.DB) {
	db := setupTestDB(t)
	matchRepo := repository.NewMatchRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	svc := service.NewMatchService(matchRepo, teamRepo)
	return svc, db
}

func createTestTeamForMatch(t *testing.T, db *gorm.DB, name string) *domain.Team {
	teamRepo := repository.NewTeamRepository(db)
	team := &domain.Team{
		Name:        name,
		City:        "City",
		FoundedYear: 2000,
		Address:     "Addr",
	}
	teamRepo.Create(team)
	return team
}

func TestMatchService_CreateMatch_Success(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
	}

	err := svc.CreateMatch(match)
	assert.NoError(t, err)
	assert.NotEmpty(t, match.ID)
	assert.Equal(t, "scheduled", match.Status)
	assert.Equal(t, 1, match.Version)
}

func TestMatchService_CreateMatch_SameTeam(t *testing.T) {
	svc, db := setupMatchService(t)
	team := createTestTeamForMatch(t, db, "Persija")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team.ID,
		AwayTeamID: team.ID, // Same team
	}

	err := svc.CreateMatch(match)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be different")
}

func TestMatchService_CreateMatch_HomeTeamNotFound(t *testing.T) {
	svc, db := setupMatchService(t)
	team := createTestTeamForMatch(t, db, "Persija")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: "00000000-0000-0000-0000-000000000000",
		AwayTeamID: team.ID,
	}

	err := svc.CreateMatch(match)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "home team not found")
}

func TestMatchService_CreateMatch_AwayTeamNotFound(t *testing.T) {
	svc, db := setupMatchService(t)
	team := createTestTeamForMatch(t, db, "Persija")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team.ID,
		AwayTeamID: "00000000-0000-0000-0000-000000000000",
	}

	err := svc.CreateMatch(match)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "away team not found")
}

func TestMatchService_GetMatch_Success(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
	}
	svc.CreateMatch(match)

	found, err := svc.GetMatch(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, "2026-09-15", found.MatchDate)
}

func TestMatchService_GetMatch_NotFound(t *testing.T) {
	svc, _ := setupMatchService(t)

	_, err := svc.GetMatch("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMatchService_ListMatches(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	// Create 3 matches
	for i := 0; i < 3; i++ {
		svc.CreateMatch(&domain.Match{
			MatchDate:  "2026-09-15",
			MatchTime:  "19:30:00",
			HomeTeamID: team1.ID,
			AwayTeamID: team2.ID,
		})
	}

	matches, total, err := svc.ListMatches(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, matches, 3)
}

func TestMatchService_UpdateMatchStatus_Success(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
	}
	svc.CreateMatch(match)

	err := svc.UpdateMatchStatus(match.ID, "cancelled")
	assert.NoError(t, err)

	updated, _ := svc.GetMatch(match.ID)
	assert.Equal(t, "cancelled", updated.Status)
	assert.Equal(t, 2, updated.Version)
}

func TestMatchService_UpdateMatchStatus_InvalidTransition(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
	}
	svc.CreateMatch(match)

	// Try to go from scheduled to completed (should be valid)
	err := svc.UpdateMatchStatus(match.ID, "completed")
	assert.NoError(t, err)

	// Try to go from completed to cancelled (should be invalid)
	err = svc.UpdateMatchStatus(match.ID, "cancelled")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
}

func TestMatchService_UpdateMatchStatus_NotFound(t *testing.T) {
	svc, _ := setupMatchService(t)

	err := svc.UpdateMatchStatus("00000000-0000-0000-0000-000000000000", "cancelled")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMatchService_GetMatchHistory(t *testing.T) {
	svc, db := setupMatchService(t)
	team1 := createTestTeamForMatch(t, db, "Persija")
	team2 := createTestTeamForMatch(t, db, "Persib")

	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
	}
	svc.CreateMatch(match)

	// Update status to create more history
	svc.UpdateMatchStatus(match.ID, "cancelled")

	history, err := svc.GetMatchHistory(match.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2) // At least create + update
}
