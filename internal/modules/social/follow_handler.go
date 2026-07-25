package social

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Follow godoc
// @Summary Follow a user
// @Description Follows another user for the authenticated user.
// @Tags social
// @Produce json
// @Security Bearer
// @Param userId path string true "User ID to follow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/follow/{userId} [post]
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

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollows another user for the authenticated user.
// @Tags social
// @Produce json
// @Security Bearer
// @Param userId path string true "User ID to unfollow"
// @Success 200 {object} map[string]interface{}
// @Router /social/follow/{userId} [delete]
func (h *Handler) Unfollow(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	followedID := c.Param("userId")
	_ = h.service.Unfollow(c.Request.Context(), userID, followedID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetFollowers godoc
// @Summary Get user followers
// @Description Returns paginated followers for a user.
// @Tags social
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} FollowersResponse
// @Failure 500 {object} map[string]interface{}
// @Router /social/followers/{userId} [get]
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

// GetFollowing godoc
// @Summary Get users followed by a user
// @Description Returns paginated list of users a specific user follows.
// @Tags social
// @Produce json
// @Param userId path string false "User ID (defaults to authenticated user)"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} FollowersResponse
// @Failure 500 {object} map[string]interface{}
// @Router /social/following/{userId} [get]
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

// IsFollowing godoc
// @Summary Check if following a user
// @Description Checks whether the authenticated user follows another user.
// @Tags social
// @Produce json
// @Security Bearer
// @Param userId path string true "Target user ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/is-following/{userId} [get]
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

// GetFeed godoc
// @Summary Get activity feed
// @Description Returns the activity feed for the authenticated user or a specified user.
// @Tags social
// @Produce json
// @Security Bearer
// @Param user_id query string false "User ID (defaults to authenticated user)"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Items to skip" default(0)
// @Param types query string false "Filter by activity types (comma-separated)"
// @Success 200 {object} ActivityResponse
// @Failure 500 {object} map[string]interface{}
// @Router /social/feed [get]
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
