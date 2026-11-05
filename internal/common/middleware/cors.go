package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CORS(allowedOrigins []string, logger *zap.Logger) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)
		origin = strings.TrimRight(origin, "/")
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	fallbackHeaders := "Authorization,Content-Type,Accept,X-Request-ID,X-Vercel-Oidc-Token"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		isPreflight := c.Request.Method == http.MethodOptions

		if isPreflight {
			requestedHeaders := c.GetHeader("Access-Control-Request-Headers")
			requestedMethod := c.GetHeader("Access-Control-Request-Method")

			matched := false
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					matched = true
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Credentials", "true")
					c.Header("Vary", "Origin, Access-Control-Request-Headers")
					c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
					if requestedHeaders != "" {
						c.Header("Access-Control-Allow-Headers", requestedHeaders)
					} else {
						c.Header("Access-Control-Allow-Headers", fallbackHeaders)
					}
					c.Header("Access-Control-Max-Age", "600")
					c.Status(http.StatusNoContent)
				}
			}
			if logger != nil {
				logger.Info("CORS preflight",
					zap.String("origin", origin),
					zap.Bool("matched", matched),
					zap.String("requested_headers", requestedHeaders),
					zap.String("requested_method", requestedMethod),
				)
			}
			if !matched {
				c.AbortWithStatus(http.StatusNoContent)
			}
			return
		}

		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Vary", "Origin")
			}
		}

		c.Next()

		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Vary", "Origin")
			}
		}
	}
}
