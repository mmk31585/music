package reactions

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	r := rg.Group("/reactions")
	r.Use(authMW)
	{
		r.POST("", handler.React)
		r.DELETE("/:targetType/:targetId", handler.RemoveReaction)
		r.GET("/:targetType/:targetId/mine", handler.GetUserReaction)
		r.GET("/:targetType/:targetId/counts", handler.GetCounts)
		r.GET("/tracks", handler.GetLikedTracks)
		r.GET("/albums", handler.GetLikedAlbums)
	}
}
