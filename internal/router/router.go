package router

import (
	"kickbase/internal/config"
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
	cfg := config.Load()
	r := gin.Default()

	// Global Middleware
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
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewResultRepository(db)
	goalRepo := repository.NewGoalRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	teamSvc := service.NewTeamService(teamRepo)
	playerSvc := service.NewPlayerService(playerRepo, teamRepo)
	matchSvc := service.NewMatchService(matchRepo, teamRepo)
	resultSvc := service.NewResultService(resultRepo, matchRepo, goalRepo, playerRepo, db)
	reportSvc := service.NewReportService(db, matchRepo, resultRepo, goalRepo, teamRepo, playerRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	playerHandler := handler.NewPlayerHandler(playerSvc)
	matchHandler := handler.NewMatchHandler(matchSvc)
	resultHandler := handler.NewResultHandler(resultSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

	// Auth guards
	authGuard := middleware.AuthMiddleware(cfg.JWTSecret)
	adminGuard := middleware.RequireRole("admin")

	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			handler.RespondSuccess(c, gin.H{"status": "ok", "database": "connected"}, "Service is healthy")
		})

		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Teams
		teams := api.Group("/teams")
		{
			// Read operations (Public)
			teams.GET("", teamHandler.ListTeams)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.GET("/:id/history", teamHandler.GetTeamHistory)

			// Admin operations (Protected)
			teams.POST("", authGuard, adminGuard, teamHandler.CreateTeam)
			teams.PUT("/:id", authGuard, adminGuard, teamHandler.UpdateTeam)
			teams.DELETE("/:id", authGuard, adminGuard, teamHandler.DeleteTeam)
			teams.POST("/:id/revert", authGuard, adminGuard, teamHandler.RevertTeam)
		}

		// Players
		players := api.Group("/players")
		{
			// Read operations (Public)
			players.GET("", playerHandler.ListPlayers)
			players.GET("/:id", playerHandler.GetPlayer)
			players.GET("/:id/history", playerHandler.GetPlayerHistory)

			// Admin operations (Protected)
			players.POST("", authGuard, adminGuard, playerHandler.CreatePlayer)
			players.PUT("/:id", authGuard, adminGuard, playerHandler.UpdatePlayer)
			players.DELETE("/:id", authGuard, adminGuard, playerHandler.DeletePlayer)
			players.POST("/:id/revert", authGuard, adminGuard, playerHandler.RevertPlayer)
		}

		// Matches
		matches := api.Group("/matches")
		{
			// Read operations (Public)
			matches.GET("", matchHandler.ListMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.GET("/:id/history", matchHandler.GetMatchHistory)

			// Admin operations (Protected)
			matches.POST("", authGuard, adminGuard, matchHandler.CreateMatch)
			matches.PATCH("/:id/status", authGuard, adminGuard, matchHandler.UpdateMatchStatus)
			matches.POST("/:id/revert", authGuard, adminGuard, matchHandler.RevertMatch)
		}

		// Results
		results := api.Group("/results")
		{
			// Read operations (Public)
			results.GET("/:matchId", resultHandler.GetMatchResult)

			// Admin operations (Protected)
			results.POST("", authGuard, adminGuard, resultHandler.CreateMatchResult)
		}

		// Reports (Public read-only aggregations)
		reports := api.Group("/reports")
		{
			reports.GET("/matches", reportHandler.ListMatchReports)
			reports.GET("/matches/:id", reportHandler.GetMatchReport)
		}
	}

	return r
}
