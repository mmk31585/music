package social

import (
	"net/http"
	"strconv"

	appErr "music/internal/common/errors"
	"music/internal/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) CreateParty(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	party, err := h.service.CreateParty(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create party"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": party})
}

func (h *Handler) GetParty(c *gin.Context) {
	id := c.Param("id")
	party, err := h.service.GetParty(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "party not found"})
		return
	}
	participants, _ := h.service.GetPartyParticipants(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": party, "participants": participants})
}

func (h *Handler) ListActiveParties(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	parties, err := h.service.ListActiveParties(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list parties"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": parties})
}

func (h *Handler) UpdatePartyStatus(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req struct {
		Status  string  `json:"status" binding:"required"`
		TrackID *string `json:"track_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdatePartyStatus(c.Request.Context(), id, req.Status, userID, req.TrackID); err != nil {
		h.logger.Error("update party status failed", zap.String("party_id", id), zap.String("status", req.Status), zap.Error(err))
		response.Error(c, appErr.Internal("failed to update party", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) JoinParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinParty(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("join party failed", zap.String("party_id", id), zap.String("user_id", userID), zap.Error(err))
		response.Error(c, appErr.Internal("failed to join party", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveParty(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// --- Live Rooms ---

func (h *Handler) CreateRoom(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := h.service.CreateRoom(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": room})
}

func (h *Handler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": room})
}

func (h *Handler) ListActiveRooms(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	rooms, err := h.service.ListActiveRooms(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rooms})
}

func (h *Handler) JoinRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinRoom(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveRoom(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetRoomParticipants(c *gin.Context) {
	id := c.Param("id")
	participants, err := h.service.GetRoomParticipants(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get participants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": participants})
}

func (h *Handler) AddToRoomQueue(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req AddToQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddToRoomQueue(c.Request.Context(), id, req.TrackID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to queue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetRoomQueue(c *gin.Context) {
	id := c.Param("id")
	queue, err := h.service.GetRoomQueue(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": queue})
}

// --- Democratic Voting Queue ---

func (h *Handler) SuggestTrack(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req struct {
		TrackID string `json:"track_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.queueEngine.SuggestTrack(c.Request.Context(), id, req.TrackID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CastVote(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	candidateID := c.Param("candidateId")
	userID := c.GetString("auth_user_id")
	if err := h.queueEngine.CastVote(c.Request.Context(), candidateID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": id})
}

func (h *Handler) RemoveVote(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	candidateID := c.Param("candidateId")
	userID := c.GetString("auth_user_id")
	if err := h.queueEngine.RemoveVote(c.Request.Context(), candidateID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetQueueState(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	state, err := h.queueEngine.GetQueueState(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": state})
}

func (h *Handler) TrackEnded(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	var req struct {
		TrackID string `json:"track_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.queueEngine.TrackEnded(c.Request.Context(), id, req.TrackID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// --- Music Clubs ---
