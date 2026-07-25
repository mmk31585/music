package ai

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service     *Service
	rateLimiter map[string]int
	rateMu      sync.Mutex
}

func NewHandler(service *Service) *Handler {
	h := &Handler{
		service:     service,
		rateLimiter: make(map[string]int),
	}
	go h.cleanupRateLimits()
	return h
}

func (h *Handler) cleanupRateLimits() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		h.rateMu.Lock()
		h.rateLimiter = make(map[string]int)
		h.rateMu.Unlock()
	}
}

func (h *Handler) checkRateLimit(userID string) bool {
	h.rateMu.Lock()
	defer h.rateMu.Unlock()
	count := h.rateLimiter[userID]
	if count >= 20 {
		return false
	}
	h.rateLimiter[userID] = count + 1
	return true
}

// GenerateEmbedding godoc
// @Summary Generate track embeddings
// @Description Generates AI embeddings for one or more tracks for similarity search.
// @Tags ai
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body EmbeddingRequest true "Embedding request with track IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/embeddings [post]
func (h *Handler) GenerateEmbedding(c *gin.Context) {
	var req EmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	if len(req.TrackIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track_ids required"})
		return
	}

	results := make([]EmbeddingResponse, 0, len(req.TrackIDs))
	for _, tid := range req.TrackIDs {
		emb, err := h.service.GenerateEmbedding(c.Request.Context(), tid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		results = append(results, EmbeddingResponse{
			TrackID:      emb.TrackID.String(),
			ModelVersion: emb.ModelVersion,
			Dimensions:   len(emb.Embedding),
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// AnalyzeMood godoc
// @Summary Analyze track mood
// @Description Analyzes and generates mood data for a specific track using AI.
// @Tags ai
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body MoodAnalysisRequest true "Mood analysis request with track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/moods [post]
func (h *Handler) AnalyzeMood(c *gin.Context) {
	var req MoodAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	mood, err := h.service.AnalyzeMood(c.Request.Context(), req.TrackID)
	if err != nil {
		if strings.Contains(err.Error(), "mood not found") {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": MoodResponse{
			TrackID:          mood.TrackID.String(),
			MoodTags:         mood.MoodTags,
			Energy:           mood.Energy,
			Valence:          mood.Valence,
			Tempo:            mood.Tempo,
			Danceability:     mood.Danceability,
			Acousticness:     mood.Acousticness,
			Instrumentalness: mood.Instrumentalness,
			Liveness:         mood.Liveness,
			Speechiness:      mood.Speechiness,
		},
	})
}

// GetMood godoc
// @Summary Get track mood
// @Description Returns cached mood data for a specific track.
// @Tags ai
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/moods/{trackId} [get]
func (h *Handler) GetMood(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "track_id required"})
		return
	}

	mood, err := h.service.GetMood(c.Request.Context(), trackID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "mood not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": MoodResponse{
			TrackID:          mood.TrackID.String(),
			MoodTags:         mood.MoodTags,
			Energy:           mood.Energy,
			Valence:          mood.Valence,
			Tempo:            mood.Tempo,
			Danceability:     mood.Danceability,
			Acousticness:     mood.Acousticness,
			Instrumentalness: mood.Instrumentalness,
			Liveness:         mood.Liveness,
			Speechiness:      mood.Speechiness,
		},
	})
}

// GeneratePlaylist godoc
// @Summary Generate AI playlist
// @Description Generates a playlist based on a text prompt, mood, or seed track using AI.
// @Tags ai
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body GeneratePlaylistRequest true "Playlist generation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/playlists/generate [post]
func (h *Handler) GeneratePlaylist(c *gin.Context) {
	var req GeneratePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if len(req.Prompt) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "prompt too long (max 500 characters)"})
		return
	}

	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "user not authenticated"})
		return
	}

	if !h.checkRateLimit(userID.String()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"success": false, "message": "rate limit exceeded (max 20 generations per minute)"})
		return
	}

	playlist, err := h.service.GeneratePlaylist(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(playlist.Tracks) > req.Limit {
		playlist.Tracks = playlist.Tracks[:req.Limit]
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": playlist})
}

// SimilarByMood godoc
// @Summary Find similar tracks by mood
// @Description Returns tracks with similar mood characteristics to a given track.
// @Tags ai
// @Produce json
// @Security Bearer
// @Param trackId path string true "Source track ID"
// @Param mood query string false "Target mood filter"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/similar/mood/{trackId} [get]
func (h *Handler) SimilarByMood(c *gin.Context) {
	trackID := c.Param("trackId")
	mood := c.Query("mood")
	limit := 20

	results, err := h.service.GetSimilarByMood(c.Request.Context(), trackID, mood, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// SimilarByEmbedding godoc
// @Summary Find similar tracks by embedding
// @Description Returns tracks with similar audio/text embeddings to a given track.
// @Tags ai
// @Produce json
// @Security Bearer
// @Param trackId path string true "Source track ID"
// @Param space query string false "Embedding space (audio, text)" default(audio)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/similar/embedding/{trackId} [get]
func (h *Handler) SimilarByEmbedding(c *gin.Context) {
	trackID := c.Param("trackId")
	limit := 20
	embeddingSpace := c.DefaultQuery("space", "audio")

	results, err := h.service.GetSimilarByEmbedding(c.Request.Context(), trackID, limit, embeddingSpace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
}

var ErrNotFound = errors.New("not found")
