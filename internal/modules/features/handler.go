package features

import (
	"net/http"

	"music/internal/config"
	"music/internal/common/response"

	"github.com/gin-gonic/gin"
)

// FeatureFlagsResponse represents the feature flags payload.
type FeatureFlagsResponse struct {
	Analytics      bool `json:"analytics" example:"true"`
	Recommendation bool `json:"recommendation" example:"true"`
	Search         bool `json:"search" example:"true"`
	Social         bool `json:"social" example:"true"`
	Reactions      bool `json:"reactions" example:"true"`
	Creator        bool `json:"creator" example:"true"`
	Moderation     bool `json:"moderation" example:"true"`
	AI             bool `json:"ai" example:"true"`
	Contribution   bool `json:"contribution" example:"true"`
	Gamification   bool `json:"gamification" example:"true"`
	Tips           bool `json:"tips" example:"true"`
	Subscription   bool `json:"subscription" example:"true"`
	Notification   bool `json:"notification" example:"true"`
}

type Handler struct {
	cfg *config.FeaturesConfig
}

func NewHandler(cfg *config.FeaturesConfig) *Handler {
	return &Handler{cfg: cfg}
}

// GetFeatures godoc
// @Summary Get feature flags
// @Description Returns the current state of all module-level feature flags.
// @Tags features
// @Produce json
// @Success 200 {object} response.SuccessResponseData[FeatureFlagsResponse] "feature flags"
// @Router /features [get]
func (h *Handler) GetFeatures(c *gin.Context) {
	response.Success(c, http.StatusOK, "feature flags", gin.H{
		"analytics":      h.cfg.Analytics,
		"recommendation": h.cfg.Recommendation,
		"search":         h.cfg.Search,
		"social":         h.cfg.Social,
		"reactions":      h.cfg.Reactions,
		"creator":        h.cfg.Creator,
		"moderation":     h.cfg.Moderation,
		"ai":             h.cfg.AI,
		"contribution":   h.cfg.Contribution,
		"gamification":   h.cfg.Gamification,
		"tips":           h.cfg.Tips,
		"subscription":   h.cfg.Subscription,
		"notification":   h.cfg.Notification,
	})
}
