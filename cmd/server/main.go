package main

import (
	"io"
	"os"

	_ "kickbase/docs"
	"kickbase/internal/config"
	"kickbase/internal/database"
	"kickbase/internal/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Kickbase Football Team Management API
// @version 1.0
// @description REST API for managing amateur football teams, players, matches, and reports for Perusahaan XYZ.
// @host localhost:8080
// @BasePath /
// @schemes http
// @produce json
// @consumes json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load config
	cfg := config.Load()

	// Setup multi-writer logger (Stdout + File for Promtail/Loki)
	logDir := "/var/log/kickbase"
	var writers []io.Writer = []io.Writer{os.Stdout}

	if err := os.MkdirAll(logDir, 0755); err == nil {
		if logFile, err := os.OpenFile(logDir+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			writers = append(writers, logFile)
		}
	}

	multiWriter := io.MultiWriter(writers...)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch cfg.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Seed initial data (Admin/Staff/Viewer & Initial Club Data)
	if err := database.Seed(db); err != nil {
		log.Warn().Err(err).Msg("Database seeding warning")
	}

	// Setup router
	r := router.Setup(db)

	// Start server
	log.Info().Str("port", cfg.ServerPort).Msg("Starting server")
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
