package recommendation

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler, onboardingHandler *OnboardingHandler, authMw gin.HandlerFunc) {
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
			recoAuth.GET("/home", handler.HomeFeed)
			recoAuth.GET("/personalized", handler.PersonalizedTracks)
			recoAuth.GET("/discover-weekly", handler.DiscoverWeekly)
			recoAuth.GET("/stats", handler.ListeningStats)
		}
	}

	// Radio mode routes — require auth
	radio := r.Group("/radio")
	radio.Use(authMw)
	{
		radio.POST("/start", handler.StartRadio)
		radio.GET("/:sessionId/next", handler.GetNextRadioBatch)
		radio.POST("/:sessionId/end", handler.EndRadio)
	}

	// Onboarding routes — require auth
	onboarding := r.Group("/onboarding")
	onboarding.Use(authMw)
	{
		onboarding.POST("/genres", onboardingHandler.SetOnboardingGenres)
	}
}
