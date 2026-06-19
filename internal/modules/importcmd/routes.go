package importcmd

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	admin := rg.Group("/admin/import")
	admin.Use(authMW, auth.RequireRole("admin"))
	{
		admin.GET("/search", h.Search)
		admin.POST("/import", h.Import)
		admin.GET("/:jobId/progress", h.GetProgress)
	}
}
