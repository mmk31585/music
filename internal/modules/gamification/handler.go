package gamification

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

func (h *Handler) AddXP(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	var req XPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	newBalance, err := h.service.AddXP(c.Request.Context(), userID, req.Amount, req.Source)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"new_balance": newBalance}})
}

func (h *Handler) GetBadges(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	allBadges, userBadges, err := h.service.GetBadges(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"all":  allBadges,
			"mine": userBadges,
		},
	})
}

func (h *Handler) GetChallenges(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	challenges, userChallenges, err := h.service.GetActiveChallenges(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	earnedXP := 0
	for _, uc := range userChallenges {
		if uc.IsCompleted && uc.Challenge != nil {
			earnedXP += uc.Challenge.XPReward
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"challenges": challenges,
			"progress":   userChallenges,
			"earned_xp":  earnedXP,
		},
	})
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	lbType := c.DefaultQuery("type", "all")
	limit := 20
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	result, err := h.service.GetLeaderboard(c.Request.Context(), lbType, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *Handler) CheckBadges(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	newBadges, err := h.service.CheckBadges(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": newBadges})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "user not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}
