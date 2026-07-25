package track

import (
	"errors"
	"log/slog"
	"music/internal/modules/catalog/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListPublic godoc
// @Summary List public tracks
// @Description Returns a paginated list of publicly visible tracks. Supports optional filtering by album_id or artist_id.
// @Tags tracks
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Param album_id query string false "Filter by album ID"
// @Param artist_id query string false "Filter by artist ID"
// @Success 200 {array} TrackResponse
// @Failure 500 {object} map[string]interface{}
// @Router /tracks [get]
func (h *Handler) ListPublic(c *gin.Context) {
	p := common.ParsePagination(c)

	var opts ListOptions
	if albumID := c.Query("album_id"); albumID != "" {
		uid, err := uuid.Parse(albumID)
		if err == nil {
			opts.AlbumID = &uid
		}
	}
	if artistID := c.Query("artist_id"); artistID != "" {
		uid, err := uuid.Parse(artistID)
		if err == nil {
			opts.ArtistID = &uid
		}
	}
	if q := c.Query("q"); q != "" {
		opts.Query = q
	}

	items, err := h.service.List(c.Request.Context(), p.Limit, p.Offset, true, opts)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}

	c.JSON(http.StatusOK, TrackListToResponse(items))
}

// ListAdmin godoc
// @Summary List all tracks
// @Description Returns a paginated list of tracks including non-public entries for administrative access.
// @Tags tracks
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} TrackResponse
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

	c.JSON(http.StatusOK, TrackListToResponse(items))
}

// Get godoc
// @Summary Get track by ID
// @Description Returns a single track by its ID.
// @Tags tracks
// @Produce json
// @Param trackID path string true "Track ID"
// @Success 200 {object} TrackResponse
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
	c.JSON(http.StatusOK, TrackToResponse(item))
}

// Create godoc
// @Summary Create track
// @Description Creates a new track.
// @Tags tracks
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Track creation payload"
// @Success 201 {object} TrackResponse
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

	if req.ArtistID == uuid.Nil && len(req.Artists) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artist_id or artists is required"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if errors.Is(err, common.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track payload"})
		return
	}
	if errors.Is(err, common.ErrForeignKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album_id, artists, or genre_ids"})
		return
	}
	if errors.Is(err, common.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "track conflict"})
		return
	}
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TrackToResponse(item))
}

// Update godoc
// @Summary Update track
// @Description Updates an existing track by ID.
// @Tags tracks
// @Accept json
// @Produce json
// @Param trackID path string true "Track ID"
// @Param request body UpdateRequest true "Track update payload"
// @Success 200 {object} TrackResponse
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
		slog.Error("failed to update track", "trackID", c.Param("trackID"), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, TrackToResponse(item))
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

// Random godoc
// @Summary Get random tracks
// @Description Returns a list of randomly selected public tracks. Useful for shuffle mode.
// @Tags tracks
// @Produce json
// @Param limit query int false "Number of random tracks to return (max 100)"
// @Success 200 {array} TrackResponse
// @Failure 500 {object} map[string]interface{}
// @Router /tracks/random [get]
func (h *Handler) Random(c *gin.Context) {
	limit := common.ParsePagination(c).Limit
	items, err := h.service.Random(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch random tracks"})
		return
	}
	c.JSON(http.StatusOK, TrackListToResponse(items))
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

	c.JSON(http.StatusOK, TrackCreditsToResponse(items))
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

	c.JSON(http.StatusOK, TrackCreditsToResponse(items))
}

func (h *Handler) Artists(c *gin.Context) {
	trackID := c.Param("trackID")

	items, err := h.service.ListArtists(c.Request.Context(), trackID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		case errors.Is(err, common.ErrNotFound):
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
		switch {
		case errors.Is(err, common.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		case errors.Is(err, common.ErrNotFound):
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
