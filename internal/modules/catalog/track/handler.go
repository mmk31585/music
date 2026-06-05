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

// ListPublic godoc
// @Summary List public tracks
// @Description Returns a paginated list of publicly visible tracks.
// @Tags tracks
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} Track
// @Failure 500 {object} map[string]interface{}
// @Router /tracks [get]
func (h *Handler) ListPublic(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset, true)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// ListAdmin godoc
// @Summary List all tracks
// @Description Returns a paginated list of tracks including non-public entries for administrative access.
// @Tags tracks
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} Track
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/tracks [get]
func (h *Handler) ListAdmin(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Get godoc
// @Summary Get track by ID
// @Description Returns a single track by its ID.
// @Tags tracks
// @Produce json
// @Param trackID path string true "Track ID"
// @Success 200 {object} Track
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracks/{trackID} [get]
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

// Create godoc
// @Summary Create track
// @Description Creates a new track.
// @Tags tracks
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Track creation payload"
// @Success 201 {object} Track
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracks [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album_id, artists, or genre_ids"})
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

// Update godoc
// @Summary Update track
// @Description Updates an existing track by ID.
// @Tags tracks
// @Accept json
// @Produce json
// @Param trackID path string true "Track ID"
// @Param request body UpdateRequest true "Track update payload"
// @Success 200 {object} Track
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracks/{trackID} [put]
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

// Delete godoc
// @Summary Delete track
// @Description Deletes a track by ID.
// @Tags tracks
// @Produce json
// @Param trackID path string true "Track ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracks/{trackID} [delete]
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
func (h *Handler) Credits(c *gin.Context) {
	items, err := h.service.ListCredits(c.Request.Context(), c.Param("trackID"))
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list track credits"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) ReplaceCredits(c *gin.Context) {
	var req []TrackCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceCredits(c.Request.Context(), c.Param("trackID"), req)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id or credits"})
		return
	}
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace track credits"})
		return
	}

	c.JSON(http.StatusOK, items)
}
func (h *Handler) Artists(c *gin.Context) {
	trackID := c.Param("trackID")

	items, err := h.service.ListArtists(c.Request.Context(), trackID)
	if err != nil {
		switch err {
		case common.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		case common.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to fetch track artists"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

func (h *Handler) ReplaceArtists(c *gin.Context) {
	trackID := c.Param("trackID")

	var req []TrackArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceArtists(c.Request.Context(), trackID, req)
	if err != nil {
		switch err {
		case common.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		case common.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to replace track artists"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}
