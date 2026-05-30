package follow

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	rg.POST("/artists/:id/follow", authMW, handler.FollowArtist)
	rg.DELETE("/artists/:id/follow", authMW, handler.UnfollowArtist)

	rg.POST("/users/:id/follow", authMW, handler.FollowUser)
	rg.DELETE("/users/:id/follow", authMW, handler.UnfollowUser)
}
