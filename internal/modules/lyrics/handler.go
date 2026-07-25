package lyrics

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service          *Service
	mlClient         *MLServiceClient
	openRouterClient *OpenRouterClient
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// SetMLClient sets the ML service client used for AI lyrics generation.
func (h *Handler) SetMLClient(client *MLServiceClient) {
	h.mlClient = client
}

// SetOpenRouterClient sets the OpenRouter client used for AI lyrics fetch/sync/review.
func (h *Handler) SetOpenRouterClient(client *OpenRouterClient) {
	h.openRouterClient = client
}

// LyricsCallbackPayload is the JSON body sent by the Python ML service
// after processing a lyrics generation job.
type LyricsCallbackPayload struct {
	TrackID             string  `json:"track_id" binding:"required"`
	LRCContent          string  `json:"lrc_content"`
	PlainText           string  `json:"plain_text"`
	Confidence          float64 `json:"confidence"`
	DetectedLanguage    string  `json:"detected_language"`
	WhisperModelVersion string  `json:"whisper_model_version"`
}

// HandleLyricsCallback godoc
// @Summary Handle AI lyrics callback
// @Description Receives the result of an AI lyrics generation job from the Python ML service and persists lyrics. Protected by HMAC signature verification (not JWT).
// @Tags lyrics
// @Accept json
// @Produce json
// @Param request body LyricsCallbackPayload true "Lyrics callback payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /internal/v1/lyrics/callback [post]
func (h *Handler) HandleLyricsCallback(c *gin.Context) {
	var payload LyricsCallbackPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Determine the lyrics type from the content
	lyricsType := "lrc"
	if !isLRCLyrics(payload.LRCContent) {
		lyricsType = "plain"
	}

	// Use the detected language, falling back to Persian (fa) as default
	// since most tracks on Muse are Persian.
	lang := payload.DetectedLanguage
	if lang == "" {
		lang = "fa"
	}

	// Build the LRC content — prefer synced, fall back to plain text
	content := payload.LRCContent
	if content == "" && payload.PlainText != "" {
		content = payload.PlainText
		lyricsType = "plain"
	}
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no lyrics content provided"})
		return
	}

	confidence := payload.Confidence

	req := CreateLyricsRequest{
		TrackID:         payload.TrackID,
		Language:        lang,
		Type:            lyricsType,
		Content:         content,
		Source:          "ai_generated",
		ConfidenceScore: &confidence,
	}

	lyrics, err := h.service.CreateLyrics(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTrackNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		case errors.Is(err, ErrLyricsAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "lyrics already exist for this track and language"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save lyrics"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "received",
		"lyrics_id": lyrics.ID,
	})
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

// FetchOrGenerateLyrics godoc
// @Summary Fetch lyrics from LRCLIB or generate via AI
// @Description Tries LRCLIB first. If no lyrics found, enqueues an AI lyrics generation job.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/fetch-or-generate/{trackId} [post]
func (h *Handler) FetchOrGenerateLyrics(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}
	if _, err := uuid.Parse(trackID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid track id"})
		return
	}

	// Step 1: Try LRCLIB first
	lyrics, err := h.service.FetchFromLRC(c.Request.Context(), trackID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"source":  "lrclib",
			"data":    ToLyricsResponse(lyrics),
		})
		return
	}

	// LRCLIB failed — check if it's a recoverable case
	if !errors.Is(err, ErrLRCLibNoLyrics) &&
		!errors.Is(err, ErrLRCLibEmptyResponse) &&
		!errors.Is(err, ErrTrackNoTitle) &&
		!errors.Is(err, ErrTrackNotFound) {
		// Genuine error (track not found, DB error, etc.)
		h.handleError(c, err)
		return
	}

	// Step 2: Try OpenRouter (AI-generated lyrics without audio)
	if h.openRouterClient != nil {
		info, infoErr := h.service.GetTrackInfo(c.Request.Context(), trackID)
		if infoErr == nil {
			lrcContent, lang, aiErr := h.openRouterClient.FetchLyrics(c.Request.Context(), info.Title, info.ArtistName)
			if aiErr == nil && lrcContent != "" {
				// Save OpenRouter-generated lyrics
				req := CreateLyricsRequest{
					TrackID:  trackID,
					Language: lang,
					Type:     "lrc",
					Content:  lrcContent,
					Source:   "ai_generated",
				}
				if req.Language == "" {
					req.Language = "fa"
				}
				lyrics, saveErr := h.service.CreateLyrics(c.Request.Context(), req)
				if saveErr == nil {
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"source":  "openrouter",
						"data":    ToLyricsResponse(lyrics),
					})
					return
				}
				// If save failed (e.g. already exists), still return the content
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"source":  "openrouter",
					"data": LyricsResponse{
						ID:        "",
						TrackID:   trackID,
						Language:  lang,
						Type:      "lrc",
						Content:   lrcContent,
						Source:    "ai_generated",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
				})
				return
			}
		}
	}

	// Step 3: Fall back to Whisper (requires audio file)
	if h.mlClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"source":  "none",
			"message": "No lyrics found on LRCLIB or OpenRouter, and Whisper AI generation is not configured",
		})
		return
	}

	// Get track info for Whisper job
	info, infoErr := h.service.GetTrackInfo(c.Request.Context(), trackID)
	if infoErr != nil {
		h.handleError(c, infoErr)
		return
	}

	// Check if track has an audio file to transcribe
	if info.AudioURL == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"source":  "none",
			"message": "No lyrics found on LRCLIB or OpenRouter, and track has no audio file for Whisper transcription",
		})
		return
	}

	// Build audio source for ML service
	var audioFilePath, audioDownloadURL string
	if strings.HasPrefix(info.AudioURL, "http://") || strings.HasPrefix(info.AudioURL, "https://") {
		audioDownloadURL = info.AudioURL
	} else {
		audioFilePath = strings.TrimLeft(info.AudioURL, "/")
	}

	jobReq := CreateLyricsJobRequest{
		TrackID:          trackID,
		AudioFilePath:    audioFilePath,
		AudioDownloadURL: audioDownloadURL,
		TrackTitle:       info.Title,
		TrackArtist:      info.ArtistName,
	}

	jobID, err := h.mlClient.EnqueueLyricsJob(c.Request.Context(), jobReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"source":  "none",
			"message": "No lyrics found on LRCLIB or OpenRouter. Whisper AI lyrics generation is unavailable: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":  true,
		"source":   "ai",
		"job_id":   jobID,
		"track_id": trackID,
		"message":  "No lyrics found on LRCLIB or OpenRouter. Whisper AI lyrics generation has been started and will be saved automatically when ready.",
	})
}

// AILyricsStatus godoc
// @Summary Check AI lyrics generation status for a track
// @Description Proxies to the ML service to get the latest AI lyrics job status for a track.
// @Tags lyrics
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/ai-status/{trackId} [get]
func (h *Handler) AILyricsStatus(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track id is required"})
		return
	}

	// First check if lyrics already exist in the DB (AI might have finished)
	lyricsResult, lyricsErr := h.service.GetTrackLyrics(c.Request.Context(), trackID, "")
	if lyricsErr == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  "completed",
			"source":  "ai",
			"lyrics":  lyricsResult,
		})
		return
	}

	// No lyrics yet — query ML service for job status
	if h.mlClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  "unavailable",
			"message": "AI lyrics generation is not configured",
		})
		return
	}

	// Proxy to ML service's by-track endpoint
	jobStatus, err := h.mlClient.GetJobByTrack(c.Request.Context(), trackID)
	if err != nil {
		// No job found or ML service unreachable
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  "not_found",
			"message": "No AI lyrics job found for this track",
		})
		return
	}

	// Map ML job status to human-readable status
	humanStatus := jobStatus.Status
	message := ""

	switch jobStatus.Status {
	case "queued":
		message = "AI lyrics job is queued — waiting for a worker..."
	case "downloading":
		message = "Downloading audio file for transcription..."
	case "transcribing":
		message = "Transcribing audio with AI (Whisper)..."
	case "formatting":
		message = "Formatting transcribed lyrics..."
	case "callback_pending":
		message = "Saving generated lyrics..."
	case "completed":
		message = "Lyrics generated successfully!"
		humanStatus = "completed"
	case "failed":
		message = "AI lyrics generation failed: " + jobStatus.ErrorMessage
	default:
		message = "AI lyrics generation is in progress..."
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"status":       humanStatus,
		"source":       "ai",
		"job_id":       jobStatus.JobID,
		"track_id":     jobStatus.TrackID,
		"message":      message,
		"progress_pct": jobProgressPercent(jobStatus.Status),
	})
}

// jobProgressPercent maps ML job status to a percentage for UI progress bars.
func jobProgressPercent(status string) int {
	switch status {
	case "queued":
		return 5
	case "downloading":
		return 15
	case "transcribing":
		return 50
	case "formatting":
		return 80
	case "callback_pending":
		return 90
	case "completed":
		return 100
	case "failed":
		return -1
	default:
		return 0
	}
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
		slog.Error("CreateLyrics: invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body: " + err.Error()})
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

// SyncWithAI godoc
// @Summary Sync plain lyrics to LRC via OpenRouter
// @Description Takes plain text lyrics and adds LRC timestamps using OpenRouter AI.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body SyncLyricsRequest true "Sync lyrics request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/sync [post]
func (h *Handler) SyncWithAI(c *gin.Context) {
	if h.openRouterClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"source":  "none",
			"message": "OpenRouter AI is not configured (set OPENROUTER_API_KEY)",
		})
		return
	}

	var req SyncLyricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	// Validate track exists and get duration
	info, infoErr := h.service.GetTrackInfo(c.Request.Context(), req.TrackID)
	if infoErr != nil {
		h.handleError(c, infoErr)
		return
	}

	lrcContent, lang, err := h.openRouterClient.SyncLyrics(
		c.Request.Context(),
		info.Title, info.ArtistName,
		req.PlainText, info.DurationSeconds,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": friendlyOpenRouterError(err),
		})
		return
	}

	if lrcContent == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": "AI returned empty results. The model may be overloaded or the input may need more context. Try again later.",
		})
		return
	}

	// Upsert: update existing lyrics or create new ones
	savedLyrics, saveErr := h.upsertLyrics(c.Request.Context(), req.TrackID, lang, lrcContent)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": "Failed to save synced lyrics: " + saveErr.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"source":  "openrouter",
		"data":    ToLyricsResponse(savedLyrics),
	})
}

// ReviewWithAI godoc
// @Summary Review and fix lyrics via OpenRouter
// @Description Reviews existing LRC lyrics and fixes spelling/timing issues using OpenRouter AI.
// @Tags lyrics
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ReviewLyricsRequest true "Review lyrics request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/lyrics/review [post]
func (h *Handler) ReviewWithAI(c *gin.Context) {
	if h.openRouterClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"source":  "none",
			"message": "OpenRouter AI is not configured (set OPENROUTER_API_KEY)",
		})
		return
	}

	var req ReviewLyricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	// Get track info for context
	info, _ := h.service.GetTrackInfo(c.Request.Context(), req.TrackID)
	title := req.TrackTitle
	artist := req.ArtistName
	if info != nil {
		if title == "" {
			title = info.Title
		}
		if artist == "" {
			artist = info.ArtistName
		}
	}

	fixedLRC, lang, err := h.openRouterClient.ReviewLyrics(
		c.Request.Context(),
		title, artist,
		req.ExistingLRC,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": friendlyOpenRouterError(err),
		})
		return
	}

	if fixedLRC == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": "AI returned empty results. The model may be overloaded or the input may need more context. Try again later.",
		})
		return
	}

	// Persist the reviewed lyrics (update in place)
	lyricsLang := lang
	if lyricsLang == "" {
		lyricsLang = "fa"
	}
	savedLyrics, saveErr := h.upsertLyrics(c.Request.Context(), req.TrackID, lyricsLang, fixedLRC)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"source":  "openrouter",
			"message": "Failed to save reviewed lyrics: " + saveErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"source":   "openrouter",
		"language": lyricsLang,
		"data":     ToLyricsResponse(savedLyrics),
	})
}

// upsertLyrics looks for existing lyrics for a track+language pair and
// updates the content, or creates a new record if none exists.
func (h *Handler) upsertLyrics(ctx context.Context, trackID, language, lrcContent string) (Lyrics, error) {
	if language == "" {
		language = "fa"
	}

	// Try to find existing lyrics
	existingResp, lookupErr := h.service.GetTrackLyrics(ctx, trackID, language)
	if lookupErr == nil && existingResp.ID != "" {
		lrcType := "lrc"
		return h.service.UpdateLyrics(ctx, existingResp.ID, UpdateLyricsRequest{
			Type:    &lrcType,
			Content: &lrcContent,
		})
	}

	// No existing lyrics — create new
	return h.service.CreateLyrics(ctx, CreateLyricsRequest{
		TrackID:  trackID,
		Language: language,
		Type:     "lrc",
		Content:  lrcContent,
		Source:   "ai_generated",
	})
}

// friendlyOpenRouterError maps common Go-level OpenRouter errors to user-facing messages.
func friendlyOpenRouterError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "TLS handshake timeout"):
		return "OpenRouter is not responding. The AI service may be overloaded. Try again later."
	case strings.Contains(msg, "i/o timeout"):
		return "OpenRouter request timed out. The free model may be slow right now. Try again later."
	case strings.Contains(msg, "no such host"):
		return "Cannot reach OpenRouter. Check your internet connection."
	case strings.Contains(msg, "connection refused"):
		return "OpenRouter refused the connection. The service may be down."
	case strings.Contains(msg, "status 429"):
		return "OpenRouter rate limit exceeded. Wait a moment and try again."
	case strings.Contains(msg, "status 401"):
		return "OpenRouter API key is invalid. Check your OPENROUTER_API_KEY."
	case strings.Contains(msg, "status 402"):
		return "OpenRouter free tier quota exceeded. Check your OpenRouter account."
	case strings.Contains(msg, "status 404"):
		return "OpenRouter model not found. The AI model may be unavailable or deprecated. Check OPENROUTER_MODEL."
	case strings.Contains(msg, "status 5"):
		return "OpenRouter server error. The AI service may be overloaded. Try again later."
	default:
		// Strip Go-internal wrapping for a cleaner message but keep the substance.
		if strings.HasPrefix(msg, "openrouter ") {
			msg = strings.TrimPrefix(msg, "openrouter ")
		}
		return "OpenRouter AI error: " + msg
	}
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLyricsNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrLyricsAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackNoTitle):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidLyricsType):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidLyricsLanguage):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrForbiddenLyricsAccess):
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrLRCLibNoLyrics):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrLRCLibEmptyResponse):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
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
