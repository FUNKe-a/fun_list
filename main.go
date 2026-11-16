package main

import (
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/configs"
	"github.com/FUNKe-a/fun_list/internal/db"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"github.com/FUNKe-a/fun_list/internal/handlers"
	"github.com/FUNKe-a/fun_list/internal/routes"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

func main() {
	db_conn := db.Init("data/data.db")
	defer db_conn.Close()

	priv, pub := configs.ParseJWTKeys()
	handler := handlers.Handler{
		DB:      db_conn,
		Queries: sqlc.New(db_conn),
		Tokens:  auth.NewAuthService(priv, pub, "fun_list", 10*time.Minute),
	}

	mux := routes.CreateMux(&handler)

	http.ListenAndServe(configs.Port, mux)
}
