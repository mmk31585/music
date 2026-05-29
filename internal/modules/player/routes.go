package player

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(router *gin.RouterGroup, handler *Handler) {
	playerGroup := router.Group("/player")
	{
		playerGroup.GET("/tracks/:id", handler.GetPlaybackTrack)
		playerGroup.GET("/tracks/:id/stream", handler.StreamTrack)
	}
}

func RegisterPrivateRoutes(router *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	playerGroup := router.Group("/player")
	playerGroup.Use(authMW)
	{
		// Add authenticated player endpoints later.
		// Example:
		// playerGroup.POST("/tracks/:id/play", handler.TrackPlayed)
	}
}
