package config

import "time"

type AppConfig struct {
	Name            string
	Env             string
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Name:            getEnv("APP_NAME", "musicapp"),
		Env:             getEnv("APP_ENV", "development"),
		Host:            getEnv("APP_HOST", "0.0.0.0"),
		Port:            getEnv("APP_PORT", "8080"),
		ReadTimeout:     getEnvAsDurationSeconds("APP_READ_TIMEOUT", 15),
		WriteTimeout:    getEnvAsDurationSeconds("APP_WRITE_TIMEOUT", 15),
		IdleTimeout:     getEnvAsDurationSeconds("APP_IDLE_TIMEOUT", 60),
		ShutdownTimeout: getEnvAsDurationSeconds("APP_SHUTDOWN_TIMEOUT", 10),
	}
}
