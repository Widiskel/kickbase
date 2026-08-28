package handler

import (
	"net/http"
	"strconv"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// PlayerHandler handles player-related HTTP requests
type PlayerHandler struct {
	playerService interfaces.PlayerService
}

// NewPlayerHandler creates a new PlayerHandler
func NewPlayerHandler(playerService interfaces.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Add a new player to a team
// @Tags Players
// @Accept json
// @Produce json
// @Param player body CreatePlayerRequest true "Player data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/players [post]
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	player := req.ToDomain()
	if err := h.playerService.CreatePlayer(player); err != nil {
		switch err.Error() {
		case "team not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "jersey number already exists in this team":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		case "invalid position":
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create player", nil)
		}
		return
	}

	RespondCreated(c, player)
}

// ListPlayers godoc
// @Summary List all players
// @Description Get a paginated list of players, optionally filtered by team
// @Tags Players
// @Produce json
// @Param team_id query string false "Filter by team ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /api/players [get]
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	teamID := c.Query("team_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if teamID != "" {
		players, total, err := h.playerService.ListPlayersByTeam(teamID, page, limit)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list players", nil)
			return
		}
		RespondPaginated(c, players, total, page, limit)
	} else {
		// List all players - for now, return empty (could add ListAllPlayers to service)
		RespondPaginated(c, []interface{}{}, 0, page, limit)
	}
}

// GetPlayer godoc
// @Summary Get a player by ID
// @Description Get a single player by its ID
// @Tags Players
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/players/{id} [get]
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")

	player, err := h.playerService.GetPlayer(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
		return
	}

	RespondSuccess(c, player, "")
}

// UpdatePlayer godoc
// @Summary Update a player
// @Description Update an existing player's information
// @Tags Players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param player body UpdatePlayerRequest true "Player data with version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/players/{id} [put]
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	player := req.ToDomain(id)
	if err := h.playerService.UpdatePlayer(player); err != nil {
		switch err.Error() {
		case "player not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "version mismatch":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		case "jersey number already exists in this team":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		case "invalid position":
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update player", nil)
		}
		return
	}

	RespondSuccess(c, player, "Player updated successfully")
}

// DeletePlayer godoc
// @Summary Delete a player
// @Description Soft delete a player (only if no goal records)
// @Tags Players
// @Produce json
// @Param id path string true "Player ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/players/{id} [delete]
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	if err := h.playerService.DeletePlayer(id); err != nil {
		switch err.Error() {
		case "player not found":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		case "cannot delete player with goal records":
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete player", nil)
		}
		return
	}

	RespondNoContent(c)
}

// GetPlayerHistory godoc
// @Summary Get player history
// @Description Get the version history of a player
// @Tags Players
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} SuccessResponse
// @Router /api/players/{id}/history [get]
func (h *PlayerHandler) GetPlayerHistory(c *gin.Context) {
	id := c.Param("id")

	history, err := h.playerService.GetPlayerHistory(id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get history", nil)
		return
	}

	RespondSuccess(c, history, "")
}

// RevertPlayer godoc
// @Summary Revert player to previous version
// @Description Revert a player to a specific previous version
// @Tags Players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param request body RevertRequest true "Target version"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/players/{id}/revert [post]
func (h *PlayerHandler) RevertPlayer(c *gin.Context) {
	id := c.Param("id")

	var req RevertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if err := h.playerService.RevertPlayer(id, req.TargetVersion); err != nil {
		switch err.Error() {
		case "player not found", "target version not found in history":
			RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		default:
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert player", nil)
		}
		return
	}

	RespondSuccess(c, nil, "Player reverted successfully")
}
