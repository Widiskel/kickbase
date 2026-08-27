package handler

import (
	"net/http"
	"strconv"

	"kickbase/internal/domain"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Team handlers

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new football team with the provided details
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body domain.Team true "Team data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/teams [post]
func CreateTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var team domain.Team
		if err := c.ShouldBindJSON(&team); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}
		if err := db.Create(&team).Error; err != nil {
			RespondError(c, http.StatusConflict, "CONFLICT", "Team name already exists", nil)
			return
		}
		RespondCreated(c, team)
	}
}

// ListTeams godoc
// @Summary List all teams
// @Description Get a paginated list of all non-deleted teams
// @Tags Teams
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /api/teams [get]
func ListTeams(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var teams []domain.Team
		var total int64

		db.Model(&domain.Team{}).Count(&total)
		db.Offset((page - 1) * limit).Limit(limit).Find(&teams)

		RespondPaginated(c, teams, total, page, limit)
	}
}

// GetTeam godoc
// @Summary Get a team by ID
// @Description Get a single team by its ID
// @Tags Teams
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/teams/{id} [get]
func GetTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var team domain.Team
		if err := db.First(&team, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Team not found", nil)
			return
		}
		RespondSuccess(c, team, "")
	}
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update an existing team's information
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param team body domain.Team true "Team data with version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/teams/{id} [put]
func UpdateTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var team domain.Team
		if err := db.First(&team, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Team not found", nil)
			return
		}

		var input domain.Team
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		if input.Version != team.Version {
			RespondError(c, http.StatusConflict, "CONFLICT", "Version mismatch", nil)
			return
		}

		updates := map[string]interface{}{
			"name":         input.Name,
			"logo_url":     input.LogoURL,
			"founded_year": input.FoundedYear,
			"address":      input.Address,
			"city":         input.City,
			"version":      team.Version + 1,
		}

		db.Model(&team).Updates(updates)
		RespondSuccess(c, team, "Team updated successfully")
	}
}

// DeleteTeam godoc
// @Summary Delete a team
// @Description Soft delete a team (only if no active players)
// @Tags Teams
// @Produce json
// @Param id path string true "Team ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/teams/{id} [delete]
func DeleteTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var team domain.Team
		if err := db.First(&team, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Team not found", nil)
			return
		}

		// Check if team has active players
		var playerCount int64
		db.Model(&domain.Player{}).Where("team_id = ?", id).Count(&playerCount)
		if playerCount > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Cannot delete team with active players", nil)
			return
		}

		db.Delete(&team)
		RespondNoContent(c)
	}
}

func GetTeamHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var history []domain.TeamHistory
		db.Where("team_id = ?", id).Order("version").Find(&history)
		RespondSuccess(c, history, "")
	}
}

func RevertTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input struct {
			TargetVersion int `json:"target_version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		teamRepo := repository.NewTeamRepository(db)
		svc := service.NewTeamService(teamRepo)

		if err := svc.RevertTeam(id, input.TargetVersion); err != nil {
			if err.Error() == "team not found" || err.Error() == "target version not found in history" {
				RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert team", nil)
			}
			return
		}

		RespondSuccess(c, nil, "Team reverted successfully")
	}
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Add a new player to a team
// @Tags Players
// @Accept json
// @Produce json
// @Param player body domain.Player true "Player data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/players [post]
func CreatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player domain.Player
		if err := c.ShouldBindJSON(&player); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		// Validate position
		validPositions := map[string]bool{
			"CF": true, "SS": true, "LWF": true, "RWF": true,
			"AMF": true, "CMF": true, "DMF": true, "LMF": true, "RMF": true,
			"CB": true, "LB": true, "RB": true, "GK": true,
		}
		if !validPositions[player.Position] {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid position. Valid positions: CF, SS, LWF, RWF, AMF, CMF, DMF, LMF, RMF, CB, LB, RB, GK", nil)
			return
		}

		// Check jersey uniqueness
		var count int64
		db.Model(&domain.Player{}).Where("team_id = ? AND jersey_number = ?", player.TeamID, player.JerseyNumber).Count(&count)
		if count > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Jersey number already exists in this team", nil)
			return
		}

		if err := db.Create(&player).Error; err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create player", nil)
			return
		}
		RespondCreated(c, player)
	}
}

func ListPlayersByTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Param("teamId")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var players []domain.Player
		var total int64

		db.Model(&domain.Player{}).Where("team_id = ?", teamID).Count(&total)
		db.Where("team_id = ?", teamID).Offset((page - 1) * limit).Limit(limit).Find(&players)

		RespondPaginated(c, players, total, page, limit)
	}
}

func ListPlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Query("team_id")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var players []domain.Player
		var total int64

		query := db.Model(&domain.Player{})
		if teamID != "" {
			query = query.Where("team_id = ?", teamID)
		}

		query.Count(&total)
		query.Offset((page - 1) * limit).Limit(limit).Find(&players)

		RespondPaginated(c, players, total, page, limit)
	}
}

func GetPlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}
		RespondSuccess(c, player, "")
	}
}

func UpdatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}

		var input domain.Player
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		if input.Version != player.Version {
			RespondError(c, http.StatusConflict, "CONFLICT", "Version mismatch", nil)
			return
		}

		updates := map[string]interface{}{
			"name":          input.Name,
			"height":        input.Height,
			"weight":        input.Weight,
			"position":      input.Position,
			"playstyle":     input.Playstyle,
			"jersey_number": input.JerseyNumber,
			"version":       player.Version + 1,
		}

		db.Model(&player).Updates(updates)
		RespondSuccess(c, player, "Player updated successfully")
	}
}

func DeletePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}

		// Check if player has goal records
		var goalCount int64
		db.Model(&domain.Goal{}).Where("player_id = ?", id).Count(&goalCount)
		if goalCount > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Cannot delete player with goal records", nil)
			return
		}

		db.Delete(&player)
		RespondNoContent(c)
	}
}

func GetPlayerHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var history []domain.PlayerHistory
		db.Where("player_id = ?", id).Order("version").Find(&history)
		RespondSuccess(c, history, "")
	}
}

func RevertPlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input struct {
			TargetVersion int `json:"target_version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		playerRepo := repository.NewPlayerRepository(db)
		teamRepo := repository.NewTeamRepository(db)
		svc := service.NewPlayerService(playerRepo, teamRepo)

		if err := svc.RevertPlayer(id, input.TargetVersion); err != nil {
			if err.Error() == "player not found" || err.Error() == "target version not found in history" {
				RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert player", nil)
			}
			return
		}

		RespondSuccess(c, nil, "Player reverted successfully")
	}
}

// CreateMatch godoc
// @Summary Schedule a new match
// @Description Create a match schedule between two teams
// @Tags Matches
// @Accept json
// @Produce json
// @Param match body domain.Match true "Match data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/matches [post]
func CreateMatch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var match domain.Match
		if err := c.ShouldBindJSON(&match); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		if match.HomeTeamID == match.AwayTeamID {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Home team and away team must be different", nil)
			return
		}

		match.Status = "scheduled"
		if err := db.Create(&match).Error; err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create match", nil)
			return
		}
		RespondCreated(c, match)
	}
}

func ListMatches(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var matches []domain.Match
		var total int64

		db.Model(&domain.Match{}).Count(&total)
		db.Offset((page - 1) * limit).Limit(limit).Find(&matches)

		RespondPaginated(c, matches, total, page, limit)
	}
}

func GetMatch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var match domain.Match
		if err := db.First(&match, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match not found", nil)
			return
		}
		RespondSuccess(c, match, "")
	}
}

func UpdateMatchStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var match domain.Match
		if err := db.First(&match, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match not found", nil)
			return
		}

		var input struct {
			Status  string `json:"status"`
			Version int    `json:"version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		if input.Version != match.Version {
			RespondError(c, http.StatusConflict, "CONFLICT", "Version mismatch", nil)
			return
		}

		db.Model(&match).Updates(map[string]interface{}{
			"status":  input.Status,
			"version": match.Version + 1,
		})
		RespondSuccess(c, match, "Match status updated")
	}
}

func GetMatchHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var history []domain.MatchHistory
		db.Where("match_id = ?", id).Order("version").Find(&history)
		RespondSuccess(c, history, "")
	}
}

func RevertMatch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input struct {
			TargetVersion int `json:"target_version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		matchRepo := repository.NewMatchRepository(db)
		teamRepo := repository.NewTeamRepository(db)
		svc := service.NewMatchService(matchRepo, teamRepo)

		if err := svc.RevertMatch(id, input.TargetVersion); err != nil {
			if err.Error() == "match not found" || err.Error() == "target version not found in history" {
				RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert match", nil)
			}
			return
		}

		RespondSuccess(c, nil, "Match reverted successfully")
	}
}

// CreateMatchResult godoc
// @Summary Report match result
// @Description Report the result of a completed match with goals
// @Tags Results
// @Accept json
// @Produce json
// @Param result body object true "Match result with goals"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/results [post]
func CreateMatchResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			MatchID   string `json:"match_id"`
			HomeScore int    `json:"home_score"`
			AwayScore int    `json:"away_score"`
			Goals     []struct {
				PlayerID string `json:"player_id"`
				GoalTime string `json:"goal_time"`
			} `json:"goals"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		// Check match exists and is scheduled
		var match domain.Match
		if err := db.First(&match, "id = ?", input.MatchID).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match not found", nil)
			return
		}

		if match.Status != "scheduled" {
			RespondError(c, http.StatusConflict, "CONFLICT", "Match is not in scheduled status", nil)
			return
		}

		// Check no existing result
		var count int64
		db.Model(&domain.MatchResult{}).Where("match_id = ?", input.MatchID).Count(&count)
		if count > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Match result already exists", nil)
			return
		}

		// Create result in transaction
		tx := db.Begin()
		result := domain.MatchResult{
			MatchID:   input.MatchID,
			HomeScore: input.HomeScore,
			AwayScore: input.AwayScore,
		}
		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create result", nil)
			return
		}

		// Create goals
		for _, g := range input.Goals {
			goal := domain.Goal{
				MatchResultID: result.ID,
				PlayerID:      g.PlayerID,
				GoalTime:      g.GoalTime,
			}
			if err := tx.Create(&goal).Error; err != nil {
				tx.Rollback()
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create goal", nil)
				return
			}
		}

		// Update match status to completed
		tx.Model(&match).Updates(map[string]interface{}{
			"status":  "completed",
			"version": match.Version + 1,
		})

		tx.Commit()
		RespondCreated(c, result)
	}
}

func GetMatchResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		matchID := c.Query("match_id")
		if matchID == "" {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "match_id query parameter is required", nil)
			return
		}

		var result domain.MatchResult
		if err := db.Where("match_id = ?", matchID).First(&result).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match result not found", nil)
			return
		}

		var goals []domain.Goal
		db.Where("match_result_id = ?", result.ID).Find(&goals)

		RespondSuccess(c, gin.H{"result": result, "goals": goals}, "")
	}
}

// ListMatchReports godoc
// @Summary List all match reports
// @Description Get a paginated list of all match reports with scores and status
// @Tags Reports
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /api/reports/matches [get]
func ListMatchReports(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		// Create service
		matchRepo := repository.NewMatchRepository(db)
		resultRepo := repository.NewResultRepository(db)
		goalRepo := repository.NewGoalRepository(db)
		teamRepo := repository.NewTeamRepository(db)
		playerRepo := repository.NewPlayerRepository(db)
		reportSvc := service.NewReportService(db, matchRepo, resultRepo, goalRepo, teamRepo, playerRepo)

		reports, total, err := reportSvc.ListMatchReports(page, limit)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list reports", nil)
			return
		}

		RespondPaginated(c, reports, total, page, limit)
	}
}

// GetMatchReport godoc
// @Summary Get a single match report
// @Description Get detailed report for a specific match including score, status, top scorer, and cumulative wins
// @Tags Reports
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/reports/matches/{id} [get]
func GetMatchReport(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Create service
		matchRepo := repository.NewMatchRepository(db)
		resultRepo := repository.NewResultRepository(db)
		goalRepo := repository.NewGoalRepository(db)
		teamRepo := repository.NewTeamRepository(db)
		playerRepo := repository.NewPlayerRepository(db)
		reportSvc := service.NewReportService(db, matchRepo, resultRepo, goalRepo, teamRepo, playerRepo)

		report, err := reportSvc.GetMatchReport(id)
		if err != nil {
			if err.Error() == "match not found" {
				RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match not found", nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get report", nil)
			}
			return
		}

		RespondSuccess(c, report, "")
	}
}
