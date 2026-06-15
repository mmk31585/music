package queue

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetQueue godoc
// @Summary Get current user's queue
// @Description Returns the current playback queue for the authenticated user.
// @Tags queue
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} QueueResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue [get]
func (h *Handler) GetQueue(c *gin.Context) {

	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queue"})
		return
	}

	c.JSON(http.StatusOK, QueueResponse{
		Items: items,
		Count: len(items),
	})
}

// AddTrack godoc
// @Summary Add a track to the queue
// @Description Adds a track to the authenticated user's queue.
// @Tags queue
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body AddTrackRequest true "Add track request"
// @Success 201 {object} AddTrackResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue/tracks [post]
func (h *Handler) AddTrack(c *gin.Context) {

	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.AddTrack(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTrackID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		case errors.Is(err, ErrInvalidQueueMode):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue mode. allowed: next, later"})
		case errors.Is(err, ErrTrackNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add track to queue"})
		}
		return
	}

	c.JSON(http.StatusCreated, AddTrackResponse{
		Item: *item,
	})
}

// RemoveTrack godoc
// @Summary Remove a track from the queue
// @Description Removes the queue item specified by id for the authenticated user.
// @Tags queue
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Queue item ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue/tracks/{id} [delete]
func (h *Handler) RemoveTrack(c *gin.Context) {

	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	itemID := c.Param("id")

	err := h.service.Remove(c.Request.Context(), userID, itemID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidQueueItemID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue item id"})
		case errors.Is(err, ErrQueueItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "queue item not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove queue item"})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "queue item removed",
	})
}

// Reorder godoc
// @Summary Reorder queue items
// @Description Reorders items in the authenticated user's queue.
// @Tags queue
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ReorderQueueRequest true "Reorder queue request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue/reorder [put]
func (h *Handler) Reorder(c *gin.Context) {

	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ReorderQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.Reorder(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyReorderList):
			c.JSON(http.StatusBadRequest, gin.H{"error": "reorder list cannot be empty"})
		case errors.Is(err, ErrInvalidQueueItemID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue item id"})
		case errors.Is(err, ErrInvalidPosition):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid position"})
		case errors.Is(err, ErrQueueItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "queue item not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder queue"})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "queue reordered",
	})
}

// Clear godoc
// @Summary Clear queue
// @Description Removes all items from the authenticated user's queue.
// @Tags queue
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} MessageResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue [delete]
func (h *Handler) Clear(c *gin.Context) {

	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.Clear(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear queue"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "queue cleared",
	})
}

// Clear godoc
// @Summary Clear the user's queue
// @Description Removes all items from the authenticated user's queue
// @Tags queue
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=queue.MessageResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security Bearer
// @Router /queue [delete]

func getUserIDFromGin(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("auth_user_id")
	if !exists {
		value, exists = c.Get("userID")
		if !exists {
			value, exists = c.Get("user_id")
			if !exists {
				return uuid.Nil, false
			}
		}
	}

	switch v := value.(type) {
	case uuid.UUID:
		return v, true

	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		return id, true

	default:
		return uuid.Nil, false
	}
}
