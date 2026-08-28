package handler

import (
	"net/http"
	"strings"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService interfaces.MatchService
}

func NewMatchHandler(matchService interfaces.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

// CreateMatch godoc
// @Summary Schedule a new match
// @Description Schedule a match between two different teams
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param match body CreateMatchRequest true "Match schedule data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/matches [post]
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	match := req.ToDomain()
	if err := h.matchService.CreateMatch(match); err != nil {
		switch err.Error() {
		case "home team and away team must be different", "home and away teams must be different", "invalid match date format", "invalid match time format":
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		case "home team not found", "away team not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to schedule match", nil)
		}
		return
	}

	RespondCreated(c, match)
}

// ListMatches godoc
// @Summary List all matches
// @Description Get a paginated list of matches with optional team, status, date filtering and sorting
// @Tags Matches
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param team_id query string false "Filter by team ID (home or away)"
// @Param status query string false "Filter by match status (scheduled, completed, cancelled, deferred)"
// @Param date_from query string false "Filter matches from date (YYYY-MM-DD)"
// @Param date_to query string false "Filter matches to date (YYYY-MM-DD)"
// @Param sort_by query string false "Sort field (match_date, created_at)" default(match_date)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} PaginatedResponse
// @Router /api/matches [get]
func (h *MatchHandler) ListMatches(c *gin.Context) {
	page, limit := GetPagination(c)

	opts := interfaces.MatchFilterOptions{
		Page:     page,
		Limit:    limit,
		TeamID:   c.Query("team_id"),
		Status:   c.Query("status"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		SortBy:   c.Query("sort_by"),
		Order:    c.Query("order"),
	}

	matches, total, err := h.matchService.ListMatches(opts)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list matches", nil)
		return
	}

	RespondPaginated(c, matches, total, page, limit)
}

// GetMatch godoc
// @Summary Get a match by ID
// @Description Get a single match by its ID
// @Tags Matches
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/matches/{id} [get]
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id := c.Param("id")

	match, err := h.matchService.GetMatch(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match not found", nil)
		return
	}

	RespondSuccess(c, match, "")
}

// UpdateMatchStatus godoc
// @Summary Update match status
// @Description Update the status of a match (scheduled, completed, cancelled, deferred)
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Param status body UpdateMatchStatusRequest true "New status"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/matches/{id}/status [patch]
func (h *MatchHandler) UpdateMatchStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if err := h.matchService.UpdateMatchStatus(id, req.Status); err != nil {
		switch err.Error() {
		case "match not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		default:
			if strings.Contains(err.Error(), "invalid status transition") {
				RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update match status", nil)
			}
		}
		return
	}

	RespondSuccess(c, nil, "Match status updated")
}

// GetMatchHistory godoc
// @Summary Get match history
// @Description Get the version history of a match
// @Tags Matches
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Router /api/matches/{id}/history [get]
func (h *MatchHandler) GetMatchHistory(c *gin.Context) {
	id := c.Param("id")

	history, err := h.matchService.GetMatchHistory(id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get history", nil)
		return
	}

	RespondSuccess(c, history, "")
}

// RevertMatch godoc
// @Summary Revert match to previous version
// @Description Revert a match to a specific previous version
// @Tags Matches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Match ID"
// @Param request body RevertRequest true "Target version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/matches/{id}/revert [post]
func (h *MatchHandler) RevertMatch(c *gin.Context) {
	id := c.Param("id")

	var req RevertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if err := h.matchService.RevertMatch(id, req.TargetVersion); err != nil {
		switch err.Error() {
		case "match not found", "target version not found in history":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert match", nil)
		}
		return
	}

	RespondSuccess(c, nil, "Match reverted successfully")
}
