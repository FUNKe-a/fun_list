package main

import (
	"errors"
	"fmt"
	"github.com/FUNKe-a/fun_list/internal/db"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"github.com/FUNKe-a/fun_list/internal/handlers"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

func main() {
	db_conn, err := db.Init("data/data.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db_conn.Close()

	q := sqlc.New(db_conn)
	mux := http.NewServeMux()
	setHandlers(mux, q)

	http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func setHandlers(mux *http.ServeMux, q *sqlc.Queries) {
	dir_h := handlers.DirectorHandler{
		Queries: q,
	}
	media_h := handlers.MediaHandler{
		Queries: q,
	}
	comment_h := handlers.CommentHandler{
		Queries: q,
	}

	mux.HandleFunc("GET /api/directors/{id}", dir_h.Get)
	mux.HandleFunc("GET /api/directors", dir_h.List)
	mux.HandleFunc("POST /api/directors", dir_h.Post)
	mux.HandleFunc("PATCH /api/directors/{id}", dir_h.Patch)
	mux.HandleFunc("DELETE /api/directors/{id}", dir_h.Delete)
	mux.HandleFunc("GET /api/media/{id}", media_h.Get)
	mux.HandleFunc("GET /api/media", media_h.List)
	mux.HandleFunc("POST /api/media", media_h.Post)
	mux.HandleFunc("PATCH /api/media/{id}", media_h.Patch)
	mux.HandleFunc("DELETE /api/media/{id}", media_h.Delete)
	mux.HandleFunc("GET /api/comments/{id}", comment_h.Get)
	mux.HandleFunc("GET /api/comments", comment_h.List)
	mux.HandleFunc("POST /api/comments", comment_h.Post)
	mux.HandleFunc("PATCH /api/comments/{id}", comment_h.Patch)
	mux.HandleFunc("DELETE /api/comments/{id}", comment_h.Delete)
}
