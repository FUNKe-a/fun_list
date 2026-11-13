package handlers

import (
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"time"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
)

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(tmplPages + "/login.tmpl"))
	tmpl.Execute(w, map[string]any{"Title": "Home Page"})
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	dbUser, _ := h.Queries.GetUserByUsername(r.Context(), r.PostFormValue("username"))
	if bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(r.PostFormValue("password"))) != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	}

	accessToken, err := h.Tokens.Issue(strconv.FormatInt(dbUser.UserID, 10), dbUser.Username, dbUser.Role)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	_, err = h.Queries.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    dbUser.UserID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(refreshTTL).Unix(),
	})
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, "access_token", accessToken, h.Tokens.TTL(), "/")
	setAuthCookie(w, "refresh_token", refreshToken, refreshTTL, "/auth/refresh")

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
