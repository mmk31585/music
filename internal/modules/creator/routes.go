package creator

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	creator := rg.Group("/creator")
	creator.Use(authMW)
	{
		creator.GET("/overview", handler.GetOverview)
		creator.GET("/daily", handler.GetDailyStats)
		creator.GET("/tracks", handler.GetTrackStats)
		creator.POST("/refresh", handler.RefreshStats)
		creator.GET("/check", handler.IsCreator)

		// Earnings
		creator.GET("/earnings", handler.GetEarnings)
		creator.GET("/earnings/payouts", handler.GetPayoutHistory)
		creator.GET("/earnings/methods", handler.GetPayoutMethods)

		// Audience
		creator.GET("/audience", handler.GetAudience)

		// Content Management
		creator.GET("/content", handler.GetContent)
		creator.PUT("/tracks/:trackId", handler.UpdateTrack)
		creator.PUT("/albums/:albumId", handler.UpdateAlbum)
		creator.DELETE("/tracks/:trackId", handler.DeleteTrack)
		creator.DELETE("/albums/:albumId", handler.DeleteAlbum)
	}
}
