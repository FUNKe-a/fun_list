package db

import (
	"database/sql"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

func Init(dbPath string) *sql.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		slog.Error("Failed to create directory for database.", "err", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("Failed to open database.", "err", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		slog.Error("Failed to apply migration.", "err", err)
	}
	slog.Info("Applied migrations.", "amount", n)

	// content_u, _ := os.ReadFile("db/test_data/app_users.sql")
	// content_d, _ := os.ReadFile("db/test_data/directors.sql")
	// content_m, _ := os.ReadFile("db/test_data/media.sql")
	// content_c, _ := os.ReadFile("db/test_data/comments.sql")
	//
	// db.Exec(string(content_u))
	// db.Exec(string(content_d))
	// db.Exec(string(content_m))
	// db.Exec(string(content_c))

	return db
}
