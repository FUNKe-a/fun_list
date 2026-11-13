
package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	oldHash := auth.HashToken(cookie.Value)

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := h.Queries.WithTx(tx)

	oldToken, err := qtx.GetActiveRefreshToken(r.Context(), oldHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := qtx.RotateRefreshToken(r.Context(), oldHash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rows != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := qtx.GetUsers(r.Context(), oldToken.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.Tokens.Issue(
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = qtx.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    user.UserID,
		TokenHash: auth.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(refreshTTL).Unix(),
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, accessCookieName, accessToken, h.Tokens.TTL(), "/")
	setAuthCookie(w, refreshCookieName, newRefreshToken, refreshTTL, "/auth/refresh")

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

