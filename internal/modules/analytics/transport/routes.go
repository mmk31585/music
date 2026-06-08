package analytics

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, optionalAuthMW ...gin.HandlerFunc) {
	analyticsGroup := rg.Group("/analytics")

	for _, mw := range optionalAuthMW {
		analyticsGroup.Use(mw)
	}

	{
		analyticsGroup.POST("/events", handler.TrackEvent)
	}
}
