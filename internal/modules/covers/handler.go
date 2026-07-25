package covers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CallbackPayload is the JSON body sent by the Python ML service
// after optimizing a cover image.
type CallbackPayload struct {
	AlbumID  string  `json:"album_id" binding:"required"`
	CoverURL string  `json:"cover_url"`
	ThumbURL string  `json:"thumb_url,omitempty"`
	MedURL   string  `json:"med_url,omitempty"`
	TrackID  *string `json:"track_id,omitempty"`
}

// Handler handles cover optimization callbacks from the ML service.
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new covers handler with a DB connection for
// directly updating album cover URLs.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// HandleCoverCallback receives the result of a cover optimization job
// from the Python ML service and updates the album's cover URL.
// This endpoint is protected by HMAC signature verification, not JWT auth.
func (h *Handler) HandleCoverCallback(c *gin.Context) {
	var payload CallbackPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Validate album_id
	albumID, err := uuid.Parse(payload.AlbumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album_id"})
		return
	}

	// Update album cover URL
	if payload.CoverURL != "" {
		result, err := h.db.ExecContext(
			c.Request.Context(),
			`UPDATE albums SET cover_url = $2, updated_at = NOW() WHERE id = $1`,
			albumID, payload.CoverURL,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update album cover"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "received",
		"album_id":  payload.AlbumID,
		"cover_url": payload.CoverURL,
	})
}
