package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
func getEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")

	if value == "" {
		return defaultValue
	}

	switch value {
	case "true", "TRUE", "1", "yes", "YES", "y", "Y":
		return true
	case "false", "FALSE", "0", "no", "NO", "n", "N":
		return false
	default:
		return defaultValue
	}
}

func getInt64Env(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsDurationSeconds(key string, fallback int) time.Duration {
	return time.Duration(getEnvAsInt(key, fallback)) * time.Second
}

func getEnvAsStringSlice(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

func mbToBytes(mb int64) int64 {
	return mb << 20
}
