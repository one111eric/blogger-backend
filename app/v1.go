package app

import (
	"net/http"

	"github.com/one111eric/blogger-backend/db"
	"github.com/one111eric/blogger-backend/logger"
)

// v1Routes registers all v1 routes to the provided ServeMux
func V1Routes(mux *http.ServeMux, database *db.Database) {
	logger.Info("Registering v1 routes", map[string]interface{}{
		"version": "v1",
	})

	// Route for getting posts
	mux.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received", map[string]interface{}{
			"method": r.Method,
			"path":   "/v1/posts",
		})

		if r.Method == http.MethodGet {
			GetPosts(w, r, database)
		} else if r.Method == http.MethodPost {
			CreatePost(w, r, database)
		} else {
			logger.Error("Unsupported HTTP method", map[string]interface{}{
				"method": r.Method,
			})

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
		}
	})

	logger.Info("Routes successfully registered", map[string]interface{}{
		"version": "v1",
	})
}
