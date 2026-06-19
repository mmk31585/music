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

		// Artist discography search
		admin.GET("/artist", h.SearchArtist)

		// Batch import
		admin.POST("/batch", h.BatchImport)
		admin.GET("/batch/:batchId/progress", h.GetBatchProgress)

		// Per-job progress (must be last due to :jobId catch-all)
		admin.GET("/:jobId/progress", h.GetProgress)
	}
}
