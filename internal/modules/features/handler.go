package features

import (
	"net/http"

	"music/internal/config"

	"github.com/gin-gonic/gin"

	"music/internal/common/response"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/features", h.List)
}

type Handler struct {
	cfg config.FeaturesConfig
}

func NewHandler(cfg config.FeaturesConfig) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, "ok", gin.H{
		"analytics":        h.cfg.Analytics,
		"recommendation":   h.cfg.Recommendation,
		"search":           h.cfg.Search,
		"social":           h.cfg.Social,
		"reactions":        h.cfg.Reactions,
		"creator":          h.cfg.Creator,
		"moderation":       h.cfg.Moderation,
		"ai":               h.cfg.AI,
		"contribution":     h.cfg.Contribution,
		"gamification":     h.cfg.Gamification,
		"tips":             h.cfg.Tips,
		"subscription":     h.cfg.Subscription,
		"notification":     h.cfg.Notification,
		"redesignedPlayer": false,
	})
}
