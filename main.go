package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/one111eric/blogger-backend/app"
	"github.com/one111eric/blogger-backend/db"
	"github.com/one111eric/blogger-backend/logger"
	"github.com/rs/cors"
)

func main() {
	// Init Logger
	logger.InitializeLogger()

	// Log the application startup
	logger.Info("Application is starting", map[string]interface{}{
		"port": 8080,
	})

	// Initialize the database
	database, err := db.Initialize("./blog.db")
	if err != nil {
		logger.Error("Error initializing database", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("Error closing database", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	router := mux.NewRouter()

	// Register v1 routes
	app.V1Routes(router, database)
	// Set up CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Trace-Id"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	logger.Info("Server listening", map[string]interface{}{
		"port": 8080,
	})
	//fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", corsHandler)
}
