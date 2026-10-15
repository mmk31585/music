package analytics

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, optionalAuthMW ...gin.HandlerFunc) {
	analyticsGroup := rg.Group("/analytics")

	for _, mw := range optionalAuthMW {
		analyticsGroup.Use(mw)
	}

	{
		analyticsGroup.POST("/events", handler.TrackEvent)
	}
}

// RegisterAdminRoutes mounts the analytics admin endpoints under the admin group.
func RegisterAdminRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	admin := rg.Group("/admin")
	admin.Use(authMW, adminMW)
	{
		admin.GET("/analytics/overview", handler.GetOverview)
	}
}