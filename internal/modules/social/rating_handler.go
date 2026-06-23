package social

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
