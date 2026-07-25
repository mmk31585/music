package reputation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"music/internal/common/response"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetUserReputation(c *gin.Context) error {
	userID := c.Param("userId")
	summary, err := h.svc.GetUserReputation(c.Request.Context(), userID)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "user reputation retrieved", summary)
	return nil
}

func (h *Handler) GetTrustTiers(c *gin.Context) error {
	tiers, err := h.svc.GetTrustTiers(c.Request.Context())
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "trust tiers retrieved", tiers)
	return nil
}

func (h *Handler) GetTopContributors(c *gin.Context) error {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	contributors, err := h.svc.GetTopContributors(c.Request.Context(), limit)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "top contributors retrieved", contributors)
	return nil
}
