package handlers

import (
	"crypto/sha256"
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	rc, err := r.Cookie("refresh_token")
	if err != nil {
		// http.Error(w, "missing refresh token", http.StatusUnauthorized)
		if r.Header.Get("HX-Request") != "" {
			w.Header().Add("HX-Redirect", "/login")
		}
		return
	}

	if ac, err := r.Cookie("access_token"); err == nil {
		if _, err := h.Tokens.Validate(ac.Value); err == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	now := time.Now()

	refTo := auth.CreateRef()
	nHash := sha256.Sum256([]byte(refTo))
	oHash := sha256.Sum256([]byte(rc.Value))

	row, err := h.Queries.RotateToken(r.Context(), sqlc.RotateTokenParams{
		NewHash: nHash[:],
		OldHash: oHash[:],
		Now:     now.Unix(),
	})
	if err != nil {
		// http.Error(w, "missing refresh token", http.StatusUnauthorized)
		if r.Header.Get("HX-Request") != "" {
			w.Header().Add("HX-Redirect", "/login")
		}
		return
	}
	setAuthCookie(w, "refresh_token", refTo, time.Unix(row.ExpiresAt, 0), "/api/auth")

	user, err := h.Queries.GetUserByUserID(r.Context(), row.UserID)
	accessToken, err := h.Tokens.IssueAcc(
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)
	setAuthCookie(w, "access_token", accessToken, now.Add(h.Tokens.AccTTL()), "/")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	dbUser, _ := h.Queries.GetUserByUsername(r.Context(), r.PostFormValue("username"))
	if bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(r.PostFormValue("password"))) != nil {
		http.Error(w, "incorrect credentials", http.StatusUnauthorized)
	}

	now := time.Now()

	accTo, err := h.Tokens.IssueAcc(strconv.FormatInt(dbUser.UserID, 10), dbUser.Username, dbUser.Role)
	if err != nil {
		http.Error(w, "failed to issue access token.", http.StatusInternalServerError)
		return
	}
	setAuthCookie(w, "access_token", accTo, now.Add(h.Tokens.AccTTL()), "/")

	refTo := auth.CreateRef()
	hash := sha256.Sum256([]byte(refTo))
	res, err := h.Queries.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
		TokenHash: hash[:],
		UserID:    dbUser.UserID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(h.Tokens.RefTTL()).Unix(),
	})
	setAuthCookie(w, "refresh_token", refTo, time.Unix(res.ExpiresAt, 0), "/api/auth")

	if r.Header.Get("HX-Request") != "" {
		w.Header().Add("HX-Redirect", "/dashboard")
	}
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

	now := time.Now()

	accessToken, err := h.Tokens.IssueAcc(
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)
	if err != nil {
		http.Error(w, "failed to issue access token.", http.StatusInternalServerError)
		return
	}
	setAuthCookie(w, "access_token", accessToken, now.Add(h.Tokens.AccTTL()), "/")

	refTo := auth.CreateRef()
	hash := sha256.Sum256([]byte(refTo))
	res, err := h.Queries.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
		TokenHash: hash[:],
		UserID:    user.UserID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(h.Tokens.RefTTL()).Unix(),
	})
	setAuthCookie(w, "refresh_token", refTo, time.Unix(res.ExpiresAt, 0), "/api/auth")

	if r.Header.Get("HX-Request") != "" {
		w.Header().Add("HX-Redirect", "/dashboard")
	}
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Logout api request.")

	setAuthCookie(w, "access_token", "", time.Now(), "/")

	if cookie, err := r.Cookie("refresh_token"); err == nil {
		hash := sha256.Sum256([]byte(cookie.Value))
		h.Queries.DeleteRefreshToken(r.Context(), hash[:])
		setAuthCookie(w, "refresh_token", "", time.Now(), "/api/auth")
	}

	if r.Header.Get("HX-Request") != "" {
		w.Header().Add("HX-Redirect", "/login")
	}
}

func setAuthCookie(w http.ResponseWriter, name, value string, expires time.Time, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  expires,
	})
}
