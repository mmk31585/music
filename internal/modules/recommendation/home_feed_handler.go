package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

func (h *Handler) HomeFeed(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	feed, err := h.homeFeedSvc.BuildFeed(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to build home feed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    feed,
	})
}
