package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const searchKeyPrefix = "cache:search:"

type SearchCache struct {
	rdb redis.UniversalClient
	ttl time.Duration
}

func NewSearchCache(rdb redis.UniversalClient, ttl time.Duration) *SearchCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &SearchCache{rdb: rdb, ttl: ttl}
}

func cacheKey(query string) string {
	return fmt.Sprintf("%s%x", searchKeyPrefix, query)
}

func (c *SearchCache) Get(ctx context.Context, query string) ([]byte, bool, error) {
	data, err := c.rdb.Get(ctx, cacheKey(query)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return data, true, nil
}

func (c *SearchCache) Set(ctx context.Context, query string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, cacheKey(query), data, c.ttl).Err()
}

func (c *SearchCache) Delete(ctx context.Context, query string) error {
	return c.rdb.Del(ctx, cacheKey(query)).Err()
}

func (c *SearchCache) TTL(ctx context.Context) time.Duration {
	return c.ttl
}
