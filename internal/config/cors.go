package config

import (
	"os"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func loadCORSConfig() CORSConfig {
	origins := []string{"http://localhost:3000", "http://localhost:5173"}
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		parts := strings.Split(v, ",")
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				trimmed = append(trimmed, s)
			}
		}
		if len(trimmed) > 0 {
			origins = trimmed
		}
	}
	return CORSConfig{AllowedOrigins: origins}
}
