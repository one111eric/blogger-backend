package app

import (
	"net/http"

	"github.com/one111eric/blogger-backend/db"
)

func V1Routes(router *http.ServeMux, database *db.Database) {
	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetPosts(w, r, database)
		} else if r.Method == http.MethodPost {
			CreatePost(w, r, database)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}
