package playlist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePlaylist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
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

func (h *Handler) UpdatePlaylist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := parseInt64Param(c, "id")
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

func (h *Handler) DeletePlaylist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := parseInt64Param(c, "id")
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

func (h *Handler) GetPlaylist(c *gin.Context) {
	playlistID, err := parseInt64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	var requesterID *int64
	if userID, ok := getUserIDFromGin(c); ok {
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

func (h *Handler) ListMyPlaylists(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
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

func (h *Handler) AddTrack(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := parseInt64Param(c, "id")
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

func (h *Handler) RemoveTrack(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := parseInt64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid playlist id"})
		return
	}

	trackID, err := parseInt64Param(c, "trackId")
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

func (h *Handler) ReorderTrack(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	playlistID, err := parseInt64Param(c, "id")
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

func parseInt64Param(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

// Replace this with your real auth helper.
func getUserIDFromGin(c *gin.Context) (int64, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := v.(int64)
	return id, ok
}
