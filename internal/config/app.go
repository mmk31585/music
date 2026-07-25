package config

import (
	"os"
	"strconv"
	"time"
)

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
	name := os.Getenv("APP_NAME")
	if name == "" {
		name = "musicapp"
	}
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	readTimeout := 30
	if v := os.Getenv("APP_READ_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			readTimeout = parsed
		}
	}
	writeTimeout := 60
	if v := os.Getenv("APP_WRITE_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			writeTimeout = parsed
		}
	}
	idleTimeout := 60
	if v := os.Getenv("APP_IDLE_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			idleTimeout = parsed
		}
	}
	shutdownTimeout := 10
	if v := os.Getenv("APP_SHUTDOWN_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			shutdownTimeout = parsed
		}
	}
	return AppConfig{
		Name:            name,
		Env:             env,
		Host:            host,
		Port:            port,
		ReadTimeout:     time.Duration(readTimeout) * time.Second,
		WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		IdleTimeout:     time.Duration(idleTimeout) * time.Second,
		ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
	}
}
