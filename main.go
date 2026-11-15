package main

import (
	"errors"
	"fmt"
	"github.com/FUNKe-a/fun_list/internal/auth"
	"github.com/FUNKe-a/fun_list/internal/configs"
	"github.com/FUNKe-a/fun_list/internal/db"
	"github.com/FUNKe-a/fun_list/internal/db/sqlc"
	"github.com/FUNKe-a/fun_list/internal/handlers"
	"github.com/FUNKe-a/fun_list/internal/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file, using environment variables: %v", err)
	}

	privBytes, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal(err)
	}

	db_conn, err := db.Init("data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db_conn.Close()

	mux := http.NewServeMux()
	handler := handlers.Handler{
		DB: db_conn,
		Queries: sqlc.New(db_conn),
		Tokens:  auth.NewTokenService(priv, pub, "fun_list", 15*time.Minute),
	}

	routes.RoutesHandle(mux, &handler)

	http.ListenAndServe(configs.Port, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}


