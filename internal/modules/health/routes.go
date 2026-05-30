package health

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	health := rg.Group("/health")
	{
		health.GET("", handler.Health)
		health.GET("/live", handler.Live)
		health.GET("/ready", handler.Ready)
	}
}
