package library

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	library := rg.Group("/library")
	library.Use(authMW)
	{
		library.GET("/tracks", handler.ListLikedTracks)
		library.POST("/tracks/like", handler.LikeTrack)
		library.DELETE("/tracks/:trackId/like", handler.UnlikeTrack)

		library.GET("/albums", handler.ListLikedAlbums)
		library.POST("/albums/like", handler.LikeAlbum)
		library.DELETE("/albums/:albumId/like", handler.UnlikeAlbum)

		library.GET("/artists", handler.ListFollowedArtists)
		library.POST("/artists/follow", handler.FollowArtist)
		library.DELETE("/artists/:artistId/follow", handler.UnfollowArtist)

		library.POST("/history", handler.AddPlayHistory)
		library.GET("/history", handler.ListPlayHistory)
		library.GET("/recently-played", handler.ListRecentlyPlayed)
	}
}
