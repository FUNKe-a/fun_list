package handlers

import (
	"context"
	"net/http"
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

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//http.Error(w, "missing access token", http.StatusUnauthorized)
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
