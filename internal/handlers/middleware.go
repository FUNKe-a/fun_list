package handlers

import (
	"net/http"
)

func (h *Handler) ValidateToken(ogHa http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("access_token")
		if err := cookie.Valid(); err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		if _, err := h.Tokens.Validate(cookie.Value); err == nil {
			ogHa.ServeHTTP(w, r)
		}
	})
}
