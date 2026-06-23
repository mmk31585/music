package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type OnboardingHandler struct {
	tasteService *TasteProfileService
}

func NewOnboardingHandler(tasteService *TasteProfileService) *OnboardingHandler {
	return &OnboardingHandler{tasteService: tasteService}
}

type onboardingGenresRequest struct {
	GenreIDs []string `json:"genre_ids" binding:"required,min=1"`
}

// SetOnboardingGenres saves the user's cold-start genre preferences.
// Called once after registration (though not strictly enforced — overwrites).
func (h *OnboardingHandler) SetOnboardingGenres(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req onboardingGenresRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: at least one genre_id required"})
		return
	}

	if err := h.tasteService.SetOnboardingGenres(c.Request.Context(), userID, req.GenreIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save genre preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
