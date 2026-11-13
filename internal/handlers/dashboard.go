package handlers

import (
	"github.com/FUNKe-a/fun_list/internal/auth"
	"html/template"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(tmplPages + "/dashboard.tmpl"))
	tmpl.Execute(w, []string{})
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		tokenHash := auth.HashToken(cookie.Value)

		_, err := h.Queries.RevokeRefreshToken(r.Context(), tokenHash)
		if err != nil {
			http.Error(w, "Logout failed", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
