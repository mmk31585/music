package social

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRating godoc
// @Summary Create a rating
// @Description Creates or updates a track rating for the authenticated user.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateRatingRequest true "Create rating request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/ratings [post]
func (h *Handler) CreateRating(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rating, err := h.service.CreateRating(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rating"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": rating})
}

// GetTrackRatings godoc
// @Summary Get track ratings
// @Description Returns ratings and average score for a track.
// @Tags social
// @Produce json
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/ratings/{trackId} [get]
func (h *Handler) GetTrackRatings(c *gin.Context) {
	trackID := c.Param("trackId")
	ratings, err := h.service.GetTrackRatings(c.Request.Context(), trackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings"})
		return
	}
	avg, count, _ := h.service.GetTrackRatingAverage(c.Request.Context(), trackID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ratings,
		"average": avg,
		"count":   count,
	})
}
