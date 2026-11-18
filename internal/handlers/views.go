package handlers

import (
	"html/template"
	"net/http"
)

var tmplPages = template.Must(template.ParseGlob("internal/assets/html/pages/*.tmpl"))

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmplPages.ExecuteTemplate(w, "register.tmpl", map[string]any{"Title": "Home Page"})
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmplPages.ExecuteTemplate(w, "login.tmpl", map[string]any{"Title": "Home Page"})
}

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	tmplPages.ExecuteTemplate(w, "dashboard.tmpl", map[string]any{"Title": "Home Page"})
}

func (h *Handler) ShowComments(w http.ResponseWriter, r *http.Request) {
	tmplPages.ExecuteTemplate(w, "comments.tmpl", []string{})
}
