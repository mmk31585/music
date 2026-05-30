package notification

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListNotifications godoc
// @Summary List notifications
// @Description Returns paginated notifications for the authenticated user.
// @Tags notifications
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} ListNotificationsResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications [get]
func (h *Handler) ListNotifications(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	res, err := h.service.ListNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Marks a specific notification as read for the authenticated user.
// @Tags notifications
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} MarkReadResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications/{id}/read [patch]
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification id is required"})
		return
	}

	err := h.service.MarkAsRead(c.Request.Context(), userID, notificationID)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, MarkReadResponse{
		Message: "Notification marked as read",
	})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Marks all notifications as read for the authenticated user.
// @Tags notifications
// @Produce json
// @Security Bearer
// @Success 200 {object} MarkReadResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications/read-all [patch]
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark all notifications as read"})
		return
	}

	c.JSON(http.StatusOK, MarkReadResponse{
		Message: "All notifications marked as read",
	})
}
