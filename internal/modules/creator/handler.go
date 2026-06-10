package creator

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

func (h *Handler) GetOverview(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	stats, err := h.service.GetOverview(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *Handler) GetDailyStats(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	from := c.Query("from")
	to := c.Query("to")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	items, err := h.service.GetDailyStats(c.Request.Context(), userID, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get daily stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetTrackStats(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	items, err := h.service.GetTrackStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get track stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) RefreshStats(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	if err := h.service.RefreshStats(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to refresh"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) IsCreator(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	isCreator, err := h.service.IsCreator(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_creator": isCreator})
}

func (h *Handler) GetEarnings(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	data, err := h.service.GetEarningsBreakdown(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get earnings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *Handler) GetPayoutHistory(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.service.GetPayoutHistory(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get payout history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetPayoutMethods(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	items, err := h.service.GetPayoutMethods(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get payout methods"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetAudience(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("top_limit", "20"))

	overview, err := h.service.GetAudienceOverview(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get audience overview"})
		return
	}

	listeners, err := h.service.GetTopListeners(c.Request.Context(), userID, limit)
	if err != nil {
		listeners = []TopListener{}
	}

	geo, err := h.service.GetGeographicStats(c.Request.Context(), userID)
	if err != nil {
		geo = []GeographicStat{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"overview":      overview,
			"top_listeners": listeners,
			"geographics":   geo,
		},
	})
}

func (h *Handler) GetContent(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	data, err := h.service.GetCreatorContent(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get content"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *Handler) UpdateTrack(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	trackID := c.Param("trackId")
	var req TrackUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateTrack(c.Request.Context(), userID, trackID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) UpdateAlbum(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	albumID := c.Param("albumId")
	var req AlbumUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateAlbum(c.Request.Context(), userID, albumID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteTrack(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	trackID := c.Param("trackId")
	if err := h.service.DeleteTrack(c.Request.Context(), userID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteAlbum(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	albumID := c.Param("albumId")
	if err := h.service.DeleteAlbum(c.Request.Context(), userID, albumID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
