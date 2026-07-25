package social

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RaiseHand godoc
// @Summary Raise hand in room
// @Description Raises the user's hand to request speaking in a room stage.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/raise-hand [post]
func (h *Handler) RaiseHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.RaiseHand(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LowerHand godoc
// @Summary Lower hand in room
// @Description Lowers the user's hand in a room stage.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/lower-hand [post]
func (h *Handler) LowerHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.LowerHand(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ApproveHand godoc
// @Summary Approve a raised hand
// @Description Approves a user's request to speak on stage (host only).
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param userId path string true "User ID to approve"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/{userId}/approve [post]
func (h *Handler) ApproveHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.ApproveHand(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DenyHand godoc
// @Summary Deny a raised hand
// @Description Denies a user's request to speak on stage (host only).
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param userId path string true "User ID to deny"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/{userId}/deny [post]
func (h *Handler) DenyHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.DenyHand(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RemoveFromStage godoc
// @Summary Remove user from stage
// @Description Removes a user from the stage (host only).
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param userId path string true "User ID to remove"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/{userId} [delete]
func (h *Handler) RemoveFromStage(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.RemoveFromStage(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LeaveStage godoc
// @Summary Leave the stage
// @Description Allows the authenticated user to leave the stage.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/leave [post]
func (h *Handler) LeaveStage(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.LeaveStage(c.Request.Context(), roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ToggleMute godoc
// @Summary Toggle user mute on stage
// @Description Mutes or unmutes a user on stage (host only).
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param userId path string true "Target user ID"
// @Param request body object{muted=bool} true "Mute toggle"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage/{userId}/mute [post]
func (h *Handler) ToggleMute(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	actorID := c.GetString("auth_user_id")
	var req struct {
		Muted bool `json:"muted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.stageManager.ToggleMute(c.Request.Context(), roomID, actorID, targetUserID, req.Muted); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetStageState godoc
// @Summary Get stage state
// @Description Returns the current stage state for a room including who is on stage.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/stage [get]
func (h *Handler) GetStageState(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	userID := c.GetString("auth_user_id")
	state, err := h.stageManager.GetStageState(c.Request.Context(), roomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": state})
}
