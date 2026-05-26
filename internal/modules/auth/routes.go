package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
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
}
