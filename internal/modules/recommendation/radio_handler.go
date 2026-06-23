package recommendation

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type StartRadioRequest struct {
	SeedTrackID string `json:"seed_track_id" binding:"required"`
	SeedType    string `json:"seed_type"`
}

func (h *Handler) StartRadio(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req StartRadioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "seed_track_id is required"})
		return
	}

	if req.SeedType == "" {
		req.SeedType = "track"
	}

	session, tracks, err := h.radioService.StartRadio(c.Request.Context(), userID, req.SeedTrackID, req.SeedType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": StartRadioResponse{
			SessionID: session.ID,
			Tracks:    tracks,
		},
	})
}

func (h *Handler) GetNextRadioBatch(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "session id is required"})
		return
	}

	// Verify ownership
	session, err := h.radioService.repo.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal error"})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "session not found"})
		return
	}
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "not your session"})
		return
	}

	count := 10
	if c.Query("count") != "" {
		if _, err := fmt.Sscanf(c.Query("count"), "%d", &count); err != nil || count < 1 {
			count = 10
		}
		if count > 50 {
			count = 50
		}
	}

	tracks, err := h.radioService.GetNextBatch(c.Request.Context(), sessionID, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": NextBatchResponse{
			Tracks:  tracks,
			HasMore: len(tracks) >= count,
		},
	})
}

func (h *Handler) EndRadio(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "session id is required"})
		return
	}

	session, err := h.radioService.repo.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal error"})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "session not found"})
		return
	}
	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "not your session"})
		return
	}

	if err := h.radioService.EndRadio(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "radio ended"})
}
