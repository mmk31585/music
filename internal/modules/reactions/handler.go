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
