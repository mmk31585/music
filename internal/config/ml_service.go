package config

import "time"

type MLServiceConfig struct {
	BaseURL           string
	WebhookHMACSecret string
	RequestTimeout    time.Duration
}

func loadMLServiceConfig() MLServiceConfig {
	return MLServiceConfig{
		BaseURL:           getEnv("ML_SERVICE_BASE_URL", "http://localhost:8000"),
		WebhookHMACSecret: getEnv("ML_SERVICE_WEBHOOK_HMAC_SECRET", ""),
		RequestTimeout:    getEnvAsDurationSeconds("ML_SERVICE_REQUEST_TIMEOUT", 30),
	}
}
