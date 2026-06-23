package config

import (
	"fmt"
	"os"
	"strings"
)

// ValidateProductionConfig checks that all critical configuration values are set
// appropriately for the current environment. In production mode it panics if any
// required secret is still using a placeholder value.
func ValidateProductionConfig(cfg *Config) {
	if cfg.App.Env != "production" {
		// In non-production, just warn about weak secrets
		checkWeakSecret(cfg.Auth.JWTAccessSecret, "JWT_ACCESS_SECRET")
		checkWeakSecret(cfg.Auth.JWTRefreshSecret, "JWT_REFRESH_SECRET")
		return
	}

	// ── Production: hard-fail on placeholder secrets ──
	assertStrongSecret(cfg.Auth.JWTAccessSecret, "JWT_ACCESS_SECRET")
	assertStrongSecret(cfg.Auth.JWTRefreshSecret, "JWT_REFRESH_SECRET")

	// Ensure TLS-friendly settings
	if strings.ToLower(cfg.App.Host) == "0.0.0.0" {
		fmt.Fprintf(os.Stderr, "WARNING: APP_HOST is 0.0.0.0 in production. Bind to a specific interface if possible.\n")
	}
}

func checkWeakSecret(value, name string) {
	if value == "" || value == "change_me_access_secret" || value == "change_me_refresh_secret" {
		fmt.Fprintf(os.Stderr, "WARNING: %s is using a placeholder or empty value. This is INSECURE for production.\n", name)
		fmt.Fprintf(os.Stderr, "  Generate a strong secret: openssl rand -hex 32\n")
	}
}

func assertStrongSecret(value, name string) {
	if value == "change_me_access_secret" || value == "change_me_refresh_secret" {
		panic(fmt.Sprintf(
			"FATAL [%s]: Default placeholder secret detected. "+
				"Set a strong random %s environment variable before starting in production.\n"+
				"  Example: export %s=$(openssl rand -hex 32)",
			name, name, name,
		))
	}
	if value == "" {
		panic(fmt.Sprintf(
			"FATAL [%s]: %s is empty and cannot be empty in production. "+
				"Set a strong random value via environment variable.",
			name, name,
		))
	}
	if len(value) < 32 {
		fmt.Fprintf(os.Stderr,
			"WARNING: %s is shorter than 32 characters (%d chars). "+
				"Consider using a stronger secret for production.\n",
			name, len(value),
		)
	}
}
