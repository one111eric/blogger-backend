package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/one111eric/blogger-backend/db"
	"github.com/one111eric/blogger-backend/logger"
)

// V1Routes registers all v1 routes to the provided router
func V1Routes(router *mux.Router, database *db.Database) {
	logger.Info("Registering v1 routes", map[string]interface{}{
		"version": "v1",
	})

	// Route for managing posts
	router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodGet, http.MethodPost)

	// Routes for individual posts using {id} as a path parameter
	router.HandleFunc("/v1/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		logger.Info("Request received", map[string]interface{}{
			"method": r.Method,
			"path":   "/v1/posts/{id}",
			"id":     id,
		})

		switch r.Method {
		case http.MethodGet:
			GetPostByIdHandler(w, r, database, id)
		case http.MethodPut:
			EditPostHandler(w, r, database, id)
		case http.MethodDelete:
			DeletePostHandler(w, r, database, id)
		default:
			logger.Error("Unsupported HTTP method", map[string]interface{}{
				"method": r.Method,
			})

			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
		}
	}).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	logger.Info("Routes successfully registered", map[string]interface{}{
		"version": "v1",
	})
}
