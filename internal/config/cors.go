package config

import (
	"os"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func loadCORSConfig() CORSConfig {
	var origins []string
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		parts := strings.Split(v, ",")
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			s := strings.TrimSpace(p)
			s = strings.TrimRight(s, "/")
			if s != "" {
				trimmed = append(trimmed, s)
			}
		}
		origins = trimmed
	}
	return CORSConfig{AllowedOrigins: origins}
}
