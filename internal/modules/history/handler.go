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
// @Summary      Get listening history
// @Description  Returns paginated listening history for the current user
// @Tags         history
// @Security     Bearer
// @Param        page  query     int  false  "Page number"
// @Param        limit query     int  false  "Items per page"
// @Success      200   {object}  response.Data{data=[]HistoryResponse}
// @Router       /history [get]
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
