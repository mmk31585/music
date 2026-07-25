package config

import "os"

type PostgresConfig struct {
	URL string
}

func loadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URL: os.Getenv("POSTGRES_URL"),
	}
}
