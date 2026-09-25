package db

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

func Init(dbPath string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory for database: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return nil, fmt.Errorf("Failed to apply migration: %v", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)

	content_d, _ := os.ReadFile("db/test_data/directors.sql")
	content_m, _ := os.ReadFile("db/test_data/media.sql")
	content_c, _ := os.ReadFile("db/test_data/comments.sql")

	db.Exec(string(content_d))
	db.Exec(string(content_m))
	db.Exec(string(content_c))

	return db, nil
}
