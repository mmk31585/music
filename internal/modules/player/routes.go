package player

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	playerGroup := router.Group("/player")
	{
		playerGroup.GET("/tracks/:id", handler.GetPlaybackTrack)
		playerGroup.GET("/tracks/:id/stream", handler.StreamTrack)
	}
}
func RegisterPrivateRoutes(router *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	playerGroup := router.Group("/admin/player")
	playerGroup.Use(authMW, auth.RequireRole("admin"))
	{
		playerGroup.GET("/tracks/:id", handler.GetAdminPlaybackTrack)
		playerGroup.GET("/tracks/:id/stream", handler.StreamAdminTrack)
	}
}
