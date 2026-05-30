package player

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	appErrors "music/internal/common/errors"
	"music/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetPlaybackTrack godoc
// @Summary      Get track playback details
// @Description  Returns track info and a streaming URL
// @Tags         player
// @Param        id   path      string  true  "Track ID"
// @Success      200  {object}  response.Data{data=PlaybackTrackResponse}
// @Failure      404  {object}  errors.AppError
// @Router       /player/tracks/{id} [get]
func (h *Handler) GetPlaybackTrack(c *gin.Context) {
	trackID := c.Param("id")

	track, err := h.service.GetPlaybackTrack(c.Request.Context(), trackID)
	if err != nil {
		log.Print("Get Playback Track Error")
		h.handleError(c, err)
		return
	}

	payload := h.service.BuildPlaybackResponse(track)

	response.Success(c, http.StatusOK, "Success", payload)
}

// StreamTrack godoc
// @Summary      Stream track audio
// @Description  Serves the actual audio file
// @Tags         player
// @Param        id   path      string  true  "Track ID"
// @Produce      audio/mpeg
// @Success      200  {file}    binary
// @Router       /player/tracks/{id}/stream [get]
func (h *Handler) StreamTrack(c *gin.Context) {
	trackID := c.Param("id")
	track, err := h.service.GetPlaybackTrack(c.Request.Context(), trackID)
	if err != nil {
		log.Print("Stream Track Error")
		h.handleError(c, err)
		return
	}

	filePath, err := h.service.ResolveAudioFilePath(track.AudioURL)
	if err != nil {
		log.Print("Stream Audio Path Error")
		h.handleError(c, err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"message": "audio file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to open audio file"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read audio file"})
		return
	}

	if stat.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"message": "audio file not found"})
		return
	}

	contentType := detectAudioContentType(filePath)
	setStreamingHeaders(c.Writer, contentType)

	go h.service.TrackPlayed(
		context.Background(),
		"",
		track,
		0,
		false,
		"stream",
	)

	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), file)
}

// Enable logging in the default case to see the ACTUAL error
func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidTrackID):
		response.Error(c, appErrors.BadRequest("invalid track id", nil))
	case errors.Is(err, ErrTrackNotFound):
		response.Error(c, appErrors.NotFound("track not found", nil))
	case errors.Is(err, ErrAudioNotFound):
		response.Error(c, appErrors.NotFound("track audio not found", nil))
	case errors.Is(err, ErrInvalidMediaURL):
		response.Error(c, appErrors.BadRequest("invalid media url", nil))
	case errors.Is(err, ErrPrivateTrack):
		response.Error(c, appErrors.Forbidden("track is private", nil))
	default:
		// 🔑 UNCOMMENT & LOG THIS TO FIND THE ROOT CAUSE
		log.Printf("[PLAYER ERROR] %v", err)
		response.Error(c, appErrors.Internal("player error lsknfd", err)) // Pass err for logging
	}
}
