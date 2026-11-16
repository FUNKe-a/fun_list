package handlers

import (
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	dbUser, _ := h.Queries.GetUserByUsername(r.Context(), r.PostFormValue("username"))
	if bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(r.PostFormValue("password"))) != nil {
		http.Error(w, "incorrect credentials", http.StatusUnauthorized)
	}

	accessToken, err := h.Tokens.Issue(strconv.FormatInt(dbUser.UserID, 10), dbUser.Username, dbUser.Role)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, "access_token", accessToken, h.Tokens.TTL(), "/")

	//w.Header().Add("HX-Redirect", "/dashboard")
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
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	setAuthCookie(w, "access_token", "", 0, "/")
}

func setAuthCookie(w http.ResponseWriter, name, value string, ttl time.Duration, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   int(ttl.Seconds()),
		Expires:  time.Now().Add(ttl),
	})
}
