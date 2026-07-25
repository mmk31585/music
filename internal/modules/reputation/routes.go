package reputation

import (
	"github.com/gin-gonic/gin"

	"music/internal/common/response"
	"music/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, tokenManager *auth.TokenManager) {
	api := rg.Group("")
	api.Use(auth.AuthMiddleware(tokenManager))

	// Public reputation queries (any authenticated user)
	api.GET("/reputation/tiers", response.Wrap(handler.GetTrustTiers))
	api.GET("/reputation/top", response.Wrap(handler.GetTopContributors))
	api.GET("/reputation/users/:userId", response.Wrap(handler.GetUserReputation))
}
