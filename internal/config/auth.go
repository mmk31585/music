package config

import (
	"os"
	"strconv"
	"time"
)

type AuthConfig struct {
	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
}

func loadAuthConfig() AuthConfig {
	jwtAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if jwtAccessSecret == "" {
		jwtAccessSecret = "change_me_access_secret"
	}
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		jwtRefreshSecret = "change_me_refresh_secret"
	}
	accessTTL := 15
	if v := os.Getenv("JWT_ACCESS_TTL_MINUTES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			accessTTL = parsed
		}
	}
	refreshTTL := 30
	if v := os.Getenv("JWT_REFRESH_TTL_DAYS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			refreshTTL = parsed
		}
	}
	return AuthConfig{
		JWTAccessSecret:  jwtAccessSecret,
		JWTRefreshSecret: jwtRefreshSecret,
		AccessTTL:        time.Duration(accessTTL) * time.Minute,
		RefreshTTL:       time.Duration(refreshTTL) * 24 * time.Hour,
	}
}
