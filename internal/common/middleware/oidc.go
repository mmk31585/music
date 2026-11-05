package middleware

import (
	"github.com/gin-gonic/gin"
	"music/internal/platform/storage"
)

// OIDCTokenMiddleware extracts the Vercel OIDC token from the
// x-vercel-oidc-token request header and stores it in the request context
// so the S3 credential provider can use it for AssumeRoleWithWebIdentity.
func OIDCTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-vercel-oidc-token")
		if token != "" {
			c.Request = c.Request.WithContext(storage.WithOIDCToken(c.Request.Context(), token))
		}
		c.Next()
	}
}
