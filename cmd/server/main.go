package main

import (
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
func main() {
	// Load config
	cfg := config.Load()

	// Setup logger
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

	// Setup router
	r := router.Setup(db)

	// Start server
	log.Info().Str("port", cfg.ServerPort).Msg("Starting server")
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
