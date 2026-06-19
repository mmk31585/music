package lyrics

import (
	"errors"
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

// FetchFromLRC godoc
// @Summary Fetch LRC lyrics from LRCLIB
// @Description Fetches synced lyrics from LRCLIB for a track and saves them.
// @Tags lyrics
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/fetch/{trackId} [post]
func (h *Handler) FetchFromLRC(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}

	lyrics, err := h.service.FetchFromLRC(c.Request.Context(), trackID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ToLyricsResponse(lyrics),
	})
}

// CreateLyrics godoc
// @Summary Create lyrics
// @Description Creates new lyrics for a track.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateLyricsRequest true "Create lyrics request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics [post]
func (h *Handler) CreateLyrics(c *gin.Context) {
	var req CreateLyricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	lyrics, err := h.service.CreateLyrics(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ToLyricsResponse(lyrics),
	})
}

// UpdateLyrics godoc
// @Summary Update lyrics
// @Description Updates existing lyrics.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Lyrics ID"
// @Param request body UpdateLyricsRequest true "Update lyrics request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/{id} [put]
func (h *Handler) UpdateLyrics(c *gin.Context) {
	lyricsID := c.Param("id")
	if lyricsID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "lyrics id is required"})
		return
	}

	var req UpdateLyricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	lyrics, err := h.service.UpdateLyrics(c.Request.Context(), lyricsID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ToLyricsResponse(lyrics),
	})
}

// DeleteLyrics godoc
// @Summary Delete lyrics
// @Description Deletes existing lyrics.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Lyrics ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/{id} [delete]
func (h *Handler) DeleteLyrics(c *gin.Context) {
	lyricsID := c.Param("id")
	if lyricsID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "lyrics id is required"})
		return
	}

	if err := h.service.DeleteLyrics(c.Request.Context(), lyricsID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "lyrics deleted",
	})
}

// GetTrackLyrics godoc
// @Summary Get track lyrics
// @Description Returns lyrics for a specific track.
// @Tags lyrics
// @Accept json
// @Produce json
// @Param trackId path string true "Track ID"
// @Param lang query string false "Language code (e.g., en, es)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tracks/{trackId}/lyrics [get]
func (h *Handler) GetTrackLyrics(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}

	language := c.Query("lang")

	lyrics, err := h.service.GetTrackLyrics(c.Request.Context(), trackID, language)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    lyrics,
	})
}

// GetLyricsByTrackID godoc
// @Summary Get lyrics by track ID
// @Description Returns all lyrics for a specific track.
// @Tags lyrics
// @Accept json
// @Produce json
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tracks/{trackId}/lyrics/all [get]
func (h *Handler) GetLyricsByTrackID(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}

	lyricsList, err := h.service.GetLyricsByTrackID(c.Request.Context(), trackID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    lyricsList,
	})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLyricsNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrLyricsAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidLyricsType):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidLyricsLanguage):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrForbiddenLyricsAccess):
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}

func parseInt64Param(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

// Replace this with your real auth helper.
func getUserIDFromGin(c *gin.Context) (int64, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := v.(int64)
	return id, ok
}
