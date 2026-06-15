package dashboard

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	admin := rg.Group("/admin/dashboard")
	admin.Use(authMW, auth.RequireRole("admin"))
	{
		admin.GET("/stats", h.GetStats)
	}
}
