package history

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	historyGroup := rg.Group("/history")
	historyGroup.Use(authMW)

	{
		historyGroup.GET("", handler.GetHistory)
		historyGroup.DELETE("/:id", handler.DeleteHistoryItem)
		historyGroup.DELETE("", handler.ClearHistory)
	}
}
