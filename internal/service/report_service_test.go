package service_test

import (
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupReportService(t *testing.T) (*service.ReportService, *gorm.DB) {
	db := setupTestDB(t)
	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewResultRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	svc := service.NewReportService(db, matchRepo, resultRepo, goalRepo, teamRepo, playerRepo)
	return svc, db
}

func TestReportService_GetMatchReport_Scheduled(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	matchRepo := repository.NewMatchRepository(db)
	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		Status:     "scheduled",
	}
	matchRepo.Create(match)

	report, err := svc.GetMatchReport(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Persija", report.HomeTeam)
	assert.Equal(t, "Persib", report.AwayTeam)
	assert.Equal(t, "scheduled", report.Status)
	assert.Nil(t, report.HomeScore)
	assert.Nil(t, report.AwayScore)
}

func TestReportService_GetMatchReport_HomeWin(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Bambang", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)

	matchRepo := repository.NewMatchRepository(db)
	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		Status:     "completed",
	}
	matchRepo.Create(match)

	resultRepo := repository.NewResultRepository(db)
	result := &domain.MatchResult{MatchID: match.ID, HomeScore: 3, AwayScore: 1, Version: 1}
	resultRepo.Create(result)

	goalRepo := repository.NewGoalRepository(db)
	goalRepo.Create(&domain.Goal{MatchResultID: result.ID, PlayerID: player1.ID, GoalTime: "23'"})
	goalRepo.Create(&domain.Goal{MatchResultID: result.ID, PlayerID: player1.ID, GoalTime: "45'"})
	goalRepo.Create(&domain.Goal{MatchResultID: result.ID, PlayerID: player1.ID, GoalTime: "67'"})

	report, err := svc.GetMatchReport(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Tim Home Menang", report.Status)
	assert.Equal(t, 3, *report.HomeScore)
	assert.Equal(t, 1, *report.AwayScore)
	assert.Len(t, report.TopScorers, 1)
	assert.Equal(t, "Bambang", report.TopScorers[0].PlayerName)
	assert.Equal(t, 3, report.TopScorers[0].Goals)
}

func TestReportService_GetMatchReport_Draw(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	matchRepo := repository.NewMatchRepository(db)
	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		Status:     "completed",
	}
	matchRepo.Create(match)

	resultRepo := repository.NewResultRepository(db)
	result := &domain.MatchResult{MatchID: match.ID, HomeScore: 2, AwayScore: 2, Version: 1}
	resultRepo.Create(result)

	report, err := svc.GetMatchReport(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Draw", report.Status)
}

func TestReportService_GetMatchReport_AwayWin(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	matchRepo := repository.NewMatchRepository(db)
	match := &domain.Match{
		MatchDate:  "2026-09-15",
		MatchTime:  "19:30:00",
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		Status:     "completed",
	}
	matchRepo.Create(match)

	resultRepo := repository.NewResultRepository(db)
	result := &domain.MatchResult{MatchID: match.ID, HomeScore: 1, AwayScore: 3, Version: 1}
	resultRepo.Create(result)

	report, err := svc.GetMatchReport(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Tim Away Menang", report.Status)
}

func TestReportService_GetMatchReport_NotFound(t *testing.T) {
	svc, _ := setupReportService(t)

	_, err := svc.GetMatchReport("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestReportService_GetMatchReport_CumulativeWins(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewResultRepository(db)

	// Match 1: Persija wins
	match1 := &domain.Match{MatchDate: "2026-09-01", MatchTime: "19:00", HomeTeamID: team1.ID, AwayTeamID: team2.ID, Status: "completed"}
	matchRepo.Create(match1)
	resultRepo.Create(&domain.MatchResult{MatchID: match1.ID, HomeScore: 2, AwayScore: 0, Version: 1})

	// Match 2: Persib wins
	match2 := &domain.Match{MatchDate: "2026-09-08", MatchTime: "19:00", HomeTeamID: team2.ID, AwayTeamID: team1.ID, Status: "completed"}
	matchRepo.Create(match2)
	resultRepo.Create(&domain.MatchResult{MatchID: match2.ID, HomeScore: 3, AwayScore: 1, Version: 1})

	// Match 3: Persija wins again
	match3 := &domain.Match{MatchDate: "2026-09-15", MatchTime: "19:00", HomeTeamID: team1.ID, AwayTeamID: team2.ID, Status: "completed"}
	matchRepo.Create(match3)
	resultRepo.Create(&domain.MatchResult{MatchID: match3.ID, HomeScore: 1, AwayScore: 0, Version: 1})

	report, err := svc.GetMatchReport(match3.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, report.CumulativeHomeWins)  // Persija won 2 times
	assert.Equal(t, 1, report.CumulativeAwayWins)   // Persib won 1 time
}

func TestReportService_ListMatchReports(t *testing.T) {
	svc, db := setupReportService(t)

	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	matchRepo := repository.NewMatchRepository(db)
	for i := 0; i < 3; i++ {
		matchRepo.Create(&domain.Match{
			MatchDate:  "2026-09-15",
			MatchTime:  "19:30:00",
			HomeTeamID: team1.ID,
			AwayTeamID: team2.ID,
			Status:     "scheduled",
		})
	}

	reports, total, err := svc.ListMatchReports(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, reports, 3)
}
