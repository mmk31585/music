package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZapLogger returns a Gin middleware that logs HTTP requests using Zap.
// It includes request ID, status, latency, and skips logging for specified paths.
func GinZapLogger(log *zap.Logger, skipPaths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool, len(skipPaths))
	for _, p := range skipPaths {
		skipMap[p] = true
	}

	return func(c *gin.Context) {
		// Skip logging for specified paths (e.g., health checks)
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		requestID := getRequestID(c)

		c.Next()

		latency := time.Since(start)

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Add user ID if authenticated (you can extend this)
		if userID, exists := c.Get("auth_user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		// Add any request errors
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			log.Error("http request", fields...)
		case status >= 400:
			log.Warn("http request", fields...)
		default:
			log.Info("http request", fields...)
		}
	}
}

// GinZapRecovery recovers from panics, logs them with stack trace, and returns 500.
func GinZapRecovery(log *zap.Logger, skipPaths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool, len(skipPaths))
	for _, p := range skipPaths {
		skipMap[p] = true
	}

	return func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			defer func() {
				if err := recover(); err != nil {
					// Log with stack trace even for skipped paths? Probably yes.
					log.Error("panic recovered",
						zap.Any("error", err),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("ip", c.ClientIP()),
						zap.Stack("stack"),
					)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "internal server error",
					})
				}
			}()
			c.Next()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.Stack("stack"), // full stack trace
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

// Helper functions (you can replace with your own UUID generator)
func generateRequestID() string {
	// Example: return uuid.New().String()
	// For simplicity, using timestamp + random – replace with actual UUID.
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	// Implement a simple random string or use crypto/rand.
	// This is a stub.
	return "rand"
}

func getRequestID(c *gin.Context) string {
	if id, exists := c.Get(string(RequestIDKey)); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
