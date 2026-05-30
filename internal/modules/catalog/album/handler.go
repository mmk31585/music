package album

import (
	"errors"
	"music/internal/modules/catalog/common"
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
// @Summary List albums
// @Description Returns a paginated list of albums.
// @Tags albums
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} Album
// @Failure 500 {object} map[string]interface{}
// @Router /albums [get]
func (h *Handler) List(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list albums"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Get godoc
// @Summary Get album by ID
// @Description Returns a single album by its ID.
// @Tags albums
// @Produce json
// @Param albumID path string true "Album ID"
// @Success 200 {object} Album
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /albums/{albumID} [get]
func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("albumID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get album"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create godoc
// @Summary Create album
// @Description Creates a new album.
// @Tags albums
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Album creation payload"
// @Success 201 {object} Album
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /albums [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist_id"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "album conflict"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release_date"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create album"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Update godoc
// @Summary Update album
// @Description Updates an existing album by ID.
// @Tags albums
// @Accept json
// @Produce json
// @Param albumID path string true "Album ID"
// @Param request body UpdateRequest true "Album update payload"
// @Success 200 {object} Album
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /albums/{albumID} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), c.Param("albumID"), req)
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "album conflict"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update album"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary Delete album
// @Description Deletes an album by ID.
// @Tags albums
// @Produce json
// @Param albumID path string true "Album ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /albums/{albumID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("albumID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete album"})
		return
	}
	c.Status(http.StatusNoContent)
}
