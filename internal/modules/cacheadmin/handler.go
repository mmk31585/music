package cacheadmin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	rdb *redis.Client
}

func NewHandler(rdb *redis.Client) *Handler {
	return &Handler{rdb: rdb}
}

// CacheStatus holds the current cache status returned by GetStatus.
type CacheStatus struct {
	Connected     bool   `json:"connected"`
	KeyCount      int64  `json:"key_count"`
	MemoryBytes   int64  `json:"memory_bytes"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

// GetStatus returns the current Redis cache status.
func (h *Handler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()

	status := CacheStatus{Connected: false}

	// Ping to check connectivity
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusOK, status)
		return
	}

	status.Connected = true

	// Get key count via DBSIZE
	keyCount, err := h.rdb.DBSize(ctx).Result()
	if err == nil {
		status.KeyCount = keyCount
	}

	// Get memory and uptime info from INFO
	info, err := h.rdb.Info(ctx, "memory", "server").Result()
	if err == nil {
		status.MemoryBytes = parseInfoInt(info, "used_memory")
		status.UptimeSeconds = parseInfoInt(info, "uptime_in_seconds")
	}

	c.JSON(http.StatusOK, status)
}

// FlushCache flushes all keys from the Redis cache.
func (h *Handler) FlushCache(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.rdb.FlushDB(ctx).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to flush cache",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache flushed",
	})
}

// parseInfoInt extracts an integer value for the given key from a Redis INFO response.
func parseInfoInt(info, key string) int64 {
	prefix := key + ":"
	for i := 0; i < len(info); i++ {
		// Look for key at the start of a line
		if i == 0 || info[i-1] == '\n' {
			if i+len(prefix) <= len(info) && info[i:i+len(prefix)] == prefix {
				start := i + len(prefix)
				end := start
				for end < len(info) && info[end] != '\r' && info[end] != '\n' {
					end++
				}
				if start < end {
					var val int64
					for _, ch := range info[start:end] {
						if ch >= '0' && ch <= '9' {
							val = val*10 + int64(ch-'0')
						}
					}
					return val
				}
			}
		}
	}
	return 0
}
