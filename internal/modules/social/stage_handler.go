package social

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
