package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ──

func setupEngine(t *testing.T) (*gin.Engine, *ResponseCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cache := NewResponseCache(rdb, 5*time.Minute)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r, cache, mr
}

func performRequest(r *gin.Engine, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	r.ServeHTTP(w, req)
	return w
}

// ── binary encoding round-trip ──

func TestCacheEntryEncodeDecode_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry *cacheEntry
	}{
		{
			name: "typical JSON response",
			entry: &cacheEntry{
				Status: http.StatusOK,
				Body:   []byte(`{"data":"hello","count":42}`),
				ETag:   `"abc123def456"`,
				Type:   "application/json",
				Enc:    "gzip",
				Lang:   "fa-IR",
			},
		},
		{
			name: "all empty fields",
			entry: &cacheEntry{
				Status: http.StatusNoContent,
				Body:   nil,
				ETag:   "",
				Type:   "",
				Enc:    "",
				Lang:   "",
			},
		},
		{
			name: "binary body",
			entry: &cacheEntry{
				Status: http.StatusOK,
				Body:   []byte{0x00, 0x01, 0xFF, 0xFE, 0x80, 0x7F},
				ETag:   `"binary"`,
				Type:   "application/octet-stream",
				Enc:    "",
				Lang:   "",
			},
		},
		{
			name: "large body",
			entry: &cacheEntry{
				Status: http.StatusOK,
				Body:   bytes.Repeat([]byte("A"), 10000),
				ETag:   `"large"`,
				Type:   "text/plain",
				Enc:    "",
				Lang:   "en",
			},
		},
		{
			name: "empty body slice (encodes as nil)",
			entry: &cacheEntry{
				Status: http.StatusOK,
				Body:   []byte{},
				ETag:   `"empty-body"`,
				Type:   "text/plain",
				Enc:    "",
				Lang:   "",
			},
		},
		{
			name: "unicode strings",
			entry: &cacheEntry{
				Status: http.StatusOK,
				Body:   []byte("سلام دنیا"),
				ETag:   `"persian"`,
				Type:   "text/plain; charset=utf-8",
				Enc:    "",
				Lang:   "fa",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.entry.Encode()
			require.NotNil(t, data, "Encode() should return non-nil bytes")

			decoded, err := decodeCacheEntry(data)
			require.NoError(t, err, "decodeCacheEntry should not error")
			require.NotNil(t, decoded, "decoded entry should not be nil")

			assert.Equal(t, tt.entry.Status, decoded.Status)
			assert.Equal(t, tt.entry.ETag, decoded.ETag)
			assert.Equal(t, tt.entry.Type, decoded.Type)
			assert.Equal(t, tt.entry.Enc, decoded.Enc)
			assert.Equal(t, tt.entry.Lang, decoded.Lang)

			if tt.entry.Body == nil || len(tt.entry.Body) == 0 {
				assert.Nil(t, decoded.Body, "empty/nil body should decode to nil")
			} else {
				assert.Equal(t, tt.entry.Body, decoded.Body, "body mismatch")
			}
		})
	}
}

func TestCacheEntryDecode_CorruptedData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty slice", []byte{}},
		{"truncated after flags", []byte{0}},
		{"truncated after status", []byte{0, 0, 0, 0, 0}},
		{"garbage bytes", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeCacheEntry(tt.data)
			assert.Error(t, err)
		})
	}
}

// ── binary helpers ──

func TestBinaryWriteReadString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"normal string", "hello world"},
		{"empty string", ""},
		{"unicode", "سلام دنیا"},
		{"long string", string(bytes.Repeat([]byte("x"), 5000))},
		{"special chars", "{\"key\":\"value\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf []byte
			binaryWriteString(&buf, tt.input)

			r := bytes.NewReader(buf)
			got, err := binaryReadString(r)
			require.NoError(t, err)
			assert.Equal(t, tt.input, got)
			assert.Equal(t, 0, r.Len(), "reader should be fully consumed")
		})
	}
}

func TestBinaryWriteReadBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"normal bytes", []byte("hello world")},
		{"nil slice", nil},
		{"empty slice (encodes/res nil)", []byte{}},
		{"binary data", []byte{0x00, 0xFF, 0xFE, 0x80}},
		{"large bytes", bytes.Repeat([]byte{0xAB}, 5000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf []byte
			binaryWriteBytes(&buf, tt.input)

			r := bytes.NewReader(buf)
			got, err := binaryReadBytes(r)
			require.NoError(t, err)

			if tt.input == nil || len(tt.input) == 0 {
				assert.Nil(t, got, "nil/empty input should decode to nil")
			} else {
				assert.Equal(t, tt.input, got)
			}
			assert.Equal(t, 0, r.Len(), "reader should be fully consumed")
		})
	}
}

func TestBinaryReadUint32_EOF(t *testing.T) {
	_, err := binaryReadUint32(bytes.NewReader([]byte{0, 0, 0}))
	assert.Error(t, err, "short read should error")
}

func TestBinaryReadString_Truncated(t *testing.T) {
	var buf []byte
	binaryWriteString(&buf, "hello world")
	// truncate the actual string content
	truncated := buf[:len(buf)-5]

	r := bytes.NewReader(truncated)
	_, err := binaryReadString(r)
	assert.Error(t, err, "truncated string should error")
}

func TestBinaryReadBytes_Truncated(t *testing.T) {
	var buf []byte
	binaryWriteBytes(&buf, []byte("hello world"))
	truncated := buf[:len(buf)-5]

	r := bytes.NewReader(truncated)
	_, err := binaryReadBytes(r)
	assert.Error(t, err, "truncated bytes should error")
}

// ── cache middleware (with Redis) ──

func TestCacheMiddleware_MissThenHit(t *testing.T) {
	r, cache, _ := setupEngine(t)

	callCount := 0
	r.GET("/test", cache.Cache(1*time.Minute), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusOK, gin.H{"data": "hello"})
	})

	// First request: MISS
	w1 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, callCount, "handler should be called once")
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	// Second request: HIT
	w2 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 1, callCount, "handler should NOT be called on cache hit")
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
	assert.Equal(t, w1.Body.String(), w2.Body.String(), "cached body should match")
}

func TestCacheMiddleware_304NotModified(t *testing.T) {
	r, cache, _ := setupEngine(t)

	r.GET("/test", cache.Cache(1*time.Minute), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello"})
	})

	// First request: MISS, get the response
	w1 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request: HIT, capture the ETag
	w2 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	etag := w2.Header().Get("ETag")
	require.NotEmpty(t, etag, "ETag should be present on cache hit")

	// Third request: If-None-Match with matching ETag -> 304
	w3 := performRequest(r, "GET", "/test", map[string]string{"If-None-Match": etag})
	assert.Equal(t, http.StatusNotModified, w3.Code)
}

func TestCacheMiddleware_OnlyCaches2xx(t *testing.T) {
	r, cache, mr := setupEngine(t)

	callCount := 0
	r.GET("/test", cache.Cache(1*time.Minute), func(c *gin.Context) {
		callCount++
		status := c.Query("status")
		switch status {
		case "400":
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad"})
		case "500":
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops"})
		default:
			c.JSON(http.StatusOK, gin.H{"data": "ok"})
		}
	})

	// 400 response should not be cached
	w1 := performRequest(r, "GET", "/test?status=400", nil)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
	assert.Equal(t, 1, callCount)

	// Same request again — should still miss because 400 wasn't cached
	w2 := performRequest(r, "GET", "/test?status=400", nil)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	assert.Equal(t, 2, callCount, "handler should be called again (no cache)")
	assert.Empty(t, w2.Header().Get("X-Cache"), "should not be a cache hit")

	mr.FlushAll()

	// 500 response should not be cached
	w3 := performRequest(r, "GET", "/test?status=500", nil)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
	assert.Equal(t, 3, callCount)

	w4 := performRequest(r, "GET", "/test?status=500", nil)
	assert.Equal(t, http.StatusInternalServerError, w4.Code)
	assert.Equal(t, 4, callCount, "handler should be called again (no cache)")

	// 200 response should be cached
	w5 := performRequest(r, "GET", "/test?status=200", nil)
	assert.Equal(t, http.StatusOK, w5.Code)
	assert.Equal(t, 5, callCount)

	w6 := performRequest(r, "GET", "/test?status=200", nil)
	assert.Equal(t, http.StatusOK, w6.Code)
	assert.Equal(t, 5, callCount, "handler should NOT be called (cache hit)")
	assert.Equal(t, "HIT", w6.Header().Get("X-Cache"))
}

func TestCacheMiddleware_OnlyCachesGET(t *testing.T) {
	r, cache, _ := setupEngine(t)

	callCount := 0
	handler := func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusOK, gin.H{"method": c.Request.Method})
	}

	// Register the same path for GET, POST, PUT, DELETE with middleware
	r.GET("/test", cache.Cache(1*time.Minute), handler)
	r.POST("/test", cache.Cache(1*time.Minute), handler)
	r.PUT("/test", cache.Cache(1*time.Minute), handler)
	r.DELETE("/test", cache.Cache(1*time.Minute), handler)

	// POST - should call handler, not cached
	w1 := performRequest(r, "POST", "/test", nil)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, callCount)
	assert.Empty(t, w1.Header().Get("X-Cache"))

	// POST again - should call handler again
	w2 := performRequest(r, "POST", "/test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 2, callCount, "POST should not be cached")
	assert.Empty(t, w2.Header().Get("X-Cache"))

	// PUT - should call handler, not cached
	w3 := performRequest(r, "PUT", "/test", nil)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, 3, callCount)

	// DELETE - should call handler, not cached
	w4 := performRequest(r, "DELETE", "/test", nil)
	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Equal(t, 4, callCount)

	// GET - MISS first, then HIT
	w5 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w5.Code)
	assert.Equal(t, 5, callCount)

	w6 := performRequest(r, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w6.Code)
	assert.Equal(t, 5, callCount, "GET should be cached")
	assert.Equal(t, "HIT", w6.Header().Get("X-Cache"))
}

func TestCacheMiddleware_CacheBypass(t *testing.T) {
	r, cache, _ := setupEngine(t)

	callCount := 0
	r.GET("/bypass", cache.Cache(1*time.Minute), func(c *gin.Context) {
		callCount++
		CacheBypass(c)
		c.JSON(http.StatusOK, gin.H{"data": "secret"})
	})

	// First request
	w1 := performRequest(r, "GET", "/bypass", nil)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, callCount)

	// Second request — should NOT be cached because CacheBypass was called
	w2 := performRequest(r, "GET", "/bypass", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 2, callCount, "handler should be called again (bypass)")
	assert.Empty(t, w2.Header().Get("X-Cache"), "should not be a cache hit")
}

func TestCacheMiddleware_AuthAwareKeys(t *testing.T) {
	r, cache, _ := setupEngine(t)

	r.Use(func(c *gin.Context) {
		if uid := c.GetHeader("X-User-ID"); uid != "" {
			c.Set("auth_user_id", uid)
		}
		c.Next()
	})

	callCount := 0
	r.GET("/profile", cache.Cache(1*time.Minute), func(c *gin.Context) {
		callCount++
		uid, _ := c.Get("auth_user_id")
		c.JSON(http.StatusOK, gin.H{"user": uid})
	})

	// User A: first request - MISS
	w1 := performRequest(r, "GET", "/profile", map[string]string{"X-User-ID": "user-a"})
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, 1, callCount)
	assert.Contains(t, w1.Body.String(), "user-a")

	// User A: second request - HIT (different key for user-a)
	w2 := performRequest(r, "GET", "/profile", map[string]string{"X-User-ID": "user-a"})
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, 1, callCount, "user-a's cached response should be used")
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
	assert.Contains(t, w2.Body.String(), "user-a")

	// User B: first request - MISS (different key because different user)
	w3 := performRequest(r, "GET", "/profile", map[string]string{"X-User-ID": "user-b"})
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, 2, callCount, "user-b should get a cache miss")
	assert.Empty(t, w3.Header().Get("X-Cache"))
	assert.Contains(t, w3.Body.String(), "user-b")
}

// ── nil Redis (no-op cache) ──

func TestNilRedis_NoOpCache(t *testing.T) {
	cache := NewResponseCache(nil, 5*time.Minute)
	gin.SetMode(gin.TestMode)
	r := gin.New()

	callCount := 0
	r.GET("/test", cache.Cache(1*time.Minute), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusOK, gin.H{"data": "hello"})
	})

	// Every request should call the handler (no caching)
	for i := 0; i < 3; i++ {
		w := performRequest(r, "GET", "/test", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, i+1, callCount, "handler should be called every time with nil Redis")
		assert.Empty(t, w.Header().Get("X-Cache"), "no X-Cache header with nil Redis")
	}
}

func TestNilRedis_NewResponseCache(t *testing.T) {
	cache := NewResponseCache(nil, 0)
	assert.NotNil(t, cache, "NewResponseCache should return non-nil even with nil rdb")
	assert.Nil(t, cache.rdb)
}

// ── key generation ──

func TestKeyGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		setup    func(*gin.Context)
		prefix   string
	}{
		{
			name: "without auth",
			setup: func(c *gin.Context) {
				c.Request = httptest.NewRequest("GET", "/api/v1/tracks/123", nil)
			},
			prefix: "cache:resp:",
		},
		{
			name: "with auth user id",
			setup: func(c *gin.Context) {
				c.Request = httptest.NewRequest("GET", "/api/v1/playlists/42", nil)
				c.Set("auth_user_id", "user-abc")
			},
			prefix: "cache:resp:",
		},
		{
			name: "with empty auth user id",
			setup: func(c *gin.Context) {
				c.Request = httptest.NewRequest("GET", "/api/v1/tracks/123", nil)
				c.Set("auth_user_id", "")
			},
			prefix: "cache:resp:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setup(c)

			cache := &ResponseCache{rdb: nil}
			key := cache.key(c)
			assert.Contains(t, key, tt.prefix)
			assert.Len(t, key, len(tt.prefix)+64, "key should be prefix + hex-encoded SHA256 hash")
		})
	}
}

func TestKeyGeneration_DifferentUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	c1, _ := gin.CreateTestContext(w)
	c1.Request = httptest.NewRequest("GET", "/api/v1/profile", nil)
	c1.Set("auth_user_id", "user-a")

	c2, _ := gin.CreateTestContext(w)
	c2.Request = httptest.NewRequest("GET", "/api/v1/profile", nil)
	c2.Set("auth_user_id", "user-b")

	cache := &ResponseCache{}
	keyA := cache.key(c1)
	keyB := cache.key(c2)
	assert.NotEqual(t, keyA, keyB, "different users should have different cache keys")
}

// ── captureWriter ──

func TestCaptureWriter(t *testing.T) {
	var buf bytes.Buffer
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	cw := &captureWriter{ResponseWriter: c.Writer, buf: &buf}

	data := []byte("hello world")
	n, err := cw.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, data, buf.Bytes(), "captureWriter should copy data to buffer")
	assert.Equal(t, data, w.Body.Bytes(), "captureWriter should write through to original writer")
}
