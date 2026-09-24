package handlers

import (
	"context"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"log/slog"
)

type DirectorHandler struct {
	Queries *sqlc.Queries
}

func (h *DirectorHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Debug("failed to parse Get request.")
	}

	user, err := h.Queries.GetUser(context.Background(), id) 
	w.Write([]byte(user.Username))
}

func (h *DirectorHandler) Post(w http.ResponseWriter, r *http.Request) {
}
