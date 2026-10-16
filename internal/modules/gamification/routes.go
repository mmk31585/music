package gamification

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
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
	admin.Use(authMW, adminMW)
	{
		admin.POST("/xp", h.AddXP)
	}
}
