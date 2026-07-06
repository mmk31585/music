package history

import (
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

// GetHistory godoc
// @Summary Get listening history
// @Description Returns paginated listening history for the current user.
// @Tags history
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} HistoryResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /history [get]
func (h *Handler) GetHistory(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	limitRaw := c.Query("limit")
	offsetRaw := c.Query("offset")

	items, pagination, err := h.service.GetHistory(
		c.Request.Context(),
		userID,
		limitRaw,
		offsetRaw,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get listening history",
		})
		return
	}

	c.JSON(http.StatusOK, HistoryResponse{
		Items:      items,
		Pagination: pagination,
	})
}

// DeleteHistoryItem godoc
// @Summary Delete a history entry
// @Description Removes a single listening history entry by its ID.
// @Tags history
// @Produce json
// @Security Bearer
// @Param id path string true "History entry ID"
// @Success 204
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /history/{id} [delete]
func (h *Handler) DeleteHistoryItem(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	historyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid history id"})
		return
	}

	if err := h.service.DeleteItem(c.Request.Context(), userID, historyID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "history entry not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ClearHistory godoc
// @Summary Clear all history
// @Description Removes all listening history for the current user.
// @Tags history
// @Produce json
// @Security Bearer
// @Success 204
// @Failure 401 {object} map[string]interface{}
// @Router /history [delete]
func (h *Handler) ClearHistory(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.ClearAll(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear history"})
		return
	}

	c.Status(http.StatusNoContent)
}

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
