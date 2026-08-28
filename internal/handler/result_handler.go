package handler

import (
	"net/http"
	"strings"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type ResultHandler struct {
	resultService interfaces.ResultService
}

func NewResultHandler(resultService interfaces.ResultService) *ResultHandler {
	return &ResultHandler{resultService: resultService}
}

// CreateMatchResult godoc
// @Summary Report match result
// @Description Record the final score and goals for a completed match
// @Tags Results
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param result body CreateResultRequest true "Match result data with goals"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/results [post]
func (h *ResultHandler) CreateMatchResult(c *gin.Context) {
	var req CreateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	result, err := h.resultService.CreateResult(req.ToServiceInput())
	if err != nil {
		switch err.Error() {
		case "match not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "match result already exists":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		case "can only report result for scheduled matches", "match is not in scheduled status":
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		case "score cannot be negative", "scores cannot be negative", "number of goals does not match total score":
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		default:
			if strings.Contains(err.Error(), "must match total score") || strings.Contains(err.Error(), "does not belong to either team") || strings.Contains(err.Error(), "player") {
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
// @Param matchId path string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/results/{matchId} [get]
func (h *ResultHandler) GetMatchResult(c *gin.Context) {
	matchID := c.Param("matchId")
	if matchID == "" {
		matchID = c.Query("match_id")
	}
	if matchID == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "match_id parameter is required", nil)
		return
	}

	result, goals, err := h.resultService.GetResult(matchID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Match result not found", nil)
		return
	}

	RespondSuccess(c, gin.H{"result": result, "goals": goals}, "")
}
