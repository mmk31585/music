package recommendation

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler, authMw gin.HandlerFunc) {
	reco := r.Group("/recommendations")
	{
		reco.GET("/popular", handler.PopularTracks)
		reco.GET("/best", handler.BestTracks)
		reco.GET("/similar/:trackId", handler.SimilarTracks)
		reco.GET("/artist/:artistId", handler.TracksByArtist)
		reco.GET("/genre/:genre", handler.TracksByGenre)

		recoAuth := reco.Group("")
		recoAuth.Use(authMw)
		{
			recoAuth.GET("/recent", handler.RecentTracks)
			recoAuth.GET("/for-you", handler.ForYou)
		}
	}
}
