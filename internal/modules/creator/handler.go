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

// GetOverview godoc
// @Summary Get creator overview
// @Description Returns an overview of the authenticated creator's stats.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/overview [get]
func (h *Handler) GetOverview(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	stats, err := h.service.GetOverview(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// GetDailyStats godoc
// @Summary Get creator daily stats
// @Description Returns daily statistics for the authenticated creator within a date range.
// @Tags creator
// @Produce json
// @Security Bearer
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Number of days" default(30)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/daily [get]
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

// GetTrackStats godoc
// @Summary Get creator track stats
// @Description Returns per-track statistics for the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/tracks [get]
func (h *Handler) GetTrackStats(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	items, err := h.service.GetTrackStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get track stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// RefreshStats godoc
// @Summary Refresh creator stats
// @Description Triggers a refresh of the authenticated creator's analytics stats.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/refresh [post]
func (h *Handler) RefreshStats(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	if err := h.service.RefreshStats(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to refresh"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// IsCreator godoc
// @Summary Check if user is a creator
// @Description Returns whether the authenticated user has creator status.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/check [get]
func (h *Handler) IsCreator(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	isCreator, err := h.service.IsCreator(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_creator": isCreator})
}

// GetEarnings godoc
// @Summary Get creator earnings
// @Description Returns the earnings breakdown for the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/earnings [get]
func (h *Handler) GetEarnings(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	data, err := h.service.GetEarningsBreakdown(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get earnings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GetPayoutHistory godoc
// @Summary Get payout history
// @Description Returns payout history for the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/earnings/payouts [get]
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

// GetPayoutMethods godoc
// @Summary Get payout methods
// @Description Returns available payout methods for the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/earnings/methods [get]
func (h *Handler) GetPayoutMethods(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	items, err := h.service.GetPayoutMethods(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get payout methods"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// GetAudience godoc
// @Summary Get creator audience
// @Description Returns audience overview, top listeners, and geographic stats for the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Param top_limit query int false "Number of top listeners" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/audience [get]
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

// GetContent godoc
// @Summary Get creator content
// @Description Returns the authenticated creator's content (tracks and albums).
// @Tags creator
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/content [get]
func (h *Handler) GetContent(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	data, err := h.service.GetCreatorContent(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get content"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// UpdateTrack godoc
// @Summary Update creator track
// @Description Updates a track owned by the authenticated creator.
// @Tags creator
// @Accept json
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Param request body TrackUpdateRequest true "Track update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/tracks/{trackId} [put]
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

// UpdateAlbum godoc
// @Summary Update creator album
// @Description Updates an album owned by the authenticated creator.
// @Tags creator
// @Accept json
// @Produce json
// @Security Bearer
// @Param albumId path string true "Album ID"
// @Param request body AlbumUpdateRequest true "Album update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/albums/{albumId} [put]
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

// DeleteTrack godoc
// @Summary Delete creator track
// @Description Deletes a track owned by the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/tracks/{trackId} [delete]
func (h *Handler) DeleteTrack(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	trackID := c.Param("trackId")
	if err := h.service.DeleteTrack(c.Request.Context(), userID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteAlbum godoc
// @Summary Delete creator album
// @Description Deletes an album owned by the authenticated creator.
// @Tags creator
// @Produce json
// @Security Bearer
// @Param albumId path string true "Album ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /creator/albums/{albumId} [delete]
func (h *Handler) DeleteAlbum(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	albumID := c.Param("albumId")
	if err := h.service.DeleteAlbum(c.Request.Context(), userID, albumID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
