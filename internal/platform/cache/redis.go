package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	MaxRetries   int
	MinIdleConns int
}

// NewRedis creates a Redis client. If Addr is empty, Redis is not configured
// and the function returns a nil client with a nil error — callers must handle
// a nil client gracefully (degrade without cache).
func NewRedis(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, fmt.Errorf("instrument redis tracing: %w", err)
	}

	return client, nil
}
