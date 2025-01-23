package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/one111eric/blogger-backend/db"
	"github.com/one111eric/blogger-backend/logger"
	"github.com/one111eric/blogger-backend/models"
)

func GetPosts(w http.ResponseWriter, r *http.Request, database *db.Database) {
	// Extract or generate traceId
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}

	logger.Info("Handling GetPosts request", map[string]interface{}{
		"traceId": traceId,
		"method":  r.Method,
		"path":    r.URL.Path,
	})

	// Return empty array if no posts found
	posts := []models.Post{}
	results, err := database.GetPosts()
	if err != nil {
		logger.Error("Failed to get posts", map[string]interface{}{
			"traceId": traceId,
			"error":   err.Error(),
		})

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

	logger.Info("Successfully retrieved posts", map[string]interface{}{
		"traceId": traceId,
		"count":   len(posts),
	})

	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Items:   posts,
		TraceId: traceId,
	}
	json.NewEncoder(w).Encode(response)
}

func GetPostByIdHandler(w http.ResponseWriter, r *http.Request, database *db.Database, idStr string) {
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}

	// Extract the ID from the request
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := database.GetPostById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post":    post,
		"traceId": traceId,
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request, database *db.Database) {
	// Extract or generate traceId
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}

	logger.Info("Handling CreatePost request", map[string]interface{}{
		"traceId": traceId,
		"method":  r.Method,
		"path":    r.URL.Path,
	})

	// Decode the request body into a Post object
	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		logger.Error("Invalid request body for CreatePost", map[string]interface{}{
			"traceId": traceId,
			"error":   err.Error(),
		})

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
		logger.Error("Error creating post in database", map[string]interface{}{
			"traceId": traceId,
			"error":   err.Error(),
		})

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

	logger.Info("Successfully created new post", map[string]interface{}{
		"traceId": traceId,
		"postId":  id,
	})

	// Respond with the newly created post and the traceId
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Items:   newPost,
		TraceId: traceId,
	}
	json.NewEncoder(w).Encode(response)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request, database *db.Database, idStr string) {
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {

		logger.Error("Invalid post ID", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		errorResponse := models.ErrorResponse{
			Error:   err.Error(),
			TraceId: traceId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rowsAffected, err := database.EditPost(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rowsAffected": rowsAffected,
		"traceId":      traceId,
	})
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request, database *db.Database, idStr string) {
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {

		logger.Error("Invalid post ID", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		errorResponse := models.ErrorResponse{
			Error:   err.Error(),
			TraceId: traceId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Attempt to delete the post
	rowsDeleted, err := database.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rowsDeleted": rowsDeleted,
		"traceId":     traceId,
	})
}
