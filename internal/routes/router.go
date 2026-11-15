package routes

import (
	"net/http"
	"github.com/FUNKe-a/fun_list/internal/handlers"
)

func RoutesHandle(mux *http.ServeMux, h *handlers.Handler) {

	// api
	mux.HandleFunc("GET /api/directors", h.ListDirectors)
	mux.HandleFunc("GET /api/directors/{id}", h.GetDirector)
	mux.HandleFunc("POST /api/directors", h.PostDirector)
	mux.HandleFunc("PUT /api/directors/{id}", h.PutDirector)
	mux.HandleFunc("DELETE /api/directors/{id}", h.DeleteDirector)

	mux.HandleFunc("GET /api/media", h.ListMedia)
	mux.HandleFunc("GET /api/media/{id}", h.GetMedia)
	mux.HandleFunc("POST /api/media", h.PostMedia)
	mux.HandleFunc("PUT /api/media/{id}", h.PutMedia)
	mux.HandleFunc("DELETE /api/media/{id}", h.DeleteMedia)

	mux.HandleFunc("GET /api/comments", h.ListComments)
	mux.HandleFunc("GET /api/comments/{id}", h.GetComment)
	mux.HandleFunc("POST /api/comments", h.PostComment)
	mux.HandleFunc("PUT /api/comments/{id}", h.PutComment)
	mux.HandleFunc("DELETE /api/comments/{id}", h.DeleteComment)

	mux.HandleFunc("GET /auth/login", h.ShowLogin)
	mux.HandleFunc("POST /auth/login", h.HandleLogin)
	mux.HandleFunc("GET /auth/register", h.ShowRegister)
	mux.HandleFunc("POST /auth/register", h.HandleRegister)
	mux.HandleFunc("POST /auth/logout", h.HandleLogout)

	mux.HandleFunc("GET /dashboard", h.ShowDashboard)
}
