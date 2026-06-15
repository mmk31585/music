package media

import (
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	admin := rg.Group("/admin/media")
	admin.Use(authMW)
	{
		admin.GET("", h.ListAdminMedia)
		admin.POST("/upload", h.UploadAdminMedia)
		admin.DELETE("/:id", h.DeleteAdminMedia)
	}
}
