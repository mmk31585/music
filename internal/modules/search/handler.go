package search

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Search godoc
// @Summary Search catalog
// @Description Searches across tracks, albums, artists, and playlists using a query string.
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of results" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /search [get]
func (h *Handler) Search(c *gin.Context) {

	query := strings.TrimSpace(c.Query("q"))
	limit := 10

	if rawLimit := c.Query("limit"); rawLimit != "" {
		if parsed, err := strconv.Atoi(rawLimit); err == nil {
			limit = parsed
		}
	}

	result, err := h.service.Search(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "search failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Search godoc
// @Summary Search across tracks, albums, artists and playlists
// @Description Performs a search using query string and optional limit
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of results"
// @Success 200 {object} response.SuccessResponse{data=search.SearchResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /search [get]
