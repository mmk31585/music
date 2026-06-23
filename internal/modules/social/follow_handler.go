package social

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Follow(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	followedID := c.Param("userId")
	if followedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		return
	}
	if err := h.service.Follow(c.Request.Context(), userID, followedID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Unfollow(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	followedID := c.Param("userId")
	_ = h.service.Unfollow(c.Request.Context(), userID, followedID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetFollowers(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.GetString("auth_user_id")
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.service.GetFollowers(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followers"})
		return
	}
	c.JSON(http.StatusOK, FollowersResponse{Items: items, TotalCount: total, Limit: limit, Offset: offset})
}

func (h *Handler) GetFollowing(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.GetString("auth_user_id")
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.service.GetFollowing(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get following"})
		return
	}
	c.JSON(http.StatusOK, FollowersResponse{Items: items, TotalCount: total, Limit: limit, Offset: offset})
}

func (h *Handler) IsFollowing(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	targetID := c.Param("userId")
	following, err := h.service.IsFollowing(c.Request.Context(), userID, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_following": following})
}

func (h *Handler) GetFeed(c *gin.Context) {
	userID := c.DefaultQuery("user_id", c.GetString("auth_user_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	types := c.Query("types")
	items, err := h.service.GetFeed(c.Request.Context(), userID, limit, offset, types)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get feed"})
		return
	}
	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	c.JSON(http.StatusOK, ActivityResponse{
		Items:      items,
		Pagination: Pagination{Limit: limit, Offset: offset, Count: len(items), HasMore: hasMore},
	})
}

// --- Listening Parties ---
