package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheControl(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			val := fmt.Sprintf("public, max-age=%d, s-maxage=%d", int(maxAge.Seconds()), int(maxAge.Seconds()*2))
			c.Header("Cache-Control", val)
		}
	}
}

func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
