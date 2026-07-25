package playlist

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW, optionalAuthMW gin.HandlerFunc) {
	playlists := api.Group("/playlists")
	{
		// Public routes (with optional auth so logged-in users can see their own private playlists)
		playlists.GET("", h.ListPublicPlaylists)
		playlists.GET("/:id", optionalAuthMW, h.GetPlaylist)

		// Auth routes
		protected := playlists.Group("")
		protected.Use(authMW)
		{
			protected.POST("", h.CreatePlaylist)
			protected.GET("/me", h.ListMyPlaylists)

			protected.PUT("/:id", h.UpdatePlaylist)
			protected.DELETE("/:id", h.DeletePlaylist)
			protected.PUT("/:id/collaborative", h.SetCollaborative)

			protected.GET("/:id/collaborators", h.ListCollaborators)
			protected.POST("/:id/tracks", h.AddTrack)
			protected.DELETE("/:id/tracks/:trackId", h.RemoveTrack)
			protected.PUT("/:id/tracks/reorder", h.ReorderTrack)
		}
	}
}
