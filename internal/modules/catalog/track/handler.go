package track

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

func (h *Handler) ListPublic(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) ListAdmin(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("trackID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get track"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist_id, album_id, or genre_ids"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "track conflict"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create track"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), c.Param("trackID"), req)
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}
	if errors.Is(err, common.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album_id or genre_ids"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update track"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("trackID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete track"})
		return
	}
	c.Status(http.StatusNoContent)
}
