package playlist

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	playlists := api.Group("/playlists")
	{
		// Public routes
		playlists.GET("", h.ListPublicPlaylists)
		playlists.GET("/:id", h.GetPlaylist)

		// Auth routes
		protected := playlists.Group("")
		protected.Use(authMW)
		{
			protected.POST("", h.CreatePlaylist)
			protected.GET("/me", h.ListMyPlaylists)

			protected.PUT("/:id", h.UpdatePlaylist)
			protected.DELETE("/:id", h.DeletePlaylist)

			protected.POST("/:id/tracks", h.AddTrack)
			protected.DELETE("/:id/tracks/:trackId", h.RemoveTrack)
			protected.PUT("/:id/tracks/reorder", h.ReorderTrack)
		}
	}
}
