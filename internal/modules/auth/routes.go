package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	rg.GET("/users/:id/profile", handler.GetPublicProfile)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
		auth.POST("/logout", handler.Logout)

		// Protected routes
		protected := auth.Group("/")
		protected.Use(authMW)
		protected.GET("/me", handler.Me)
	}

	// Admin user management
	admin := rg.Group("/admin/users")
	admin.Use(authMW, RequireRole("admin"))
	{
		admin.GET("", handler.AdminListUsers)
		admin.GET("/:id", handler.AdminGetUser)
		admin.PATCH("/:id", handler.AdminUpdateUser)
		admin.DELETE("/:id", handler.AdminDeleteUser)
	}
}
