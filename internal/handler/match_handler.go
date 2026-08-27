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
