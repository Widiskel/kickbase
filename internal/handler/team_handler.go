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
