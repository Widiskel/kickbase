package handler

import (
	"net/http"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// TeamHandler handles team-related HTTP requests
type TeamHandler struct {
	teamService interfaces.TeamService
}

// NewTeamHandler creates a new TeamHandler
func NewTeamHandler(teamService interfaces.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new football team with the provided details
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body CreateTeamRequest true "Team data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	team := req.ToDomain()
	if err := h.teamService.CreateTeam(team); err != nil {
		if err.Error() == "team name already exists" {
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		} else {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create team", nil)
		}
		return
	}

	RespondCreated(c, team)
}

// ListTeams godoc
// @Summary List all teams
// @Description Get a paginated list of teams with optional search and sorting
// @Tags Teams
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param name query string false "Filter by team name"
// @Param city query string false "Filter by headquarters city"
// @Param sort_by query string false "Sort field (name, city, founded_year, created_at)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Success 200 {object} PaginatedResponse
// @Router /api/teams [get]
func (h *TeamHandler) ListTeams(c *gin.Context) {
	page, limit := GetPagination(c)
	opts := interfaces.TeamFilterOptions{
		Page:    page,
		Limit:   limit,
		Name:    c.Query("name"),
		City:    c.Query("city"),
		SortBy:  c.Query("sort_by"),
		Order:   c.Query("order"),
	}

	teams, total, err := h.teamService.ListTeams(opts)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list teams", nil)
		return
	}

	RespondPaginated(c, teams, total, page, limit)
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
func (h *TeamHandler) GetTeam(c *gin.Context) {
	id := c.Param("id")

	team, err := h.teamService.GetTeam(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Team not found", nil)
		return
	}

	RespondSuccess(c, team, "")
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update an existing team's information
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param team body UpdateTeamRequest true "Team data with version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/teams/{id} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	team := req.ToDomain(id)
	if err := h.teamService.UpdateTeam(team); err != nil {
		switch err.Error() {
		case "team not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "version mismatch":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update team", nil)
		}
		return
	}

	RespondSuccess(c, team, "Team updated successfully")
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
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	if err := h.teamService.DeleteTeam(id); err != nil {
		switch err.Error() {
		case "team not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "cannot delete team with active players":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete team", nil)
		}
		return
	}

	RespondNoContent(c)
}

// GetTeamHistory godoc
// @Summary Get team history
// @Description Get the version history of a team
// @Tags Teams
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} SuccessResponse
// @Router /api/teams/{id}/history [get]
func (h *TeamHandler) GetTeamHistory(c *gin.Context) {
	id := c.Param("id")

	history, err := h.teamService.GetTeamHistory(id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get history", nil)
		return
	}

	RespondSuccess(c, history, "")
}

// RevertTeam godoc
// @Summary Revert team to previous version
// @Description Revert a team to a specific previous version
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param request body RevertRequest true "Target version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/teams/{id}/revert [post]
func (h *TeamHandler) RevertTeam(c *gin.Context) {
	id := c.Param("id")

	var req RevertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if err := h.teamService.RevertTeam(id, req.TargetVersion); err != nil {
		switch err.Error() {
		case "team not found", "target version not found in history":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert team", nil)
		}
		return
	}

	RespondSuccess(c, nil, "Team reverted successfully")
}
