package handlers

import (
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

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
