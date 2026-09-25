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
)

type DirectorHandler struct {
	Queries *sqlc.Queries
}

func (h *DirectorHandler) Get(w http.ResponseWriter, r *http.Request) {
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



func (h *DirectorHandler) List(w http.ResponseWriter, r *http.Request) {
    directors, err := h.Queries.ListDirectors(r.Context())
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(directors)
}

func (h *DirectorHandler) Post(w http.ResponseWriter, r *http.Request) {
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

func (h *DirectorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *DirectorHandler) Patch(w http.ResponseWriter, r *http.Request) {
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
