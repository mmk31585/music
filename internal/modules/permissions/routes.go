package permissions

import (
	"github.com/gin-gonic/gin"

	"music/internal/common/response"
	"music/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authSvc *auth.TokenManager, permSvc ServiceInterface) {
	// Public: none (all endpoints require auth)

	// Authenticated group
	api := rg.Group("")
	api.Use(auth.AuthMiddleware(authSvc))

	// --- Current user's access ---
	api.GET("/permissions/me", response.Wrap(handler.GetMyAccess))

	// --- Permissions catalog (admin) ---
	admin := api.Group("/permissions")
	admin.Use(auth.RequireRole("admin"))
	{
		admin.GET("", response.Wrap(handler.ListPermissions))
	}

	// --- Roles (admin) ---
	roles := api.Group("/roles")
	adminRoles := roles
	adminRoles.Use(auth.RequireRole("admin"))
	{
		adminRoles.GET("", response.Wrap(handler.ListRoles))
		adminRoles.POST("", response.Wrap(handler.CreateRole))
		adminRoles.GET("/:slug", response.Wrap(handler.GetRole))
		adminRoles.PUT("/:slug", response.Wrap(handler.UpdateRole))
		adminRoles.DELETE("/:slug", response.Wrap(handler.DeleteRole))
		adminRoles.PUT("/:slug/permissions", response.Wrap(handler.SetRolePermissions))
	}

	// --- User role management (admin) ---
	users := api.Group("/users")
	adminUsers := users
	adminUsers.Use(auth.RequireRole("admin"))
	{
		adminUsers.GET("/:id/access", response.Wrap(handler.GetUserAccess))
		adminUsers.POST("/:id/access/check", response.Wrap(handler.CheckPermission))
		adminUsers.POST("/:id/roles", response.Wrap(handler.AssignUserRole))
		adminUsers.DELETE("/:id/roles/:roleSlug", response.Wrap(handler.RemoveUserRole))
		adminUsers.GET("/:id/overrides", response.Wrap(handler.GetUserOverrides))
		adminUsers.POST("/:id/overrides/grant", response.Wrap(handler.GrantUserPermission))
		adminUsers.POST("/:id/overrides/revoke", response.Wrap(handler.RevokeUserPermission))
	}

	// --- Middleware for downstream use ---
	// Register the InjectUserAccess middleware on all authenticated routes
	// This is done at the app level; see internal/app/routes.go
}
