package handlers

import (
	"html/template"
	"net/http"
)

const (
	tmplPages = "internal/assets/html/pages"
)

var tmpl = template.Must(template.ParseGlob(tmplPages + "/*.tmpl"))

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register.tmpl", map[string]any{"Title": "Home Page"})
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.tmpl", map[string]any{"Title": "Home Page"})
}

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "dashboard.tmpl", map[string]any{"Title": "Home Page"})
}
