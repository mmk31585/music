package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"music/migrations"
)

// resolveURL resolves ${VAR} references in a connection string.
// If the entire value is ${VAR_NAME}, it is replaced with os.Getenv(VAR_NAME).
func resolveURL(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		name := value[2 : len(value)-1]
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return value
}

func RunMigrations(databaseURL string) error {
	databaseURL = resolveURL(databaseURL)
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	// Use embedded migrations so the binary does not depend on the
	// migrations/ directory existing in the deployment artifact.
	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
