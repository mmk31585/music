package health

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	rg.GET("/health", handler.Health)
	rg.GET("/live", handler.Live)
	rg.GET("/ready", handler.Ready)
}
