package contribution

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc, modMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	contributions := api.Group("/contributions")
	{
		contributions.GET("/leaderboard", h.GetLeaderboard)

		protected := contributions.Group("")
		protected.Use(authMW)
		{
			protected.POST("", h.Create)
			protected.GET("", h.ListMyContributions)
			protected.GET("/:id", h.GetByID)
			protected.GET("/:id/history", h.GetHistory)
		}

		// Moderator/admin review routes
		modProtected := contributions.Group("")
		modProtected.Use(authMW, modMW)
		{
			modProtected.GET("/pending", h.ListPending)
			modProtected.POST("/:id/review", h.Review)
		}
	}

	content := api.Group("/content")
	content.Use(authMW)
	{
		content.GET("/:targetType/:targetID/contributions", h.ListByTarget)
		content.GET("/:targetType/:targetID/versions", h.GetContentVersions)
	}

	// Admin contribution routes
	adminContributions := api.Group("/admin/contributions")
	adminContributions.Use(authMW, adminMW)
	{
		adminContributions.GET("", h.ListByStatus)
		adminContributions.POST("/:id/apply", h.Apply)
	}
}