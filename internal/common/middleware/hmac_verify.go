package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// VerifyMLServiceWebhook returns a Gin middleware that verifies the HMAC-SHA256
// signature on callbacks from the Python ML service.
//
// The signing scheme is defined in docs/webhook-contract.md (moja-ml-service):
//
//	signed = "{timestamp}.{raw_body}"
//	signature = hex(HMAC-SHA256(secret, signed))
//
// Replay protection: timestamp must be within 300 seconds of the server clock.
// Body is read before any JSON binding and restored for downstream handlers.
func VerifyMLServiceWebhook(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "webhook not configured"})
			return
		}

		timestamp := c.GetHeader("X-Moja-Timestamp")
		signature := c.GetHeader("X-Moja-Signature")

		if timestamp == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing signature headers"})
			return
		}

		// Replay protection: reject timestamps outside 300-second window (±5 min)
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || time.Since(time.Unix(ts, 0)) > 5*time.Minute || time.Since(time.Unix(ts, 0)) < -5*time.Minute {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "stale or invalid timestamp"})
			return
		}

		// Read raw body BEFORE JSON binding and restore it
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		// Recompute signature: HMAC-SHA256("{timestamp}.{raw_body}")
		signedPayload := timestamp + "." + string(rawBody)
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(signedPayload))
		expected := hex.EncodeToString(mac.Sum(nil))

		// Constant-time comparison
		if !hmac.Equal([]byte(signature), []byte(expected)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		c.Next()
	}
}
