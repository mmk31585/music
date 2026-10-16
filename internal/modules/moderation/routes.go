package moderation

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc, modMW gin.HandlerFunc) {
	r := rg.Group("/moderation")
	{
		r.POST("/report", authMW, handler.Report)
		r.GET("/pending", authMW, modMW, handler.ListPending)
		r.GET("/status/:status", authMW, modMW, handler.ListByStatus)
		r.GET("/reports/:id", authMW, handler.GetReport)
		r.POST("/resolve/:id", authMW, modMW, handler.Resolve)
		r.POST("/bulk", authMW, modMW, handler.BulkAction)
		r.POST("/flag", authMW, handler.FlagContent)
		r.GET("/flags", authMW, handler.ListFlags)
		r.GET("/stats", authMW, modMW, handler.GetStats)
		r.GET("/actions", authMW, modMW, handler.GetActions)
	}
}