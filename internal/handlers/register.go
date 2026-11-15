package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(tmplPages + "/register.tmpl"))
	tmpl.Execute(w, map[string]any{"Title": "Home Page"})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")

	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(r.PostFormValue("password")),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	params := sqlc.CreateUserParams{
		Username:     r.PostFormValue("username"),
		Email:        &email,
		PasswordHash: string(passHash),
	}

	user, err := h.Queries.CreateUser(r.Context(), params)
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.Tokens.Issue(
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, "access_token", accessToken, h.Tokens.TTL(), "/")

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
