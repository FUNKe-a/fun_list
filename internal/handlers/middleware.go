package handlers

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"net/http"
	"strconv"
	"time"
)

type claimsKey struct{}

func (h *Handler) ValidateToken(ogHa http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("access_token"); err == nil {
			if claims, err := h.Tokens.Validate(cookie.Value); err == nil {
				ctx := context.WithValue(r.Context(), claimsKey{}, claims)
				newReq := r.WithContext(ctx)
				ogHa.ServeHTTP(w, newReq)
				return
			}
		}

		//http.Redirect(w, r, "/api/auth/refresh", http.StatusTemporaryRedirect)
		http.Error(w, "missing access token", http.StatusUnauthorized)
	})
}

func (h *Handler) ValidateRole(role string, ogHa http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("access_token"); err == nil {
			if claims, err := h.Tokens.Validate(cookie.Value); err == nil {
				if claims.Role == role {
					ctx := context.WithValue(r.Context(), claimsKey{}, claims)
					newReq := r.WithContext(ctx)
					ogHa.ServeHTTP(w, newReq)
				} else {
					http.Error(w, "Incorrect role", http.StatusUnauthorized)
				}
				return
			}
		}

		//http.Error(w, "missing access token", http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})
}

func (h *Handler) refreshToken(reT string, now time.Time) (string, string, error) {
	if reT == "" {
		return "", "", errors.New("missing refresh token")
	}

	refTo := auth.CreateRef()
	nHash := sha256.Sum256([]byte(refTo))
	oHash := sha256.Sum256([]byte(reT))

	row, err := h.Queries.RotateToken(context.Background(), sqlc.RotateTokenParams{
		NewHash: nHash[:],
		OldHash: oHash[:],
		Now:     now.Unix(),
	})
	if err != nil {
		return "", "", errors.New("failed to rotate token")
	}

	user, err := h.Queries.GetUserByUserID(context.Background(), row.UserID)
	accessToken, err := h.Tokens.IssueAcc(
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)

	return accessToken, refTo, nil
}
