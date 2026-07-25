package gamification

import (
	"github.com/gin-gonic/gin"
	"music/internal/modules/auth"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	gamification := api.Group("/gamification")
	gamification.Use(authMW)
	{
		gamification.GET("/profile", h.GetProfile)
		gamification.GET("/badges", h.GetBadges)
		gamification.GET("/challenges", h.GetChallenges)
		gamification.GET("/leaderboard", h.GetLeaderboard)
		gamification.POST("/badges/check", h.CheckBadges)
	}

	// XP management (admin only)
	admin := api.Group("/admin/gamification")
	admin.Use(authMW, auth.RequireRole("admin"))
	{
		admin.POST("/xp", h.AddXP)
	}
}
