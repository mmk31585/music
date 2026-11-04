// Command migrate applies the embedded Goose migrations to the database.
//
// It is run as a one-shot container (see the `migrate` service in
// docker-compose.yml) before the API/worker start. It exists so migrations can
// be applied through the same self-contained artifact used in production
// (embed.FS) without depending on an external `goose`/`psql` toolchain.
//
// Migrations must never run inside an HTTP request (see AGENTS.md §7).
package main

import (
	"fmt"
	"os"

	"music/internal/platform/database"
)

// resolveURL picks the first non-empty Postgres connection string from the
// environment. The unpooled (direct) connection is preferred because some
// migration DDL cannot run through a transaction pooler like pgbouncer.
func resolveURL() string {
	for _, key := range []string{
		"POSTGRES_URL_NON_POOLING",
		"DATABASE_URL_UNPOOLED",
		"POSTGRES_URL",
		"DATABASE_URL",
	} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func main() {
	url := resolveURL()
	if url == "" {
		fmt.Fprintln(os.Stderr, "migrate: set POSTGRES_URL (or POSTGRES_URL_NON_POOLING)")
		os.Exit(1)
	}

	if err := database.RunMigrations(url); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("migrate: up to date")
}
