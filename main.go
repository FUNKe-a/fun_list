package main

import (
	"errors"
	"fmt"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"github.com/FUNKe-a/fun_list/internal/handlers"
	"github.com/FUNKe-a/fun_list/internal/db"
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
	dir_h := handlers.DirectorHandler{
		Queries: q,
	}

	http.HandleFunc("GET /api/directors/{id}", dir_h.Get) 

	http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
