package player

import (
	"context"
	"errors"
	"io"
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
//	@Failure		400	{object}	ErrorResponse	"Invalid Track ID"
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
//	@Failure		400	{object}	ErrorResponse	"Invalid Track ID"
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
//	@Description	Serves the audio file directly instead of redirecting to a storage URL.
//	@Tags			Player
//	@Produce		plain
//	@Param			id	path		string	true	"Track ID"
//	@Success		200	"Audio file content"
//	@Failure		400	{object}	ErrorResponse	"Invalid Track ID"
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

	stream, contentType, err := h.service.GetAudioContent(c.Request.Context(), track.AudioURL)
	if err != nil {
		log.Println("GetAudioContent error:", err)
		h.handleError(c, err)
		return
	}
	defer stream.Close()

	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "public, max-age=3600, immutable")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "Accept-Ranges, Content-Range, Content-Length")
	c.Header("Cross-Origin-Resource-Policy", "cross-origin")
	// Stream directly from S3 to HTTP response without loading entire file into RAM
	io.Copy(c.Writer, stream)
}

// StreamAdminTrack godoc
//
//	@Summary		Stream admin track audio
//	@Description	Serves the audio file directly including private tracks.
//	@Tags			Player Admin
//	@Produce		plain
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Track ID"
//	@Success		200	"Audio file content"
//	@Failure		400	{object}	ErrorResponse	"Invalid Track ID"
//	@Failure		404	{object}	ErrorResponse	"Track or audio not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/admin/player/tracks/{id}/stream [get]
func (h *Handler) StreamAdminTrack(c *gin.Context) {
	track, err := h.service.GetAdminPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}

	stream, contentType, err := h.service.GetAudioContent(c.Request.Context(), track.AudioURL)
	if err != nil {
		log.Println("GetAudioContent error:", err)
		h.handleError(c, err)
		return
	}
	defer stream.Close()

	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "public, max-age=3600, immutable")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "Accept-Ranges, Content-Range, Content-Length")
	c.Header("Cross-Origin-Resource-Policy", "cross-origin")
	// Stream directly from S3 to HTTP response without loading entire file into RAM
	io.Copy(c.Writer, stream)
}

// StreamCoverTrack godoc
//
//	@Summary		Stream track cover image
//	@Description	Proxies the track cover image from storage so it is served
//	@Description	through the backend instead of a direct private-bucket URL.
//	@Tags			Player
//	@Produce		image/jpeg, image/png, image/webp
//	@Param			id	path		string	true	"Track ID"
//	@Success		200	{file}	bytes	"Cover image content"
//	@Failure		400	{object}	ErrorResponse	"Invalid Track ID"
//	@Failure		404	{object}	ErrorResponse	"Track or cover not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/player/tracks/{id}/cover [get]
func (h *Handler) StreamCoverTrack(c *gin.Context) {
	track, err := h.service.GetPlaybackTrack(c.Request.Context(), c.Param("id"))
	if err != nil {
		log.Println("StreamCoverTrack GetPlaybackTrack error:", err)
		h.handleError(c, err)
		return
	}

	if track.CoverURL == nil || *track.CoverURL == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "track has no cover image",
		})
		return
	}

	stream, contentType, err := h.service.GetCoverContent(c.Request.Context(), *track.CoverURL)
	if err != nil {
		log.Println("GetCoverContent error:", err)
		h.handleError(c, err)
		return
	}
	defer stream.Close()

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=3600, immutable")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cross-Origin-Resource-Policy", "cross-origin")
	io.Copy(c.Writer, stream)
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

	case errors.Is(err, context.Canceled):
		c.JSON(499, gin.H{
			"success": false,
			"message": "request cancelled",
		})

	case errors.Is(err, context.DeadlineExceeded):
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"success": false,
			"message": "request timeout",
		})

	default:
		log.Println("player handler unexpected error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "internal server error",
		})
	}
}
