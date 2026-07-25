package reactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// React godoc
// @Summary Create or update reaction
// @Description Reacts to a target (like/unlike) for the authenticated user.
// @Tags reactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ReactRequest true "Reaction payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions [post]
func (h *Handler) React(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req ReactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.React(c.Request.Context(), userID, req.TargetID, req.TargetType, req.Type); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RemoveReaction godoc
// @Summary Remove reaction
// @Description Removes a reaction from a target for the authenticated user.
// @Tags reactions
// @Produce json
// @Security Bearer
// @Param targetType path string true "Target type (track, album, playlist)"
// @Param targetId path string true "Target ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions/{targetType}/{targetId} [delete]
func (h *Handler) RemoveReaction(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	targetID := c.Param("targetId")
	targetType := c.Param("targetType")
	if err := h.service.RemoveReaction(c.Request.Context(), userID, targetID, targetType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetUserReaction godoc
// @Summary Get user reaction
// @Description Returns the authenticated user's reaction to a specific target.
// @Tags reactions
// @Produce json
// @Security Bearer
// @Param targetType path string true "Target type (track, album, playlist)"
// @Param targetId path string true "Target ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions/{targetType}/{targetId}/mine [get]
func (h *Handler) GetUserReaction(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	targetID := c.Param("targetId")
	targetType := c.Param("targetType")
	reaction, err := h.service.GetUserReaction(c.Request.Context(), userID, targetID, targetType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"reaction": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reaction": reaction})
}

// GetCounts godoc
// @Summary Get reaction counts
// @Description Returns reaction counts for a specific target.
// @Tags reactions
// @Produce json
// @Param targetType path string true "Target type (track, album, playlist)"
// @Param targetId path string true "Target ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions/{targetType}/{targetId}/counts [get]
func (h *Handler) GetCounts(c *gin.Context) {
	targetID := c.Param("targetId")
	targetType := c.Param("targetType")
	counts, err := h.service.GetCounts(c.Request.Context(), targetID, targetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get counts"})
		return
	}
	c.JSON(http.StatusOK, counts)
}

// GetLikedTracks godoc
// @Summary Get liked tracks
// @Description Returns the authenticated user's liked tracks with pagination.
// @Tags reactions
// @Produce json
// @Security Bearer
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions/tracks [get]
func (h *Handler) GetLikedTracks(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetUserLikedTracks(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get liked tracks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset, "count": len(items)})
}

// GetLikedAlbums godoc
// @Summary Get liked albums
// @Description Returns the authenticated user's liked albums with pagination.
// @Tags reactions
// @Produce json
// @Security Bearer
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reactions/albums [get]
func (h *Handler) GetLikedAlbums(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetUserLikedAlbums(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get liked albums"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset, "count": len(items)})
}
