package handlers

import (
	"net/http"
)

func (h *Handler) ValidateToken(ogHa http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("access_token"); err == nil {
			if _, err := h.Tokens.Validate(cookie.Value); err == nil {
				ogHa.ServeHTTP(w, r)
				return
			}
		}

		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	})
}
