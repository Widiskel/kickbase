package handler

import (
	"net/http"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	reportService interfaces.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService interfaces.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// ListMatchReports godoc
// @Summary List all match reports
// @Description Get a paginated list of match reports with optional team and result status filtering
// @Tags Reports
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param team_id query string false "Filter by team ID"
// @Param status query string false "Filter by match status / outcome (Tim Home Menang, Tim Away Menang, Draw, scheduled)"
// @Param sort_by query string false "Sort field (match_date, created_at)" default(match_date)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Success 200 {object} PaginatedResponse
// @Router /api/reports/matches [get]
func (h *ReportHandler) ListMatchReports(c *gin.Context) {
	page, limit := GetPagination(c)
	opts := interfaces.ReportFilterOptions{
		Page:   page,
		Limit:  limit,
		TeamID: c.Query("team_id"),
		Status: c.Query("status"),
		SortBy: c.Query("sort_by"),
		Order:  c.Query("order"),
	}

	reports, total, err := h.reportService.ListMatchReports(opts)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list reports", nil)
		return
	}

	RespondPaginated(c, reports, total, page, limit)
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
func (h *ReportHandler) GetMatchReport(c *gin.Context) {
	id := c.Param("id")

	report, err := h.reportService.GetMatchReport(id)
	if err != nil {
		if err.Error() == "match not found" {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		} else {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get report", nil)
		}
		return
	}

	RespondSuccess(c, report, "")
}
