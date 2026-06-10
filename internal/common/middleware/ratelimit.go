package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		key := "ratelimit:" + c.ClientIP() + ":" + c.Request.URL.Path
		now := time.Now().Unix()

		pipe := rdb.Pipeline()
		countCmd := pipe.Get(c.Request.Context(), key)
		ttlCmd := pipe.TTL(c.Request.Context(), key)
		_, _ = pipe.Exec(c.Request.Context())

		count, _ := strconv.Atoi(countCmd.Val())
		ttl := ttlCmd.Val()

		if ttl <= 0 {
			rdb.Set(c.Request.Context(), key, 1, time.Minute)
			count = 0
		}

		if count >= requestsPerMinute {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    "RATE_LIMITED",
				"message": "too many requests, try again later",
			})
			return
		}

		rdb.Incr(c.Request.Context(), key)

		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(requestsPerMinute-count-1))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now+60, 10))

		c.Next()
	}
}

func RateLimitByUser(rdb *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		userID, exists := c.Get("auth_user_id")
		if !exists {
			c.Next()
			return
		}

		key := "ratelimit:user:" + userID.(string) + ":" + c.Request.URL.Path

		count, err := rdb.Get(c.Request.Context(), key).Int()
		if err == redis.Nil {
			rdb.Set(c.Request.Context(), key, 1, time.Minute)
			c.Next()
			return
		}
		if err != nil {
			c.Next()
			return
		}

		if count >= requestsPerMinute {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    "RATE_LIMITED",
				"message": "too many requests, try again later",
			})
			return
		}

		rdb.Incr(c.Request.Context(), key)
		c.Next()
	}
}
