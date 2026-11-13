package handlers

import (
	"database/sql"
	"net/http"
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"time"
)

const (
	tmplPages = "internal/assets/html/pages"

	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
	refreshTTL        = 30 * 24 * time.Hour
)

type Handler struct {
	DB      *sql.DB
	Queries *sqlc.Queries
	Tokens  *auth.TokenService
}

func setAuthCookie(w http.ResponseWriter, name, value string, ttl time.Duration, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   int(ttl.Seconds()),
		Expires:  time.Now().Add(ttl),
	})
}

