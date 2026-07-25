package stats

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler serves the admin catalog stats endpoint.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get returns aggregate counts for tracks, albums, artists, and genres.
// GET /admin/catalog/stats
func (h *Handler) Get(c *gin.Context) {
	stats, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		slog.Error("failed to get catalog stats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get catalog stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
