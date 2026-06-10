package ai

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) GeneratePlaylist(c *gin.Context) {
	var req GeneratePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	userID, _ := web.GetRequiredUserUUID(c)
	ctx := c.Request.Context()

	playlist, err := h.service.GeneratePlaylist(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(playlist.Tracks) > req.Limit {
		playlist.Tracks = playlist.Tracks[:req.Limit]
	}

	_ = userID
	c.JSON(http.StatusOK, gin.H{"success": true, "data": playlist})
}

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

func (h *Handler) SimilarByEmbedding(c *gin.Context) {
	trackID := c.Param("trackId")
	limit := 20

	results, err := h.service.GetSimilarByEmbedding(c.Request.Context(), trackID, limit)
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
