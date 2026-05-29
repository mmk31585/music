package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
)

// Context keys for Gin
const (
	UserIDContextKey = "auth_user_id"
	UserRoleKey      = "auth_user_role"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// and stores user ID and role in the Gin context.
func AuthMiddleware(tokens *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, ErrUnauthorized())
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, apperrors.Unauthorized("invalid authorization header", nil))
			c.Abort()
			return
		}

		claims, err := tokens.ParseAccessToken(parts[1])
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// Store claims in Gin context
		c.Set(UserIDContextKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole returns a Gin middleware that ensures the authenticated user
// has the required role.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole := UserRoleFromContext(c)
		if currentRole != role {
			response.Error(c, apperrors.Forbidden("insufficient permissions", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserIDFromContext extracts the user ID from the Gin context.
func UserIDFromContext(c *gin.Context) string {
	value, exists := c.Get(UserIDContextKey)
	if !exists {
		return ""
	}
	userID, ok := value.(string)
	if !ok {
		return ""
	}
	return userID
}

// UserRoleFromContext extracts the user role from the Gin context.
func UserRoleFromContext(c *gin.Context) string {
	value, exists := c.Get(UserRoleKey)
	if !exists {
		return ""
	}
	role, ok := value.(string)
	if !ok {
		return ""
	}
	return role
}

// ErrUnauthorized returns a standard unauthorized error.
func ErrUnauthorized() error {
	return apperrors.Unauthorized("authentication required", nil)
}
