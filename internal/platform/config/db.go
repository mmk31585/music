package config

type PostgresConfig struct {
	URL string
}

func loadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URL: getEnv("POSTGRES_URL", ""),
	}
}
