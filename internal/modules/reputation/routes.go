package reputation

import (
	"github.com/gin-gonic/gin"

	"music/internal/common/response"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	api := rg.Group("")
	api.Use(authMW)

	// Public reputation queries (any authenticated user)
	api.GET("/reputation/tiers", response.Wrap(handler.GetTrustTiers))
	api.GET("/reputation/top", response.Wrap(handler.GetTopContributors))
	api.GET("/reputation/users/:userId", response.Wrap(handler.GetUserReputation))
}