package recommendation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func parseLimit(c *gin.Context) (int, error) {
	raw := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(raw)
	if err != nil {
		return 0, ErrInvalidLimit
	}
	return limit, nil
}

// PopularTracks godoc
// @Summary Get popular tracks
// @Description Returns a list of popular tracks.
// @Tags recommendation
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/popular [get]
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

// BestTracks godoc
// @Summary Get best tracks
// @Description Returns a list of best tracks.
// @Tags recommendation
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/best [get]
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

// RecentTracks godoc
// @Summary Get recent tracks for current user
// @Description Returns recently played tracks for the authenticated user.
// @Tags recommendation
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/recent [get]
func (h *Handler) RecentTracks(c *gin.Context) {

	userID, ok := web.GetUserIDString(c)
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

// SimilarTracks godoc
// @Summary Get similar tracks
// @Description Returns tracks similar to the provided track ID.
// @Tags recommendation
// @Accept json
// @Produce json
// @Param trackId path string true "Track ID"
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/{trackId}/similar [get]
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

// TracksByArtist godoc
// @Summary Get tracks by artist
// @Description Returns recommended tracks for the given artist.
// @Tags recommendation
// @Accept json
// @Produce json
// @Param artistId path string true "Artist ID"
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/artist/{artistId} [get]
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

// TracksByGenre godoc
// @Summary Get tracks by genre
// @Description Returns recommended tracks for the given genre.
// @Tags recommendation
// @Accept json
// @Produce json
// @Param genre path string true "Genre"
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/genre/{genre} [get]
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

// ForYou godoc
// @Summary Get personalized recommendations
// @Description Returns personalized recommendations for the authenticated user.
// @Tags recommendation
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recommendations/for-you [get]
func (h *Handler) ForYou(c *gin.Context) {

	userID, ok := web.GetUserIDString(c)
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

// ForYou godoc
// @Summary Get personalized recommendations for current user
// @Description Returns personalized recommendations for the authenticated user
// @Tags recommendation
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items"
// @Success 200 {object} response.SuccessResponse{data=recommendation.RecommendationResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security Bearer
// @Router /recommendations/for_you [get]
