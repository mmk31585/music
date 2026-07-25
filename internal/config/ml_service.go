package config

import (
	"os"
	"strconv"
	"time"
)

type MLServiceConfig struct {
	BaseURL           string
	WebhookHMACSecret string
	RequestTimeout    time.Duration
}

func loadMLServiceConfig() MLServiceConfig {
	baseURL := os.Getenv("ML_SERVICE_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}
	requestTimeout := 30
	if v := os.Getenv("ML_SERVICE_REQUEST_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			requestTimeout = parsed
		}
	}
	return MLServiceConfig{
		BaseURL:           baseURL,
		WebhookHMACSecret: os.Getenv("ML_SERVICE_WEBHOOK_HMAC_SECRET"),
		RequestTimeout:    time.Duration(requestTimeout) * time.Second,
	}
}
