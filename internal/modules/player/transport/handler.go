package player

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type PlaybackResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetPlaybackTrack godoc
//
//	@Summary		Get public track playback information
//	@Description	Returns public playback metadata for a track
//	@Tags			Player
//	@Produce		json
//	@Param			id	path		string	true	"Track ID"
//	@Success		200	{object}	PlaybackResponse
//	@Failure		400	{object}	ErrorResponse	"Invalid track ID"
//	@Failure		403	{object}	ErrorResponse	"Track is private"
//	@Failure		404	{object}	ErrorResponse	"Track not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/player/tracks/{id} [get]
func (h *Handler) GetPlaybackTrack(c *gin.Context) {
	track, err := h.service.GetPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.service.BuildPlaybackResponse(track),
	})
}

// GetAdminPlaybackTrack godoc
//
//	@Summary		Get admin playback track
//	@Description	Returns playback metadata including private tracks
//	@Tags			Player Admin
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Track ID"
//	@Success		200	{object}	PlaybackResponse
//	@Failure		400	{object}	ErrorResponse	"Invalid track ID"
//	@Failure		404	{object}	ErrorResponse	"Track not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/admin/player/tracks/{id} [get]
func (h *Handler) GetAdminPlaybackTrack(c *gin.Context) {
	track, err := h.service.GetAdminPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.service.BuildAdminPlaybackResponse(track),
	})
}

// StreamTrack godoc
//
//	@Summary		Stream track audio
//	@Description	Redirects client to the actual audio file URL
//	@Tags			Player
//	@Produce		plain
//	@Param			id	path	string	true	"Track ID"
//	@Success		302	"Redirect to audio file"
//	@Failure		400	{object}	ErrorResponse	"Invalid track ID"
//	@Failure		403	{object}	ErrorResponse	"Track is private"
//	@Failure		404	{object}	ErrorResponse	"Track or audio not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/player/tracks/{id}/stream [get]
func (h *Handler) StreamTrack(c *gin.Context) {
	log.Println("StreamTrack handler hit, id:", c.Param("id"))

	track, err := h.service.GetPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		log.Println("GetPlaybackTrack error:", err)
		h.handleError(c, err)
		return
	}

	log.Println("StreamTrack track audio_url:", track.AudioURL)

	audioURL, err := h.service.ResolveAudioURL(c.Request.Context(), track.AudioURL)
	if err != nil {
		log.Println("ResolveAudioURL error:", err)
		h.handleError(c, err)
		return
	}

	log.Println("StreamTrack redirecting to:", audioURL)

	c.Redirect(http.StatusFound, audioURL)
}

// StreamAdminTrack godoc
//
//	@Summary		Stream admin track audio
//	@Description	Redirects client to the actual audio file URL including private tracks
//	@Tags			Player Admin
//	@Produce		plain
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Track ID"
//	@Success		302	"Redirect to audio file"
//	@Failure		400	{object}	ErrorResponse	"Invalid track ID"
//	@Failure		404	{object}	ErrorResponse	"Track or audio not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/admin/player/tracks/{id}/stream [get]
func (h *Handler) StreamAdminTrack(c *gin.Context) {
	track, err := h.service.GetAdminPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}

	audioURL, err := h.service.ResolveAudioURL(c.Request.Context(), track.AudioURL)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Redirect(http.StatusFound, audioURL)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidTrackID):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid track id",
		})

	case errors.Is(err, ErrTrackNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "track not found",
		})

	case errors.Is(err, ErrAudioNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "track audio not found",
		})

	case errors.Is(err, ErrPrivateTrack):
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "track is private",
		})

	case errors.Is(err, ErrInvalidMediaURL):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid media url",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "internal server error",
		})
	}
}
