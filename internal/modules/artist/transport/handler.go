package artist

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
// @Summary List artists
// @Description Returns a paginated list of artists.
// @Tags artists
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} Artist
// @Failure 500 {object} map[string]interface{}
// @Router /artists [get]
func (h *Handler) List(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artists"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Get godoc
// @Summary Get artist by ID
// @Description Returns a single artist by its ID.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 200 {object} Artist
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [get]
func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create godoc
// @Summary Create artist
// @Description Creates a new artist.
// @Tags artists
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Artist creation payload"
// @Success 201 {object} Artist
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "artist already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create artist"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Update godoc
// @Summary Update artist
// @Description Updates an existing artist by ID.
// @Tags artists
// @Accept json
// @Produce json
// @Param artistID path string true "Artist ID"
// @Param request body UpdateRequest true "Artist update payload"
// @Success 200 {object} Artist
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), c.Param("artistID"), req)
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "artist conflict"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update artist"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary Delete artist
// @Description Deletes an artist by ID.
// @Tags artists
// @Produce json
// @Param artistID path string true "Artist ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /artists/{artistID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete artist"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Overview(c *gin.Context) {
	item, err := h.service.GetOverview(c.Request.Context(), c.Param("artistID"))
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist overview"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) Tracks(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListTracks(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artist tracks"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) Albums(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListAlbums(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artist albums"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) Singles(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListSingles(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artist singles"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) AppearsOn(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListAppearsOn(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artist appears-on tracks"})
		return
	}

	c.JSON(http.StatusOK, items)
}
func (h *Handler) TopTracks(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListTopTracks(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artist top tracks"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) Related(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListRelated(c.Request.Context(), c.Param("artistID"), p.Limit, p.Offset)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list related artists"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) ReplaceRelated(c *gin.Context) {
	var req []RelatedArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceRelated(c.Request.Context(), c.Param("artistID"), req)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id or related artists"})
		return
	}
	if errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace related artists"})
		return
	}

	c.JSON(http.StatusOK, items)
}
func (h *Handler) ReplaceTopTracks(c *gin.Context) {
	artistID := c.Param("artistID")

	var req []ArtistTopTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	items, err := h.service.ReplaceTopTracks(c.Request.Context(), artistID, req)
	if err != nil {
		switch err {
		case common.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		case common.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to replace artist top tracks"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}
