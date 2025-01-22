package app

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/one111eric/blogger-backend/db"
	"github.com/one111eric/blogger-backend/models"
	"github.com/rs/zerolog/log"
)

func GetPosts(w http.ResponseWriter, r *http.Request, database *db.Database) {
	// Extract or generate traceId
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}

	log.Info().
		Str("traceId", traceId).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Handling GetPosts request")

	// Return empty array if no posts found
	posts := []models.Post{}
	results, err := database.GetPosts()
	if err != nil {
		log.Error().
			Str("traceId", traceId).
			Err(err).
			Msg("Error fetching posts from database")

		errorResponse := models.ErrorResponse{
			Error:   err.Error(),
			TraceId: traceId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if results != nil {
		posts = results
	}

	log.Info().
		Str("traceId", traceId).
		Int("postCount", len(posts)).
		Msg("Successfully fetched posts")

	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Items:   posts,
		TraceId: traceId,
	}
	json.NewEncoder(w).Encode(response)
}

func CreatePost(w http.ResponseWriter, r *http.Request, database *db.Database) {
	// Extract or generate traceId
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}

	log.Info().
		Str("traceId", traceId).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Handling CreatePost request")

	// Decode the request body into a Post object
	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		log.Error().
			Str("traceId", traceId).
			Err(err).
			Msg("Invalid request body for CreatePost")

		errorResponse := models.ErrorResponse{
			Error:   "Invalid request body",
			TraceId: traceId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Call the database layer to create a new post
	id, err := database.CreatePost(&newPost)
	if err != nil {
		log.Error().
			Str("traceId", traceId).
			Err(err).
			Msg("Error creating post in database")

		errorResponse := models.ErrorResponse{
			Error:   err.Error(),
			TraceId: traceId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Assign the generated ID to the new post
	newPost.ID = id

	log.Info().
		Str("traceId", traceId).
		Int("postId", id).
		Msg("Successfully created new post")

	// Respond with the newly created post and the traceId
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Items:   newPost,
		TraceId: traceId,
	}
	json.NewEncoder(w).Encode(response)
}
