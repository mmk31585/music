package genre

import (
	"errors"
	"music/internal/modules/catalog/transport"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List godoc
// @Summary List genres
// @Description Returns a paginated list of genres.
// @Tags genres
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} Genre
// @Failure 500 {object} map[string]interface{}
// @Router /genres [get]
func (h *Handler) List(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list genres"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Get godoc
// @Summary Get genre by ID
// @Description Returns a single genre by its ID.
// @Tags genres
// @Produce json
// @Param genreID path string true "Genre ID"
// @Success 200 {object} Genre
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /genres/{genreID} [get]
func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("genreID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create godoc
// @Summary Create genre
// @Description Creates a new genre.
// @Tags genres
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Genre creation payload"
// @Success 201 {object} Genre
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /genres [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "genre already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create genre"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Update godoc
// @Summary Update genre
// @Description Updates an existing genre by ID.
// @Tags genres
// @Accept json
// @Produce json
// @Param genreID path string true "Genre ID"
// @Param request body UpdateRequest true "Genre update payload"
// @Success 200 {object} Genre
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /genres/{genreID} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), c.Param("genreID"), req)
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "genre conflict"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update genre"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary Delete genre
// @Description Deletes a genre by ID.
// @Tags genres
// @Produce json
// @Param genreID path string true "Genre ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /genres/{genreID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("genreID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete genre"})
		return
	}
	c.Status(http.StatusNoContent)
}
