package router

import (
	"kickbase/internal/handler"
	"kickbase/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.Logger())

	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			handler.RespondSuccess(c, gin.H{"status": "ok", "database": "connected"}, "Service is healthy")
		})

		// Teams
		teams := api.Group("/teams")
		{
			teams.POST("", handler.CreateTeam(db))
			teams.GET("", handler.ListTeams(db))
			teams.GET("/:id", handler.GetTeam(db))
			teams.PUT("/:id", handler.UpdateTeam(db))
			teams.DELETE("/:id", handler.DeleteTeam(db))
			teams.GET("/:id/history", handler.GetTeamHistory(db))
			teams.POST("/:id/revert", handler.RevertTeam(db))
		}

		// Players (separate group to avoid route conflict)
		players := api.Group("/players")
		{
			players.POST("", handler.CreatePlayer(db))
			players.GET("", handler.ListPlayers(db))
			players.GET("/:id", handler.GetPlayer(db))
			players.PUT("/:id", handler.UpdatePlayer(db))
			players.DELETE("/:id", handler.DeletePlayer(db))
			players.GET("/:id/history", handler.GetPlayerHistory(db))
			players.POST("/:id/revert", handler.RevertPlayer(db))
		}

		// Matches
		matches := api.Group("/matches")
		{
			matches.POST("", handler.CreateMatch(db))
			matches.GET("", handler.ListMatches(db))
			matches.GET("/:id", handler.GetMatch(db))
			matches.PATCH("/:id/status", handler.UpdateMatchStatus(db))
			matches.GET("/:id/history", handler.GetMatchHistory(db))
			matches.POST("/:id/revert", handler.RevertMatch(db))
		}

		// Match Results
		results := api.Group("/results")
		{
			results.POST("", handler.CreateMatchResult(db))
			results.GET("/:matchId", handler.GetMatchResult(db))
		}

		// Reports
		reports := api.Group("/reports")
		{
			reports.GET("/matches", handler.ListMatchReports(db))
			reports.GET("/matches/:id", handler.GetMatchReport(db))
		}
	}

	return r
}
