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

type MediaHandler struct {
	Queries *sqlc.Queries
}

func (h *MediaHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *MediaHandler) List(w http.ResponseWriter, r *http.Request) {
    media, err := h.Queries.ListMedia(r.Context())
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) Post(w http.ResponseWriter, r *http.Request) {
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

func (h *MediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *MediaHandler) Patch(w http.ResponseWriter, r *http.Request) {
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
