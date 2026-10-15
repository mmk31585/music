package dashboard

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts public maintenance status endpoint.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	status := rg.Group("/status")
	{
		status.GET("/maintenance", h.GetMaintenanceStatus)
	}
}

// RegisterAdminRoutes mounts admin-only dashboard endpoints.
func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	admin := rg.Group("/admin")
	admin.Use(authMW, adminMW)
	{
		admin.GET("/dashboard/stats", h.GetStats)
		admin.POST("/dashboard/maintenance", h.ToggleMaintenance)
	}
}