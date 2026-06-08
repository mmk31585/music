package lyrics

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	// Public routes
	api.GET("/tracks/:trackId/lyrics", h.GetTrackLyrics)
	api.GET("/tracks/:trackId/lyrics/all", h.GetLyricsByTrackID)

	// Admin routes (protected)
	admin := api.Group("/admin")
	admin.Use(authMW)
	{
		admin.POST("/lyrics", h.CreateLyrics)
		admin.PUT("/lyrics/:id", h.UpdateLyrics)
		admin.DELETE("/lyrics/:id", h.DeleteLyrics)
	}
}
