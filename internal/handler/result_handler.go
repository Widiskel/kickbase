package handler

import (
	"net/http"

	"kickbase/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
