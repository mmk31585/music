package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

const csrfCookieName = "_csrf_token"
const csrfHeaderName = "X-CSRF-Token"

var csrfKey = func() []byte {
	b := make([]byte, 32)
	rand.Read(b)
	return b
}()

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			setCSRFCookie(c)
			c.Next()
			return
		}

		token := c.GetHeader(csrfHeaderName)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "CSRF_MISSING",
				"message": "missing CSRF token",
			})
			return
		}

		cookie, err := c.Cookie(csrfCookieName)
		if err != nil || !validateCSRF(token, cookie) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "CSRF_INVALID",
				"message": "invalid CSRF token",
			})
			return
		}

		c.Next()
	}
}

func setCSRFCookie(c *gin.Context) {
	token := generateCSRFToken()
	c.SetCookie(csrfCookieName, token, 3600, "/", "", false, true)
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	mac := hmac.New(sha256.New, csrfKey)
	mac.Write(b)
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(append(b, sig...))
}

func validateCSRF(token, cookie string) bool {
	if token == "" || cookie == "" {
		return false
	}
	return hmac.Equal([]byte(token), []byte(cookie))
}
