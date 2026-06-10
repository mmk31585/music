package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (c *Cache) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	if err := c.Get(ctx, key, dest); err == nil {
		return nil
	}

	data, err := fn()
	if err != nil {
		return err
	}

	if err := c.Set(ctx, key, data, ttl); err != nil {
		return err
	}

	raw, _ := json.Marshal(data)
	return json.Unmarshal(raw, dest)
}

func (c *Cache) RememberStruct(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	return c.Remember(ctx, key, ttl, fn, dest)
}

func (c *Cache) RememberSlice(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	return c.Remember(ctx, key, ttl, fn, dest)
}

func (c *Cache) Exists(ctx context.Context, keys ...string) (bool, error) {
	n, err := c.client.Exists(ctx, keys...).Result()
	return n > 0, err
}

func (c *Cache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, pattern).Result()
}

func (c *Cache) Client() *redis.Client {
	return c.client
}
