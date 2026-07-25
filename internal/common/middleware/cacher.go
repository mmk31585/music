package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ResponseCache provides Redis-backed caching for GET endpoints.
// It reduces database and compute load for frequently accessed,
// infrequently changing data — artist/album/track pages, genre
// listings, catalog queries, feature flags, etc.
//
// Use cases and recommended TTLs:
//   - Artist / album / track pages: 5 minutes
//   - Genre listings: 10 minutes
//   - Search results: 1 minute
//   - Feature flags / config: 5 minutes
//   - User library (playlists, favorites): 30 seconds
//
// Cache keys include the request path + query string and are
// auth-user-aware: if a user is authenticated, their user ID is
// mixed into the key so private data isn't served to other users.
type ResponseCache struct {
	rdb        *redis.Client
	defaultTTL time.Duration
}

// NewResponseCache creates a new response cache.
// Pass nil for rdb to create a no-op cache (graceful degradation
// when Redis is unavailable). This is safe for local development.
func NewResponseCache(rdb *redis.Client, defaultTTL time.Duration) *ResponseCache {
	return &ResponseCache{
		rdb:        rdb,
		defaultTTL: defaultTTL,
	}
}

// Cache returns a Gin middleware that caches GET responses in Redis.
// Only 2xx responses are cached. Handlers can opt out per-request
// by calling middleware.CacheBypass(c).
//
// Usage:
//
//	router.GET("/api/v1/artists/:id", cache.Cache(5*time.Minute), handler)
//
// The middleware checks Redis before calling the handler. On a cache
// HIT (including ETag-based conditional requests), the handler is
// skipped entirely and the cached response is written directly.
func (rc *ResponseCache) Cache(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rc.rdb == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		cacheKey := rc.key(c)
		if rc.serveCached(c, cacheKey) {
			return // cache HIT, response already written and chain aborted
		}

		// Capture the response body written by the handler
		buf := new(bytes.Buffer)
		cw := &captureWriter{ResponseWriter: c.Writer, buf: buf}
		c.Writer = cw

		c.Next()

		// Only cache successful 2xx responses
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		// Check per-request bypass (e.g., for user-specific data)
		if c.GetBool("_cache_bypass") {
			return
		}

		// Build the cached response entry
		body := buf.Bytes()
		etag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))

		entry := &cacheEntry{
			Status:  c.Writer.Status(),
			Body:    body,
			ETag:    etag,
			Type:    c.Writer.Header().Get("Content-Type"),
			Enc:     c.Writer.Header().Get("Content-Encoding"),
			Lang:    c.Writer.Header().Get("Content-Language"),
		}

		ttlToUse := ttl
		if ttlToUse <= 0 {
			ttlToUse = rc.defaultTTL
		}

		ctx := c.Request.Context()
		if err := rc.rdb.Set(ctx, cacheKey, entry.Encode(), ttlToUse).Err(); err != nil {
			// Non-fatal: next request will be a cache miss
			return
		}
	}
}

// CacheBypass marks the current request as uncacheable.
// Call from inside a handler:
//
//	middleware.CacheBypass(c)
func CacheBypass(c *gin.Context) {
	c.Set("_cache_bypass", true)
}

// key builds a deterministic, collision-resistant cache key.
// Authenticated users get their own key to prevent data leaks.
func (rc *ResponseCache) key(c *gin.Context) string {
	raw := c.Request.URL.RequestURI()
	if userID, exists := c.Get("auth_user_id"); exists {
		if id, ok := userID.(string); ok && id != "" {
			raw = id + ":" + raw
		}
	}
	h := sha256.Sum256([]byte(raw))
	return "cache:resp:" + hex.EncodeToString(h[:])
}

// serveCached attempts to serve a cached response. Returns true on hit.
func (rc *ResponseCache) serveCached(c *gin.Context, key string) bool {
	ctx := c.Request.Context()

	data, err := rc.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}

	entry, err := decodeCacheEntry(data)
	if err != nil {
		return false
	}

	// Support conditional GET (If-None-Match → 304)
	if match := c.GetHeader("If-None-Match"); match == entry.ETag {
		c.Status(http.StatusNotModified)
		c.Abort()
		return true
	}

	// Write cached response headers
	if entry.Type != "" {
		c.Header("Content-Type", entry.Type)
	}
	if entry.Enc != "" {
		c.Header("Content-Encoding", entry.Enc)
	}
	if entry.Lang != "" {
		c.Header("Content-Language", entry.Lang)
	}
	c.Header("ETag", entry.ETag)
	c.Header("X-Cache", "HIT")

	c.Data(entry.Status, entry.Type, entry.Body)
	c.Abort()
	return true
}

// ── binary cache entry format ──
//
// To avoid the overhead and ambiguity of JSON encoding (especially for
// binary response bodies that would get base64-expanded in JSON), we
// use a compact binary format:
//
//	[1B: flags][4B: status][4B: typeLen][N: type][4B: encLen][N: enc][4B: langLen][N: lang][4B: etagLen][N: etag][4B: bodyLen][N: body]

type cacheEntry struct {
	Status int
	Body   []byte
	ETag   string
	Type   string
	Enc    string
	Lang   string
}

func (e *cacheEntry) Encode() []byte {
	buf := make([]byte, 0, 128+len(e.Body))

	// Flags byte (reserved for future use)
	buf = append(buf, 0)

	// Status
	buf = binary.BigEndian.AppendUint32(buf, uint32(e.Status))

	// Content-Type
	binaryWriteString(&buf, e.Type)

	// Content-Encoding
	binaryWriteString(&buf, e.Enc)

	// Content-Language
	binaryWriteString(&buf, e.Lang)

	// ETag
	binaryWriteString(&buf, e.ETag)

	// Body length + body
	binaryWriteBytes(&buf, e.Body)

	return buf
}

func decodeCacheEntry(data []byte) (*cacheEntry, error) {
	r := bytes.NewReader(data)
	e := &cacheEntry{}

	// Flags (1 byte)
	if _, err := r.ReadByte(); err != nil {
		return nil, fmt.Errorf("cache: read flags: %w", err)
	}

	// Status (4 bytes)
	status, err := binaryReadUint32(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read status: %w", err)
	}
	e.Status = int(status)

	// Content-Type
	e.Type, err = binaryReadString(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read type: %w", err)
	}

	// Content-Encoding
	e.Enc, err = binaryReadString(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read enc: %w", err)
	}

	// Content-Language
	e.Lang, err = binaryReadString(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read lang: %w", err)
	}

	// ETag
	e.ETag, err = binaryReadString(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read etag: %w", err)
	}

	// Body
	e.Body, err = binaryReadBytes(r)
	if err != nil {
		return nil, fmt.Errorf("cache: read body: %w", err)
	}

	return e, nil
}

// ── binary helpers ──

func binaryWriteString(buf *[]byte, s string) {
	*buf = binary.BigEndian.AppendUint32(*buf, uint32(len(s)))
	*buf = append(*buf, s...)
}

func binaryWriteBytes(buf *[]byte, b []byte) {
	*buf = binary.BigEndian.AppendUint32(*buf, uint32(len(b)))
	*buf = append(*buf, b...)
}

func binaryReadUint32(r *bytes.Reader) (uint32, error) {
	var scratch [4]byte
	if _, err := io.ReadFull(r, scratch[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(scratch[:]), nil
}

func binaryReadString(r *bytes.Reader) (string, error) {
	n, err := binaryReadUint32(r)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func binaryReadBytes(r *bytes.Reader) ([]byte, error) {
	n, err := binaryReadUint32(r)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// ── response body capture ──

type captureWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *captureWriter) Write(data []byte) (int, error) {
	w.buf.Write(data)
	return w.ResponseWriter.Write(data)
}
