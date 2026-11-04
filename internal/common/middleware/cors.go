package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{})

	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)

		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Handle CORS preflight before auth middleware: 204 with the right headers.
		if c.Request.Method == http.MethodOptions {
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Vary", "Origin")
				}
				// Never combine a wildcard origin with credentials: only set
				// Allow-Credentials when a specific origin is allowed.
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID, X-Device-ID, Range")
				c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges, Content-Type")
				c.Header("Access-Control-Max-Age", "86400")
				c.Status(http.StatusNoContent)
				return
			}
			c.Next()
			return
		}

		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID, X-Device-ID, Range")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		c.Next()
	}
}
