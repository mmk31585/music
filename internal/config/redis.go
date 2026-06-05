package config

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}
}
