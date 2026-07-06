package moderation

import (
	"github.com/gin-gonic/gin"
	"music/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	r := rg.Group("/moderation")
	{
		r.POST("/report", authMW, handler.Report)
		r.GET("/pending", authMW, auth.RequireRole("moderator", "admin"), handler.ListPending)
		r.GET("/status/:status", authMW, auth.RequireRole("moderator", "admin"), handler.ListByStatus)
		r.GET("/reports/:id", authMW, handler.GetReport)
		r.POST("/resolve/:id", authMW, auth.RequireRole("moderator", "admin"), handler.Resolve)
		r.POST("/bulk", authMW, auth.RequireRole("moderator", "admin"), handler.BulkAction)
		r.POST("/flag", authMW, handler.FlagContent)
		r.GET("/flags", authMW, handler.ListFlags)
		r.GET("/stats", authMW, auth.RequireRole("moderator", "admin"), handler.GetStats)
		r.GET("/actions", authMW, auth.RequireRole("moderator", "admin"), handler.GetActions)
	}
}
