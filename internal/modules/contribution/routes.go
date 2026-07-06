package contribution

import (
	"github.com/gin-gonic/gin"
	"music/internal/modules/auth"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
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
		modProtected.Use(authMW, auth.RequireRole("moderator", "admin"))
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
	adminContributions.Use(authMW, auth.RequireRole("admin"))
	{
		adminContributions.GET("", h.ListByStatus)
		adminContributions.POST("/:id/apply", h.Apply)
	}
}
