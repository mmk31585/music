package queue

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	q := rg.Group("/queue")
	q.Use(authMW)

	{
		q.GET("", handler.GetQueue)

		q.POST("/tracks", handler.AddTrack)
		q.DELETE("/tracks/:id", handler.RemoveTrack)

		q.PUT("/reorder", handler.Reorder)

		q.DELETE("", handler.Clear)
	}
}
