package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitIP applies a rate limit per client IP address.
func RateLimitIP(rdb *redis.Client, requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:ip:%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		count, err := rdb.Get(c.Request.Context(), key).Int()
		if err == redis.Nil {
			rdb.Set(c.Request.Context(), key, 1, window)
			c.Next()
			return
		}
		if err != nil {
			c.Next()
			return
		}

		if count >= requests {
			ttl, ttlErr := rdb.TTL(c.Request.Context(), key).Result()
			if ttlErr == nil {
				seconds := int(ttl.Seconds())
				if seconds < 0 {
					seconds = int(window.Seconds())
				}
				if seconds > 0 {
					c.Header("Retry-After", strconv.Itoa(seconds))
				}
			}
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

// RateLimitOptional applies a rate limit per user ID when available,
// falling back to IP address for unauthenticated callers. This is designed
// for routes using OptionalAuth to prevent abuse by anonymous users.
func RateLimitOptional(rdb *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		// Build a unique key scoped to user+route to prevent cross-route exhaustion
		var key string
		userIDRaw, exists := c.Get("auth_user_id")
		if exists {
			// Safe type assertion with ok check — prevents panic if value is not string
			userIDStr, ok := userIDRaw.(string)
			if ok && userIDStr != "" {
				key = fmt.Sprintf("ratelimit:user:%s:%s:%s", userIDStr, c.Request.Method, c.Request.URL.Path)
			} else {
				key = fmt.Sprintf("ratelimit:ip:%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
			}
		} else {
			key = fmt.Sprintf("ratelimit:ip:%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		}

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
