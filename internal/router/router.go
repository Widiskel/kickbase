package router

import (
	"kickbase/internal/handler"
	"kickbase/internal/middleware"
	"kickbase/internal/repository"
	"kickbase/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.Logger())
	r.Use(middleware.PrometheusMetrics())

	// Metrics endpoint
	r.GET("/metrics", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, middleware.GetMetrics())
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize repositories
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewResultRepository(db)
	goalRepo := repository.NewGoalRepository(db)

	// Initialize services
	teamSvc := service.NewTeamService(teamRepo)
	playerSvc := service.NewPlayerService(playerRepo, teamRepo)
	matchSvc := service.NewMatchService(matchRepo, teamRepo)
	resultSvc := service.NewResultService(resultRepo, matchRepo, goalRepo, playerRepo, db)
	reportSvc := service.NewReportService(db, matchRepo, resultRepo, goalRepo, teamRepo, playerRepo)

	// Initialize handlers
	teamHandler := handler.NewTeamHandler(teamSvc)
	playerHandler := handler.NewPlayerHandler(playerSvc)
	matchHandler := handler.NewMatchHandler(matchSvc)
	resultHandler := handler.NewResultHandler(resultSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

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
			teams.POST("", teamHandler.CreateTeam)
			teams.GET("", teamHandler.ListTeams)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.PUT("/:id", teamHandler.UpdateTeam)
			teams.DELETE("/:id", teamHandler.DeleteTeam)
			teams.GET("/:id/history", teamHandler.GetTeamHistory)
			teams.POST("/:id/revert", teamHandler.RevertTeam)
		}

		// Players
		players := api.Group("/players")
		{
			players.POST("", playerHandler.CreatePlayer)
			players.GET("", playerHandler.ListPlayers)
			players.GET("/:id", playerHandler.GetPlayer)
			players.PUT("/:id", playerHandler.UpdatePlayer)
			players.DELETE("/:id", playerHandler.DeletePlayer)
			players.GET("/:id/history", playerHandler.GetPlayerHistory)
			players.POST("/:id/revert", playerHandler.RevertPlayer)
		}

		// Matches
		matches := api.Group("/matches")
		{
			matches.POST("", matchHandler.CreateMatch)
			matches.GET("", matchHandler.ListMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.PATCH("/:id/status", matchHandler.UpdateMatchStatus)
			matches.GET("/:id/history", matchHandler.GetMatchHistory)
			matches.POST("/:id/revert", matchHandler.RevertMatch)
		}

		// Results
		results := api.Group("/results")
		{
			results.POST("", resultHandler.CreateMatchResult)
			results.GET("/:matchId", resultHandler.GetMatchResult)
		}

		// Reports
		reports := api.Group("/reports")
		{
			reports.GET("/matches", reportHandler.ListMatchReports)
			reports.GET("/matches/:id", reportHandler.GetMatchReport)
		}
	}

	return r
}
