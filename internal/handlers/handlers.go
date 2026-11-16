package handlers

import (
	"database/sql"
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
)

const (
	refreshCookieName = "refresh_token"
)

type Handler struct {
	DB      *sql.DB
	Queries *sqlc.Queries
	Tokens  *auth.AuthService
}


