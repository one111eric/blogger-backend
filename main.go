package main

import (
	"net/http"
	"os"
	"time"

	"github.com/one111eric/blogger-backend/app"
	"github.com/one111eric/blogger-backend/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting server initialization...")

	// Initialize the database
	database, err := db.Initialize("./blog.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Warn().Err(err).Msg("Error closing database")
		}
	}()

	// Create a default HTTP server
	mux := http.NewServeMux()

	// Register v1 routes (pass the database connection)
	app.V1Routes(mux, database)

	// Start the server
	address := ":8080"
	log.Info().Str("address", address).Msg("Server listening")
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
