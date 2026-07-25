package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/modules/auth"
)

// RequirePermission returns a middleware that checks if the authenticated user
// has the specified permission via their role + overrides.
func RequirePermission(svc ServiceInterface, permSlug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.UserIDFromContext(c)
		if userID == "" {
			response.Error(c, apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil))
			c.Abort()
			return
		}

		ok, err := svc.HasPermission(c.Request.Context(), userID, permSlug)
		if err != nil {
			response.Error(c, apperrors.New(http.StatusInternalServerError, apperrors.CodeInternal, "failed to check permissions", nil))
			c.Abort()
			return
		}

		if !ok {
			response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", map[string]string{
				"required": permSlug,
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission returns a middleware that checks if the authenticated user
// has at least one of the specified permissions.
func RequireAnyPermission(svc ServiceInterface, permSlugs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.UserIDFromContext(c)
		if userID == "" {
			response.Error(c, apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil))
			c.Abort()
			return
		}

		ok, err := svc.HasAnyPermission(c.Request.Context(), userID, permSlugs...)
		if err != nil {
			response.Error(c, apperrors.New(http.StatusInternalServerError, apperrors.CodeInternal, "failed to check permissions", nil))
			c.Abort()
			return
		}

		if !ok {
			response.Error(c, apperrors.New(http.StatusForbidden, apperrors.CodeForbidden, "insufficient permissions", map[string]any{
				"required_any": permSlugs,
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}

// InjectUserAccess is a middleware that resolves and stores the full user access
// record in the Gin context for downstream handlers.
func InjectUserAccess(svc ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.UserIDFromContext(c)
		if userID == "" {
			c.Next()
			return
		}

		access, err := svc.ResolveUserAccess(c.Request.Context(), userID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_access", access)
		c.Next()
	}
}

// UserAccessFromContext extracts the resolved UserAccess from the Gin context.
func UserAccessFromContext(c *gin.Context) (UserAccess, bool) {
	access, exists := c.Get("user_access")
	if !exists {
		return UserAccess{}, false
	}
	ua, ok := access.(UserAccess)
	return ua, ok
}
