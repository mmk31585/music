package notification

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	notifications := rg.Group("/notifications")
	notifications.Use(authMW)
	{
		notifications.GET("", handler.ListNotifications)
		notifications.POST("/:id/read", handler.MarkAsRead)
		notifications.POST("/read-all", handler.MarkAllAsRead)
	}
}
