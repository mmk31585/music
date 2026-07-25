package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRequest(r *gin.Engine, method, path, clientIP string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	if clientIP != "" {
		req.Header.Set("X-Forwarded-For", clientIP)
	}
	r.ServeHTTP(w, req)
	return w
}

func TestRateLimitIP_PassesRequestsWithinLimit(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := range 3 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code, "request %d should pass", i+1)

		var resp map[string]bool
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["ok"])
	}
}

func TestRateLimitIP_BlocksExceedingLimit(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for range 3 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := makeRequest(r, "GET", "/test", "")
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "RATE_LIMITED", resp["code"])
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "too many requests, try again later", resp["message"])
}

func TestRateLimitIP_SetsRetryAfterHeader(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for range 3 {
		makeRequest(r, "GET", "/test", "")
	}

	w := makeRequest(r, "GET", "/test", "")
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	ra := w.Header().Get("Retry-After")
	assert.NotEmpty(t, ra, "Retry-After header must be set")

	seconds, err := strconv.Atoi(ra)
	require.NoError(t, err, "Retry-After must be a numeric value")
	assert.GreaterOrEqual(t, seconds, 55, "Retry-After should be near the window duration")
}

func TestRateLimitIP_PerIPIsolation(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	ipA, ipB := "192.168.1.1", "192.168.1.2"

	for i := range 3 {
		w := makeRequest(r, "GET", "/test", ipA)
		assert.Equal(t, http.StatusOK, w.Code, "IP A request %d", i+1)
	}

	for i := range 3 {
		w := makeRequest(r, "GET", "/test", ipB)
		assert.Equal(t, http.StatusOK, w.Code, "IP B request %d", i+1)
	}

	w := makeRequest(r, "GET", "/test", ipA)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "IP A should be rate-limited")

	w = makeRequest(r, "GET", "/test", ipB)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "IP B should be rate-limited")
}

func TestRateLimitIP_PerRouteIsolation(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/path/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"path": "a"})
	})
	r.GET("/path/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"path": "b"})
	})

	for range 3 {
		w := makeRequest(r, "GET", "/path/a", "")
		assert.Equal(t, http.StatusOK, w.Code, "/path/a should pass within limit")
	}

	w := makeRequest(r, "GET", "/path/a", "")
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "/path/a should be exhausted")

	for range 3 {
		w := makeRequest(r, "GET", "/path/b", "")
		assert.Equal(t, http.StatusOK, w.Code, "/path/b should still pass (independent counter)")
	}

	w = makeRequest(r, "GET", "/path/b", "")
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "/path/b should also be exhausted")
}

func TestRateLimitIP_WindowReset(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, 2*time.Second))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for range 3 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := makeRequest(r, "GET", "/test", "")
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "should be rate-limited within window")

	mr.FastForward(3 * time.Second)

	for range 3 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code, "request should pass after window reset")
	}
}

func TestRateLimitIP_NilRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(nil, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := range 10 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code, "request %d should pass with nil Redis", i+1)
	}
}

func TestRateLimitIP_GracefulDegradationOnRedisError(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:1",
		MaxRetries:   0,
		DialTimeout:  5 * time.Millisecond,
		ReadTimeout:  5 * time.Millisecond,
		WriteTimeout: 5 * time.Millisecond,
		PoolTimeout:  5 * time.Millisecond,
		PoolSize:     1,
		MinIdleConns: 0,
	})
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitIP(rdb, 3, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := range 5 {
		w := makeRequest(r, "GET", "/test", "")
		assert.Equal(t, http.StatusOK, w.Code, "request %d should pass despite Redis error", i+1)
	}
}
