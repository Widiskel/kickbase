package handler

import (
	"net/http"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// ResultHandler handles match result-related HTTP requests
type ResultHandler struct {
	resultService interfaces.ResultService
}

// NewResultHandler creates a new ResultHandler
func NewResultHandler(resultService interfaces.ResultService) *ResultHandler {
	return &ResultHandler{resultService: resultService}
}

// CreateMatchResult godoc
// @Summary Report match result
// @Description Report the result of a completed match with goals
// @Tags Results
// @Accept json
// @Produce json
// @Param result body CreateResultRequest true "Match result with goals"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/results [post]
func (h *ResultHandler) CreateMatchResult(c *gin.Context) {
	var req CreateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	input := req.ToServiceInput()
	result, err := h.resultService.CreateResult(input)
	if err != nil {
		switch err.Error() {
		case "match not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "match is not in scheduled status", "match result already exists":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		default:
			if contains(err.Error(), "does not belong to either team") || contains(err.Error(), "player") {
				RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create result", nil)
			}
		}
		return
	}

	RespondCreated(c, result)
}

// GetMatchResult godoc
// @Summary Get match result
// @Description Get the result of a specific match
// @Tags Results
// @Produce json
// @Param matchId query string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/results/{matchId} [get]
func (h *ResultHandler) GetMatchResult(c *gin.Context) {
	matchID := c.Query("match_id")
	if matchID == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "match_id query parameter is required", nil)
		return
	}

	result, goals, err := h.resultService.GetResult(matchID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match result not found", nil)
		return
	}

	RespondSuccess(c, gin.H{"result": result, "goals": goals}, "")
}
