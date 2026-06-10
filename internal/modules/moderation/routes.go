package moderation

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	r := rg.Group("/moderation")
	{
		r.POST("/report", handler.Report)
		r.GET("/pending", authMW, handler.ListPending)
		r.GET("/status/:status", authMW, handler.ListByStatus)
		r.GET("/reports/:id", authMW, handler.GetReport)
		r.POST("/resolve/:id", authMW, handler.Resolve)
		r.POST("/bulk", authMW, handler.BulkAction)
		r.POST("/flag", authMW, handler.FlagContent)
		r.GET("/flags", authMW, handler.ListFlags)
		r.GET("/stats", authMW, handler.GetStats)
		r.GET("/actions", authMW, handler.GetActions)
	}
}
