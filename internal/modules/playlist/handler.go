package playlist

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreatePlaylist godoc
// @Summary Create playlist
// @Description Creates a new playlist for the authenticated user.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreatePlaylistRequest true "Create playlist request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists [post]
func (h *Handler) CreatePlaylist(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	p, err := h.service.CreatePlaylist(c.Request.Context(), req, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ToPlaylistResponse(p, nil),
	})
}

// UpdatePlaylist godoc
// @Summary Update playlist
// @Description Updates an existing playlist owned by the authenticated user.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Playlist ID"
// @Param request body UpdatePlaylistRequest true "Update playlist request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id} [put]
func (h *Handler) UpdatePlaylist(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	p, err := h.service.UpdatePlaylist(c.Request.Context(), playlistID, userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ToPlaylistResponse(p, nil),
	})
}

// DeletePlaylist godoc
// @Summary Delete playlist
// @Description Deletes an existing playlist owned by the authenticated user.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Playlist ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id} [delete]
func (h *Handler) DeletePlaylist(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	if err := h.service.DeletePlaylist(c.Request.Context(), playlistID, userID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "playlist deleted",
	})
}

// GetPlaylist godoc
// @Summary Get playlist
// @Description Returns playlist details and tracks. Private playlists require owner access.
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id} [get]
func (h *Handler) GetPlaylist(c *gin.Context) {
	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	var requesterID *uuid.UUID
	if userID, ok := web.GetRequiredUserUUID(c); ok {
		requesterID = &userID
	}

	p, tracks, err := h.service.GetPlaylist(c.Request.Context(), playlistID, requesterID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ToPlaylistResponse(p, tracks),
	})
}

// ListPublicPlaylists godoc
// @Summary List public playlists
// @Description Returns all public playlists.
// @Tags playlists
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/public [get]
func (h *Handler) ListPublicPlaylists(c *gin.Context) {
	items, err := h.service.ListPublicPlaylists(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

// ListMyPlaylists godoc
// @Summary List my playlists
// @Description Returns playlists owned by the authenticated user.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/me [get]
func (h *Handler) ListMyPlaylists(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	items, err := h.service.ListMyPlaylists(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

// AddTrack godoc
// @Summary Add track to playlist
// @Description Adds a track to an authenticated user's playlist.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Playlist ID"
// @Param request body AddTrackRequest true "Add track request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id}/tracks [post]
func (h *Handler) AddTrack(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	var req AddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.AddTrack(c.Request.Context(), playlistID, userID, req.TrackID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "track added to playlist",
	})
}

// RemoveTrack godoc
// @Summary Remove track from playlist
// @Description Removes a track from an authenticated user's playlist.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Playlist ID"
// @Param trackId path string true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id}/tracks/{trackId} [delete]
func (h *Handler) RemoveTrack(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	trackID, err := uuid.Parse(c.Param("trackId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid track id"})
		return
	}

	if err := h.service.RemoveTrack(c.Request.Context(), playlistID, userID, trackID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "track removed from playlist",
	})
}

// ReorderTrack godoc
// @Summary Reorder playlist track
// @Description Updates the position of a track inside an authenticated user's playlist.
// @Tags playlists
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Playlist ID"
// @Param request body ReorderTrackRequest true "Reorder track request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /playlists/{id}/tracks/reorder [put]
func (h *Handler) ReorderTrack(c *gin.Context) {
	userID, ok := web.GetRequiredUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	var req ReorderTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.ReorderTrack(c.Request.Context(), playlistID, userID, req.TrackID, req.NewPosition); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "track reordered",
	})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidPlaylistName):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrPlaylistNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrPlaylistTrackNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackAlreadyInPlaylist):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrForbiddenPlaylistAccess):
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidTrackPosition):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}
