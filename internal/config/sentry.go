package config

import "os"

type SentryConfig struct {
	DSN string
}

func loadSentryConfig() SentryConfig {
	return SentryConfig{
		DSN: os.Getenv("SENTRY_DSN"),
	}
}
