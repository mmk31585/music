package config

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	// Addr, Password, DB are populated when a Redis URL is configured.
	// All zero values mean Redis is not configured and should be skipped.
	Addr     string
	Password string
	DB       int

	// DialTimeout limits how long we wait for the first connection.
	// MaxRetries limits reconnect attempts so startup is not blocked.
	DialTimeout time.Duration
	MaxRetries  int
}

func loadRedisConfig() RedisConfig {
	// Precedence: REDIS_URL, then music_REDIS_URL (Vercel integration prefix),
	// then KV_URL. Parse with redis.ParseURL which supports rediss://.
	var redisURL string
	for _, env := range []string{"REDIS_URL", "music_REDIS_URL", "KV_URL"} {
		if v := os.Getenv(env); v != "" {
			redisURL = v
			break
		}
	}

	if redisURL == "" {
		// No Redis configured — skip entirely. Never dial localhost in production.
		return RedisConfig{}
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// Malformed URL — treat as not configured rather than dialing localhost.
		return RedisConfig{}
	}

	return RedisConfig{
		Addr:      opts.Addr,
		Password:  opts.Password,
		DB:        opts.DB,
		DialTimeout: 1 * time.Second,
		MaxRetries:  1,
	}
}
