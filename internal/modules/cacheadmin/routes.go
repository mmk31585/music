package cacheadmin

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

// RegisterAdminRoutes mounts the cache admin endpoints under the admin group.
func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	admin := rg.Group("/admin")
	admin.Use(authMW, auth.RequireRole("admin"))
	{
		admin.GET("/cache", h.GetStatus)
		admin.POST("/cache/flush", h.FlushCache)
	}
}
