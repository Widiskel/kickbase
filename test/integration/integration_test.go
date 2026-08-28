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
	db.Exec("DROP TABLE IF EXISTS goals, match_results, matches, players, teams, users, refresh_tokens, team_histories, player_histories, match_histories, match_result_histories, goal_histories CASCADE")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
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

func getTokenForUser(t *testing.T, r *gin.Engine, username, password, name, role string) (string, string) {
	// Register user
	regBody := `{"username":"` + username + `","password":"` + password + `","name":"` + name + `","role":"` + role + `"}`
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Login
	loginBody := `{"username":"` + username + `","password":"` + password + `"}`
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)
	refreshToken := data["refresh_token"].(string)
	return accessToken, refreshToken
}

func TestIntegration_AuthAndTokenRefresh(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	// 1. Register & Login Admin
	accessToken, refreshToken := getTokenForUser(t, r, "admin_test", "password123", "Admin Test", "admin")
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 2. Refresh Token Flow
	refreshBody := `{"refresh_token":"` + refreshToken + `"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	newAccessToken := data["access_token"].(string)
	newRefreshToken := data["refresh_token"].(string)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, refreshToken, newRefreshToken) // Token rotated

	// 3. Old refresh token should now be rejected
	req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBufferString(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_GranularRolePermissions(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	// Get tokens for staff and viewer
	staffToken, _ := getTokenForUser(t, r, "staff_user", "password123", "Staff User", "staff")
	viewerToken, _ := getTokenForUser(t, r, "viewer_user", "password123", "Viewer User", "viewer")

	// 1. Staff can create team (has teams:create)
	teamBody := `{"name":"Persija Jakarta","logo_url":"https://example.com/logo.png","founded_year":1928,"address":"Senayan","city":"Jakarta"}`
	req := httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(teamBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+staffToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	teamID := resp["data"].(map[string]interface{})["id"].(string)

	// 2. Staff CANNOT delete team (lacks teams:delete) -> 403 Forbidden
	req = httptest.NewRequest("DELETE", "/api/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+staffToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Viewer CANNOT create team (lacks teams:create) -> 403 Forbidden
	viewerTeamBody := `{"name":"Persib Bandung","logo_url":"https://example.com/logo.png","founded_year":1933,"address":"Gedebage","city":"Bandung"}`
	req = httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(viewerTeamBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+viewerToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 4. Viewer CAN read teams (has teams:read) -> 200 OK
	req = httptest.NewRequest("GET", "/api/teams", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TeamCRUD(t *testing.T) {
	r, _ := setupIntegrationRouter(t)
	token, _ := getTokenForUser(t, r, "admin_test", "password123", "Admin Test", "admin")

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
	token, _ := getTokenForUser(t, r, "admin_test", "password123", "Admin Test", "admin")

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
	token, _ := getTokenForUser(t, r, "admin_test", "password123", "Admin Test", "admin")

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

	var matchResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &matchResp)
	matchData := matchResp["data"].(map[string]interface{})
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
	var reportResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &reportResp)
	reportData := reportResp["data"].(map[string]interface{})
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

func TestIntegration_FilteringAndSorting(t *testing.T) {
	r, db := setupIntegrationRouter(t)
	token, _ := getTokenForUser(t, r, "admin_filter", "password123", "Admin Filter", "admin")

	// Create two teams
	body1 := `{"name":"Persija Jakarta","founded_year":1928,"address":"Senayan","city":"Jakarta"}`
	body2 := `{"name":"Persib Bandung","founded_year":1933,"address":"Gedebage","city":"Bandung"}`
	
	req1 := httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(body1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	req2 := httptest.NewRequest("POST", "/api/teams", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	// Filter by city: Jakarta
	req := httptest.NewRequest("GET", "/api/teams?city=Jakarta", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["total"])

	// Sort by founded_year desc (Persib 1933 should be first)
	req = httptest.NewRequest("GET", "/api/teams?sort_by=founded_year&order=desc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["data"].([]interface{})
	assert.Equal(t, "Persib Bandung", items[0].(map[string]interface{})["name"])

	_ = db
}

func TestIntegration_LivePrometheusMetrics(t *testing.T) {
	r, _ := setupIntegrationRouter(t)

	// Scrape /metrics
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	bodyStr := w.Body.String()
	assert.Contains(t, bodyStr, "http_requests_total")
	assert.Contains(t, bodyStr, "teams_total")
	assert.Contains(t, bodyStr, "players_total")
	assert.Contains(t, bodyStr, "matches_total")
}
