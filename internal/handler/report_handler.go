package handler

import (
	"net/http"
	"strconv"

	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
