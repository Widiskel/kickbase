package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) *gorm.DB {
	dsn := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=kickbase sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available, skipping integration test")
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

func setupIntegrationRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db := setupIntegrationDB(t)
	r := router.Setup(db)
	return r, db
}

func TestIntegration_TeamCRUD(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	// Create team
	body := `{"name":"Persija Jakarta","logo_url":"https://example.com/logo.png","founded_year":1928,"address":"Jl. Pintu Satu Senayan","city":"Jakarta"}`
	req := httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))

	teamData := resp["data"].(map[string]interface{})
	teamID := teamData["id"].(string)
	assert.NotEmpty(t, teamID)

	// Get team
	req = httptest.NewRequest("GET", "/api/teams/"+teamID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Persija Jakarta", resp["data"].(map[string]interface{})["name"])

	// Update team
	body = `{"name":"Persija Jakarta Updated","city":"Jakarta Pusat","version":1}`
	req = httptest.NewRequest("PUT", "/api/teams/"+teamID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// List teams
	req = httptest.NewRequest("GET", "/api/teams", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["total"])

	// Delete team
	req = httptest.NewRequest("DELETE", "/api/teams/"+teamID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify deleted
	req = httptest.NewRequest("GET", "/api/teams", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["total"])
}

func TestIntegration_PlayerCRUD(t *testing.T) {
	r, db := setupIntegrationRouter(t)

	// Create team first
	teamRepo := repository.NewTeamRepository(db)
	team := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team)

	// Create player
	body := `{"team_id":"` + team.ID + `","name":"Bambang Pamungkas","height":178.5,"weight":72.0,"position":"CF","playstyle":"goal_poacher","jersey_number":20}`
	req := httptest.NewRequest("POST", "/api/players", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	playerData := resp["data"].(map[string]interface{})
	playerID := playerData["id"].(string)

	// Get player
	req = httptest.NewRequest("GET", "/api/players/"+playerID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Try duplicate jersey
	body = `{"team_id":"` + team.ID + `","name":"Another Player","height":180.0,"weight":75.0,"position":"CMF","jersey_number":20}`
	req = httptest.NewRequest("POST", "/api/players", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestIntegration_MatchAndResult(t *testing.T) {
	r, db := setupIntegrationRouter(t)

	// Create teams
	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Addr"}
	teamRepo.Create(team1)
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Addr"}
	teamRepo.Create(team2)

	// Create players
	playerRepo := repository.NewPlayerRepository(db)
	player1 := &domain.Player{TeamID: team1.ID, Name: "Bambang", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10, Version: 1}
	playerRepo.Create(player1)

	// Schedule match
	body := `{"match_date":"2026-09-15","match_time":"19:30:00","home_team_id":"` + team1.ID + `","away_team_id":"` + team2.ID + `"}`
	req := httptest.NewRequest("POST", "/api/matches", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	matchData := resp["data"].(map[string]interface{})
	matchID := matchData["id"].(string)

	// Report result
	body = `{"match_id":"` + matchID + `","home_score":3,"away_score":1,"goals":[{"player_id":"` + player1.ID + `","goal_time":"23'"},{"player_id":"` + player1.ID + `","goal_time":"45'"},{"player_id":"` + player1.ID + `","goal_time":"67'"}]}`
	req = httptest.NewRequest("POST", "/api/results", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Get match report
	req = httptest.NewRequest("GET", "/api/reports/matches/"+matchID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	reportData := resp["data"].(map[string]interface{})
	assert.Equal(t, "Tim Home Menang", reportData["status"])
	assert.Equal(t, float64(3), reportData["home_score"])
	assert.Equal(t, float64(1), reportData["away_score"])
}

func TestIntegration_HealthCheck(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "ok", resp["data"].(map[string]interface{})["status"])
}
