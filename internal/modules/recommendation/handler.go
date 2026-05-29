package recommendation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func getUserIDFromGin(c *gin.Context) (string, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return "", false
	}

	userID, ok := v.(string)
	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}

func parseLimit(c *gin.Context) (int, error) {
	raw := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(raw)
	if err != nil {
		return 0, ErrInvalidLimit
	}
	return limit, nil
}

func (h *Handler) PopularTracks(c *gin.Context) {
	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.PopularTracks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch popular tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "popular",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) BestTracks(c *gin.Context) {
	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.BestTracks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch best tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "best",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) RecentTracks(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.RecentTracks(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch recent tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "recent",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) SimilarTracks(c *gin.Context) {
	trackID := c.Param("trackId")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "trackId is required"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.SimilarTracks(c.Request.Context(), trackID, limit)
	if err != nil {
		if errors.Is(err, ErrTrackNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "track not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch similar tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "similar",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) TracksByArtist(c *gin.Context) {
	artistID := c.Param("artistId")
	if artistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "artistId is required"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.TracksByArtist(c.Request.Context(), artistID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch artist tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "artist",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) TracksByGenre(c *gin.Context) {
	genre := c.Param("genre")
	if genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "genre is required"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.TracksByGenre(c.Request.Context(), genre, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch genre tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "genre",
			Items: items,
			Limit: limit,
		},
	})
}

func (h *Handler) ForYou(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	items, err := h.service.ForYou(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RecommendationResponse{
			Type:  "for_you",
			Items: items,
			Limit: limit,
		},
	})
}
