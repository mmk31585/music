package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"music/internal/common/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc, rdb *redis.Client) {
	rg.GET("/users/:id/profile", handler.GetPublicProfile)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimitIP(rdb, 3, time.Minute), handler.Register)
		auth.POST("/login", middleware.RateLimitIP(rdb, 5, time.Minute), handler.Login)
	auth.POST("/forgot-password", middleware.RateLimitIP(rdb, 3, time.Minute), handler.ForgotPassword)
		auth.POST("/refresh", middleware.RateLimitIP(rdb, 5, time.Minute), handler.Refresh)
		auth.POST("/logout", authMW, handler.Logout)

		// Protected routes
		protected := auth.Group("/")
		protected.Use(authMW)
		protected.GET("/me", handler.Me)
	}

	// User profile (authenticated)
	users := rg.Group("/users")
	users.Use(authMW)
	{
		users.PUT("/me/profile", handler.UpdateProfile)
		users.PUT("/me/password", handler.ChangePassword)
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
