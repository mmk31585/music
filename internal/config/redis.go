package config

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func loadRedisConfig() RedisConfig {
	// REDIS_URL takes precedence (e.g. redis://user:pass@host:port/db)
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if opts, err := redis.ParseURL(redisURL); err == nil {
			return RedisConfig{
				Addr:     opts.Addr,
				Password: opts.Password,
				DB:       opts.DB,
			}
		}
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	db := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			db = parsed
		}
	}
	return RedisConfig{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}
