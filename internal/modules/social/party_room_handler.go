package social

import (
	"net/http"
	"strconv"

	appErr "music/internal/common/errors"
	"music/internal/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateParty godoc
// @Summary Create a listening party
// @Description Creates a new listening party for the authenticated user.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreatePartyRequest true "Create party request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/parties [post]
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

// GetParty godoc
// @Summary Get party by ID
// @Description Returns a listening party by its ID.
// @Tags social
// @Produce json
// @Param id path string true "Party ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/parties/{id} [get]
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

// ListActiveParties godoc
// @Summary List active parties
// @Description Returns a paginated list of active listening parties.
// @Tags social
// @Produce json
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/parties [get]
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

// UpdatePartyStatus godoc
// @Summary Update party status
// @Description Updates the status of a listening party (host only).
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Party ID"
// @Param request body object{status=string,track_id=string} true "Party status update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/parties/{id}/status [put]
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
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to update party", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// JoinParty godoc
// @Summary Join a party
// @Description Allows the authenticated user to join a listening party.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Party ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/parties/{id}/join [post]
func (h *Handler) JoinParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinParty(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("join party failed", zap.String("party_id", id), zap.String("user_id", userID), zap.Error(err))
		response.Error(c, appErr.New(http.StatusInternalServerError, appErr.CodeInternal, "failed to join party", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LeaveParty godoc
// @Summary Leave a party
// @Description Allows the authenticated user to leave a listening party.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Party ID"
// @Success 200 {object} map[string]interface{}
// @Router /social/parties/{id}/leave [post]
func (h *Handler) LeaveParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveParty(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// --- Live Rooms ---

// CreateRoom godoc
// @Summary Create a live room
// @Description Creates a new live listening room for the authenticated user.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateRoomRequest true "Create room request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms [post]
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

// GetRoom godoc
// @Summary Get room by ID
// @Description Returns a live room by its ID.
// @Tags social
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms/{id} [get]
func (h *Handler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": room})
}

// ListActiveRooms godoc
// @Summary List active rooms
// @Description Returns a paginated list of active live rooms.
// @Tags social
// @Produce json
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms [get]
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

// JoinRoom godoc
// @Summary Join a room
// @Description Allows the authenticated user to join a live room.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms/{id}/join [post]
func (h *Handler) JoinRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinRoom(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LeaveRoom godoc
// @Summary Leave a room
// @Description Allows the authenticated user to leave a live room.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Router /social/rooms/{id}/leave [post]
func (h *Handler) LeaveRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveRoom(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetRoomParticipants godoc
// @Summary Get room participants
// @Description Returns the participants of a live room.
// @Tags social
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms/{id}/participants [get]
func (h *Handler) GetRoomParticipants(c *gin.Context) {
	id := c.Param("id")
	participants, err := h.service.GetRoomParticipants(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get participants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": participants})
}

// AddToRoomQueue godoc
// @Summary Add track to room queue
// @Description Adds a track to a live room's playback queue.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param request body AddToQueueRequest true "Add to queue request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue [post]
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

// GetRoomQueue godoc
// @Summary Get room queue
// @Description Returns the current playback queue of a live room.
// @Tags social
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue [get]
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

// SuggestTrack godoc
// @Summary Suggest a track for room queue
// @Description Suggests a track for the democratic voting queue in a room.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param request body object{track_id=string} true "Track suggestion"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue/suggest [post]
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

// CastVote godoc
// @Summary Cast a vote on a queued track
// @Description Casts a vote for a track candidate in the democratic queue.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param candidateId path string true "Candidate track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue/{candidateId}/vote [post]
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

// RemoveVote godoc
// @Summary Remove vote from a queued track
// @Description Removes a vote from a track candidate in the democratic queue.
// @Tags social
// @Produce json
// @Security Bearer
// @Param candidateId path string true "Candidate track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue/{candidateId}/vote [delete]
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

// GetQueueState godoc
// @Summary Get room queue state
// @Description Returns the current democratic queue state including suggestions and votes.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/queue/state [get]
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

// TrackEnded godoc
// @Summary Notify track ended
// @Description Notifies the queue engine that a track has finished playing.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Param request body object{track_id=string} true "Track ended notification"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 501 {object} map[string]interface{}
// @Router /social/rooms/{id}/track-ended [post]
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
