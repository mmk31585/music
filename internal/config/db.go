package config

import (
	"os"
)

type PostgresConfig struct {
	// URL is the pooled connection string (via pgbouncer) used by the app's
	// connection pool. Suitable for normal query traffic.
	URL string

	// UnpooledURL is the direct connection string (no pgbouncer) used for
	// migrations and operations that can't go through a connection pooler.
	// Falls back to URL if not set.
	UnpooledURL string
}

func loadPostgresConfig() PostgresConfig {
	// Prefer DATABASE_URL / POSTGRES_URL (pooled via pgbouncer)
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = os.Getenv("POSTGRES_URL")
	}

	// DATABASE_URL_UNPOOLED / POSTGRES_URL_NON_POOLING for direct connections
	unpooledURL := os.Getenv("DATABASE_URL_UNPOOLED")
	if unpooledURL == "" {
		unpooledURL = os.Getenv("POSTGRES_URL_NON_POOLING")
	}
	if unpooledURL == "" {
		unpooledURL = url
	}

	return PostgresConfig{
		URL:         url,
		UnpooledURL: unpooledURL,
	}
}
