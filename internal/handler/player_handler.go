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

func CreatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player domain.Player
		if err := c.ShouldBindJSON(&player); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		// Validate position
		validPositions := map[string]bool{
			"CF": true, "SS": true, "LWF": true, "RWF": true,
			"AMF": true, "CMF": true, "DMF": true, "LMF": true, "RMF": true,
			"CB": true, "LB": true, "RB": true, "GK": true,
		}
		if !validPositions[player.Position] {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid position. Valid positions: CF, SS, LWF, RWF, AMF, CMF, DMF, LMF, RMF, CB, LB, RB, GK", nil)
			return
		}

		// Check jersey uniqueness
		var count int64
		db.Model(&domain.Player{}).Where("team_id = ? AND jersey_number = ?", player.TeamID, player.JerseyNumber).Count(&count)
		if count > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Jersey number already exists in this team", nil)
			return
		}

		if err := db.Create(&player).Error; err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create player", nil)
			return
		}
		RespondCreated(c, player)
	}
}

func ListPlayersByTeam(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Param("teamId")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var players []domain.Player
		var total int64

		db.Model(&domain.Player{}).Where("team_id = ?", teamID).Count(&total)
		db.Where("team_id = ?", teamID).Offset((page - 1) * limit).Limit(limit).Find(&players)

		RespondPaginated(c, players, total, page, limit)
	}
}

func ListPlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Query("team_id")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 { page = 1 }
		if limit < 1 { limit = 10 }

		var players []domain.Player
		var total int64

		query := db.Model(&domain.Player{})
		if teamID != "" {
			query = query.Where("team_id = ?", teamID)
		}

		query.Count(&total)
		query.Offset((page - 1) * limit).Limit(limit).Find(&players)

		RespondPaginated(c, players, total, page, limit)
	}
}

func GetPlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}
		RespondSuccess(c, player, "")
	}
}

func UpdatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}

		var input domain.Player
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		if input.Version != player.Version {
			RespondError(c, http.StatusConflict, "CONFLICT", "Version mismatch", nil)
			return
		}

		updates := map[string]interface{}{
			"name":          input.Name,
			"height":        input.Height,
			"weight":        input.Weight,
			"position":      input.Position,
			"playstyle":     input.Playstyle,
			"jersey_number": input.JerseyNumber,
			"version":       player.Version + 1,
		}

		db.Model(&player).Updates(updates)
		RespondSuccess(c, player, "Player updated successfully")
	}
}

func DeletePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player domain.Player
		if err := db.First(&player, "id = ?", id).Error; err != nil {
			RespondError(c, http.StatusNotFound, "NOT_FOUND", "Player not found", nil)
			return
		}

		// Check if player has goal records
		var goalCount int64
		db.Model(&domain.Goal{}).Where("player_id = ?", id).Count(&goalCount)
		if goalCount > 0 {
			RespondError(c, http.StatusConflict, "CONFLICT", "Cannot delete player with goal records", nil)
			return
		}

		db.Delete(&player)
		RespondNoContent(c)
	}
}

func GetPlayerHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var history []domain.PlayerHistory
		db.Where("player_id = ?", id).Order("version").Find(&history)
		RespondSuccess(c, history, "")
	}
}

func RevertPlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input struct {
			TargetVersion int `json:"target_version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
			return
		}

		playerRepo := repository.NewPlayerRepository(db)
		teamRepo := repository.NewTeamRepository(db)
		svc := service.NewPlayerService(playerRepo, teamRepo)

		if err := svc.RevertPlayer(id, input.TargetVersion); err != nil {
			if err.Error() == "player not found" || err.Error() == "target version not found in history" {
				RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			} else {
				RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revert player", nil)
			}
			return
		}

		RespondSuccess(c, nil, "Player reverted successfully")
	}
}

// CreateMatch godoc
// @Summary Schedule a new match
// @Description Create a match schedule between two teams
// @Tags Matches
// @Accept json
// @Produce json
// @Param match body domain.Match true "Match data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
