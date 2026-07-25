package analytics

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetOverview godoc
// @Summary Get platform analytics overview
// @Description Returns platform-wide stats: total tracks, users, albums, plays, active users, storage.
// @Tags analytics
// @Produce json
// @Success 200 {object} OverviewResponse
// @Failure 500 {object} map[string]interface{}
// @Router /admin/analytics/overview [get]
func (h *Handler) GetOverview(c *gin.Context) {
	ov, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch analytics overview"})
		return
	}

	c.JSON(http.StatusOK, OverviewResponse{
		TotalTracks:    ov.TotalTracks,
		TotalUsers:     ov.TotalUsers,
		TotalAlbums:    ov.TotalAlbums,
		TotalPlays:     ov.TotalPlays,
		ActiveUsers24h: ov.ActiveUsers24h,
		StorageUsedMB:  ov.StorageUsedMB,
	})
}

// TrackEvent godoc
// @Summary Track analytics event
// @Description Records an analytics event for an authenticated or anonymous user.
// @Tags analytics
// @Accept json
// @Produce json
// @Param request body TrackEventRequest true "Analytics event payload"
// @Success 201 {object} TrackEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /analytics/events [post]
func (h *Handler) TrackEvent(c *gin.Context) {
	var req TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userID := web.GetOptionalUserUUID(c)

	err := h.service.TrackEvent(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidEventType):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event type"})
		case errors.Is(err, ErrInvalidTrackID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		case errors.Is(err, ErrInvalidArtistID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		case errors.Is(err, ErrInvalidAlbumID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album id"})
		case errors.Is(err, ErrInvalidPlaylistID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		case errors.Is(err, ErrSearchQueryRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		case errors.Is(err, ErrTrackIDRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "track id is required for this event"})
		case errors.Is(err, ErrTrackNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		case errors.Is(err, ErrArtistNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		case errors.Is(err, ErrAlbumNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		case errors.Is(err, ErrPlaylistNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record analytics event"})
		}
		return
	}

	c.JSON(http.StatusCreated, TrackEventResponse{
		Message: "analytics event recorded",
	})
}
