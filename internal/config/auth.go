package config

import "time"

type AuthConfig struct {
	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
}

func loadAuthConfig() AuthConfig {
	// NOTE: In production, JWT_ACCESS_SECRET and JWT_REFRESH_SECRET MUST be set
	// to strong random values (at least 256 bits). Do NOT use defaults.
	return AuthConfig{
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "change_me_access_secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change_me_refresh_secret"),
		AccessTTL:        time.Duration(getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		RefreshTTL:       time.Duration(getEnvAsInt("JWT_REFRESH_TTL_DAYS", 30)) * 24 * time.Hour,
	}
}
