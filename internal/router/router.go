package router

import (
	"kickbase/internal/config"
	"kickbase/internal/domain"
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
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewResultRepository(db)
	goalRepo := repository.NewGoalRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
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

	// Auth Guard Middleware
	authGuard := middleware.AuthMiddleware(cfg.JWTSecret)

	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			handler.RespondSuccess(c, gin.H{"status": "ok", "database": "connected"}, "Service is healthy")
		})

		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Teams
		teams := api.Group("/teams")
		{
			// Read operations (Public)
			teams.GET("", teamHandler.ListTeams)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.GET("/:id/history", teamHandler.GetTeamHistory)

			// Permission-protected mutations
			teams.POST("", authGuard, middleware.RequirePermission(domain.PermTeamsCreate), teamHandler.CreateTeam)
			teams.PUT("/:id", authGuard, middleware.RequirePermission(domain.PermTeamsUpdate), teamHandler.UpdateTeam)
			teams.DELETE("/:id", authGuard, middleware.RequirePermission(domain.PermTeamsDelete), teamHandler.DeleteTeam)
			teams.POST("/:id/revert", authGuard, middleware.RequirePermission(domain.PermTeamsRevert), teamHandler.RevertTeam)
		}

		// Players
		players := api.Group("/players")
		{
			// Read operations (Public)
			players.GET("", playerHandler.ListPlayers)
			players.GET("/:id", playerHandler.GetPlayer)
			players.GET("/:id/history", playerHandler.GetPlayerHistory)

			// Permission-protected mutations
			players.POST("", authGuard, middleware.RequirePermission(domain.PermPlayersCreate), playerHandler.CreatePlayer)
			players.PUT("/:id", authGuard, middleware.RequirePermission(domain.PermPlayersUpdate), playerHandler.UpdatePlayer)
			players.DELETE("/:id", authGuard, middleware.RequirePermission(domain.PermPlayersDelete), playerHandler.DeletePlayer)
			players.POST("/:id/revert", authGuard, middleware.RequirePermission(domain.PermPlayersRevert), playerHandler.RevertPlayer)
		}

		// Matches
		matches := api.Group("/matches")
		{
			// Read operations (Public)
			matches.GET("", matchHandler.ListMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.GET("/:id/history", matchHandler.GetMatchHistory)

			// Permission-protected mutations
			matches.POST("", authGuard, middleware.RequirePermission(domain.PermMatchesCreate), matchHandler.CreateMatch)
			matches.PATCH("/:id/status", authGuard, middleware.RequirePermission(domain.PermMatchesUpdate), matchHandler.UpdateMatchStatus)
			matches.POST("/:id/revert", authGuard, middleware.RequirePermission(domain.PermMatchesRevert), matchHandler.RevertMatch)
		}

		// Results
		results := api.Group("/results")
		{
			// Read operations (Public)
			results.GET("/:matchId", resultHandler.GetMatchResult)

			// Permission-protected mutations
			results.POST("", authGuard, middleware.RequirePermission(domain.PermResultsCreate), resultHandler.CreateMatchResult)
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
