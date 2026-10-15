package auth

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"net/http"
)

// Context keys for Gin
const (
	UserIDContextKey      = "auth_user_id"
	UserRoleKey           = "auth_user_role"
	UserHierarchyLevelKey = "auth_user_hierarchy_level"
)

// RoleHierarchyFetcher is a function type for fetching role hierarchy level
type RoleHierarchyFetcher func(ctx context.Context, roleSlug string) (int, error)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// and stores user ID and role in the Gin context.
func AuthMiddleware(tokens *TokenManager, fetchHierarchy ...RoleHierarchyFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, ErrUnauthorized())
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "invalid authorization header", nil))
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

		// Fetch and store role hierarchy level
		var roleLevel int
		if claims.Role != "" && len(fetchHierarchy) > 0 && fetchHierarchy[0] != nil {
			if level, err := fetchHierarchy[0](c.Request.Context(), claims.Role); err == nil {
				roleLevel = level
			}
		}
		c.Set(UserHierarchyLevelKey, roleLevel)

		c.Next()
	}
}

// RequireRole returns a Gin middleware that ensures the authenticated user
// has a role with sufficient hierarchy level (greater than or equal to any of the required roles).
func RequireRoleHierarchical(fetchHierarchy RoleHierarchyFetcher, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current role hierarchy level from context
		currentLevel, _ := c.Get(UserHierarchyLevelKey)
		currentLevelInt, ok := currentLevel.(int)
		if !ok {
			response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", nil))
			c.Abort()
			return
		}

		// If no roles are required, allow access
		if len(roles) == 0 {
			c.Next()
			return
		}

		// Find the minimum hierarchy level among the required roles
		var minRequiredLevel int
		first := true

		for _, roleSlug := range roles {
			if roleSlug == "" {
				continue
			}
			if level, err := fetchHierarchy(c.Request.Context(), roleSlug); err == nil {
				if first || level < minRequiredLevel {
					minRequiredLevel = level
					first = false
				}
			}
			// If we can't find a role, treat it as requiring the highest level (secure by default)
			if first {
				// If this is the first role and we couldn't find it, we can't determine a minimum
				// For security, we'll deny access unless we can validate at least one role
				continue
			}
		}

		// If we couldn't validate any roles, deny access for security
		if first {
			response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", nil))
			c.Abort()
			return
		}

		// Check if current role level meets or exceeds the minimum required level
		if currentLevelInt >= minRequiredLevel {
			c.Next()
			return
		}

		response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", nil))
		c.Abort()
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

// OptionalAuthMiddleware returns a Gin middleware that validates JWT tokens
// but does NOT return 401 on missing/invalid token. Instead it sets a nil user
// in context and calls next(), allowing unauthenticated access.
func OptionalAuthMiddleware(tokens *TokenManager, hierarchyFetcher RoleHierarchyFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Next()
			return
		}

		claims, err := tokens.ParseAccessToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserRoleKey, claims.Role)

		// Fetch and store role hierarchy level
		var roleLevel int
		if claims.Role != "" {
			if level, err := hierarchyFetcher(c.Request.Context(), claims.Role); err == nil {
				roleLevel = level
			}
		}
		c.Set(UserHierarchyLevelKey, roleLevel)

		c.Next()
	}

}

// ErrUnauthorized returns a standard unauthorized error.
func ErrUnauthorized() error {
	return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
}

// RequireRole provides a simpler role (slug-based, equality) guard for module routes
// when a hierarchy fetcher is not available.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		current := UserRoleFromContext(c)
		for _, r := range roles {
			if current == r {
				c.Next()
				return
			}
		}
		response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", nil))
		c.Abort()
	}
}
