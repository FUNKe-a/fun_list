package routes

import (
	"net/http"
	"github.com/FUNKe-a/fun_list/internal/handlers"
)

func CreateMux(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", viewMux(h))
	mux.Handle("/api/", http.StripPrefix("/api", apiMux(h)))

	// static files
	sf := http.FileServer(http.Dir("internal/assets/static"))
	mux.Handle("/static/",  http.StripPrefix("/static/", sf))

	return mux
}

func apiMux(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// db api
	mux.HandleFunc("GET /directors", h.ListDirectors)
	mux.HandleFunc("GET /directors/{id}", h.GetDirector)
	mux.HandleFunc("POST /directors", h.PostDirector)
	mux.HandleFunc("PUT /directors/{id}", h.PutDirector)
	mux.HandleFunc("DELETE /directors/{id}", h.DeleteDirector)

	mux.HandleFunc("GET /media", h.ListMedia)
	mux.HandleFunc("GET /media/{id}", h.GetMedia)
	mux.HandleFunc("POST /media", h.PostMedia)
	mux.HandleFunc("PUT /media/{id}", h.PutMedia)
	mux.HandleFunc("DELETE /media/{id}", h.DeleteMedia)

	mux.HandleFunc("GET /comments", h.ListComments)
	mux.HandleFunc("GET /comments/{id}", h.GetComment)
	mux.HandleFunc("POST /comments", h.PostComment)
	mux.HandleFunc("PUT /comments/{id}", h.PutComment)
	mux.HandleFunc("DELETE /comments/{id}", h.DeleteComment)

	//  login api
	mux.HandleFunc("POST /auth/login", h.HandleLogin)
	mux.HandleFunc("POST /auth/register", h.HandleRegister)
	mux.HandleFunc("POST /auth/logout", h.HandleLogout)

	return mux
}

func viewMux(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", h.ShowLogin)
	mux.HandleFunc("GET /register", h.ShowRegister)
	mux.Handle("GET /dashboard", h.ValidateToken(http.HandlerFunc(h.ShowDashboard)))

	return mux
}
