package config

type LogConfig struct {
	Level string
}

func loadLogConfig() LogConfig {
	return LogConfig{
		Level: getEnv("LOG_LEVEL", "debug"),
	}
}
