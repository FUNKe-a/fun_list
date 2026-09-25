package handlers

import (
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"log/slog"
	"encoding/json"
	"fmt"
	"database/sql"
	"strings"
)

type CommentHandler struct {
	Queries *sqlc.Queries
}

func (h *CommentHandler) Get(w http.ResponseWriter, r *http.Request) {
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



func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
    comments, err := h.Queries.ListComments(r.Context())
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) Post(w http.ResponseWriter, r *http.Request) {
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

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *CommentHandler) Patch(w http.ResponseWriter, r *http.Request) {
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
