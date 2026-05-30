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

func getUserIDFromGin(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("userID")
	if !exists {
		value, exists = c.Get("user_id")
		if !exists {
			return uuid.Nil, false
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
