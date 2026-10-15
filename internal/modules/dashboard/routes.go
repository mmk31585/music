package dashboard

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	status := rg.Group("/status")
	{
		status.GET("/maintenance", h.GetMaintenanceStatus)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	admin := rg.Group("/admin")
	admin.Use(authMW, adminMW)
	{
		admin.GET("/dashboard/stats", h.GetStats)
		admin.POST("/dashboard/maintenance", h.ToggleMaintenance)
	}
}
