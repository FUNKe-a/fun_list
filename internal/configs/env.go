package configs

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func ParseJWTKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file, using environment variables: %v", err)
	}

	privBytes, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		slog.Error("Private JWT key not found.")
	}

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		slog.Error("Couldn't parse private JWT key.")
	}

	pubBytes, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY_PATH"))
	if err != nil {
		slog.Error("Public JWT key not found.")
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		slog.Error("Couldn't parse public JWT key.")
	}

	return priv, pub
}
