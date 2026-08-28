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
	db.Exec("DROP TABLE IF EXISTS goals, match_results, matches, players, teams, users, team_histories, player_histories, match_histories, match_result_histories, goal_histories CASCADE")

	err = db.AutoMigrate(
		&domain.User{},
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

func getAdminToken(t *testing.T, r *gin.Engine) string {
	// Register admin
	regBody := `{"username":"admin_test","password":"password123","name":"Admin Test","role":"admin"}`
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Login
	loginBody := `{"username":"admin_test","password":"password123"}`
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	return data["token"].(string)
}

func TestIntegration_AuthFlow(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	// 1. Register User
	regBody := `{"username":"viewer_user","password":"password123","name":"Viewer User","role":"viewer"}`
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Login Viewer User
	loginBody := `{"username":"viewer_user","password":"password123"}`
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	viewerToken := data["token"].(string)
	assert.NotEmpty(t, viewerToken)

	// 3. Unauthorized access (without token)
	teamBody := `{"name":"Persija Jakarta","logo_url":"https://example.com/logo.png","founded_year":1928,"address":"Senayan","city":"Jakarta"}`
	req = httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(teamBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 4. Forbidden access (viewer trying to mutate admin resource)
	req = httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(teamBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+viewerToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestIntegration_TeamCRUD(t *testing.T) {
	r, _ := setupIntegrationRouter(t)
	token := getAdminToken(t, r)

	// Create team
	body := `{"name":"Persija Jakarta","logo_url":"https://example.com/logo.png","founded_year":1928,"address":"Jl. Pintu Satu Senayan","city":"Jakarta"}`
	req := httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))

	teamData := resp["data"].(map[string]interface{})
	teamID := teamData["id"].(string)
	assert.NotEmpty(t, teamID)

	// Get team (Public)
	req = httptest.NewRequest("GET", "/api/teams/"+teamID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Persija Jakarta", resp["data"].(map[string]interface{})["name"])

	// Update team (Protected)
	body = `{"name":"Persija Jakarta Updated","city":"Jakarta Pusat","version":1}`
	req = httptest.NewRequest("PUT", "/api/teams/"+teamID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// List teams (Public)
	req = httptest.NewRequest("GET", "/api/teams", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["total"])

	// Delete team (Protected)
	req = httptest.NewRequest("DELETE", "/api/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// List should be empty (soft deleted)
	req = httptest.NewRequest("GET", "/api/teams", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["total"])
}

func TestIntegration_PlayerCRUD(t *testing.T) {
	r, db := setupIntegrationRouter(t)
	token := getAdminToken(t, r)

	// Create a team first
	teamRepo := repository.NewTeamRepository(db)
	team := &domain.Team{Name: "Persija Jakarta", City: "Jakarta", FoundedYear: 1928, Address: "Senayan"}
	teamRepo.Create(team)

	// Create player (Protected)
	body := `{"team_id":"` + team.ID + `","name":"Bambang Pamungkas","height":178,"weight":72,"position":"CF","jersey_number":20}`
	req := httptest.NewRequest("POST", "/api/players", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	playerData := resp["data"].(map[string]interface{})
	playerID := playerData["id"].(string)
	assert.NotEmpty(t, playerID)

	// Get player (Public)
	req = httptest.NewRequest("GET", "/api/players/"+playerID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Duplicate jersey number in same team should fail
	body = `{"team_id":"` + team.ID + `","name":"Another Player","height":180,"weight":75,"position":"CMF","jersey_number":20}`
	req = httptest.NewRequest("POST", "/api/players", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestIntegration_MatchAndResult(t *testing.T) {
	r, db := setupIntegrationRouter(t)
	token := getAdminToken(t, r)

	// Create teams
	teamRepo := repository.NewTeamRepository(db)
	team1 := &domain.Team{Name: "Persija", City: "Jakarta", FoundedYear: 1928, Address: "Senayan"}
	team2 := &domain.Team{Name: "Persib", City: "Bandung", FoundedYear: 1933, Address: "Gedebage"}
	teamRepo.Create(team1)
	teamRepo.Create(team2)

	// Create players
	playerRepo := repository.NewPlayerRepository(db)
	p1 := &domain.Player{TeamID: team1.ID, Name: "Player 1", Height: 175, Weight: 70, Position: "CF", JerseyNumber: 10}
	p2 := &domain.Player{TeamID: team2.ID, Name: "Player 2", Height: 180, Weight: 75, Position: "CMF", JerseyNumber: 8}
	playerRepo.Create(p1)
	playerRepo.Create(p2)

	// Schedule match (Protected)
	matchBody := `{"match_date":"2026-09-15","match_time":"19:30:00","home_team_id":"` + team1.ID + `","away_team_id":"` + team2.ID + `"}`
	req := httptest.NewRequest("POST", "/api/matches", bytes.NewBufferString(matchBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	matchData := resp["data"].(map[string]interface{})
	matchID := matchData["id"].(string)

	// Report result (Protected)
	resultBody := `{"match_id":"` + matchID + `","home_score":2,"away_score":1,"goals":[{"player_id":"` + p1.ID + `","goal_time":"15'"},{"player_id":"` + p1.ID + `","goal_time":"45'"},{"player_id":"` + p2.ID + `","goal_time":"60'"}]}`
	req = httptest.NewRequest("POST", "/api/results", bytes.NewBufferString(resultBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Get report (Public)
	req = httptest.NewRequest("GET", "/api/reports/matches/"+matchID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	reportData := resp["data"].(map[string]interface{})
	assert.Equal(t, "Tim Home Menang", reportData["status"])
	assert.Equal(t, float64(2), reportData["home_score"])
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
	assert.Equal(t, "Service is healthy", resp["message"])
}
