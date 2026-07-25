package config

import "os"

type LogConfig struct {
	Level string
}

func loadLogConfig() LogConfig {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "debug"
	}
	return LogConfig{
		Level: level,
	}
}
