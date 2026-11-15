package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"log/slog"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"strings"
)

// Comment table handlers
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	comment, err := h.Queries.GetComment(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Queries.ListComments(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var params sqlc.CreateCommentParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var errs []string
	if params.UserID <= 0 {
		errs = append(errs, "Field 'user_id' cannot be empty and must be above 0.")
	}
	if params.MediaID <= 0 {
		errs = append(errs, "Field 'media_id' cannot be empty and must be above 0.")
	}
	if len(params.Content) == 0 {
		errs = append(errs, "Field 'content' is required and cannot be empty.")
	}
	if params.Rating == 0 {
		errs = append(errs, "Field 'rating' cannot be empty and must be above 0.")
	}
	if len(errs) > 0 {
		errorMessage := strings.Join(errs, "\n")
		http.Error(w, errorMessage, http.StatusUnprocessableEntity)
		return
	}

	comment, err := h.Queries.CreateComment(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/api/comments/%d", comment.CommentID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteComment(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var params sqlc.UpdateCommentParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params.CommentID = id

	var errs []string
	if params.UserID <= 0 {
		errs = append(errs, "Field 'user_id' cannot be empty and must be above 0.")
	}
	if params.MediaID <= 0 {
		errs = append(errs, "Field 'media_id' cannot be empty and must be above 0.")
	}
	if len(params.Content) == 0 {
		errs = append(errs, "Field 'content' is required and cannot be empty.")
	}
	if params.Rating == 0 {
		errs = append(errs, "Field 'rating' cannot be empty and must be above 0.")
	}
	if len(errs) > 0 {
		errorMessage := strings.Join(errs, "\n")
		http.Error(w, errorMessage, http.StatusUnprocessableEntity)
		return
	}

	comment, err := h.Queries.UpdateComment(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to update comment", "error", err)
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// Director table handlers
func (h *Handler) GetDirector(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	director, err := h.Queries.GetDirector(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Director not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(director)
}

func (h *Handler) ListDirectors(w http.ResponseWriter, r *http.Request) {
	directors, err := h.Queries.ListDirectors(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directors)
}

func (h *Handler) PostDirector(w http.ResponseWriter, r *http.Request) {
	var params sqlc.CreateDirectorParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(params.Name) == 0 {
		http.Error(w, "Field 'name' is required and cannot be empty", http.StatusUnprocessableEntity)
		return
	}

	director, err := h.Queries.CreateDirector(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to create director", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/api/directors/%d", director.DirectorID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(director)
}

func (h *Handler) DeleteDirector(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteDirector(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete director", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutDirector(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var params sqlc.UpdateDirectorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params.DirectorID = id

	if len(params.Name) == 0 {
		http.Error(w, "Field 'name' is required and cannot be empty", http.StatusUnprocessableEntity)
		return
	}

	director, err := h.Queries.UpdateDirector(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Director not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to update director", "error", err)
		http.Error(w, "Failed to update director", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(director)
}

// Media table handlers
func (h *Handler) GetMedia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	media, err := h.Queries.GetMedia(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Media not found", http.StatusNotFound)
			return
		}
		fmt.Println(media)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *Handler) ListMedia(w http.ResponseWriter, r *http.Request) {
	media, err := h.Queries.ListMedia(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *Handler) PostMedia(w http.ResponseWriter, r *http.Request) {
	var params sqlc.CreateMediaParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var errs []string
	if len(params.Title) == 0 {
		errs = append(errs, "Field 'title' is required and cannot be empty.")
	}
	if params.DirectorID <= 0 {
		errs = append(errs, "Field 'director_id' cannot be empty and must be above 0.")
	}
	if len(errs) > 0 {
		http.Error(w, strings.Join(errs, "\n"), http.StatusUnprocessableEntity)
		return
	}

	media, err := h.Queries.CreateMedia(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to create media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/api/media/%d", media.MediaID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}

func (h *Handler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteMedia(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete media", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutMedia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var params sqlc.UpdateMediaParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params.MediaID = id

	var errs []string
	if len(params.Title) == 0 {
		errs = append(errs, "Field 'title' is required and cannot be empty.")
	}
	if params.DirectorID <= 0 {
		errs = append(errs, "Field 'director_id' cannot be empty and must be above 0.")
	}
	if len(errs) > 0 {
		http.Error(w, strings.Join(errs, "\n"), http.StatusUnprocessableEntity)
		return
	}

	media, err := h.Queries.UpdateMedia(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Media not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to update media", "error", err)
		http.Error(w, "Failed to update media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}
