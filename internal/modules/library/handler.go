package library

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

func (h *Handler) LikeTrack(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req LikeTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.LikeTrack(c.Request.Context(), userID, req.TrackID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "track liked"})
}

func (h *Handler) UnlikeTrack(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	trackID, err := parseInt64Param(c, "trackId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid track id"})
		return
	}

	if err := h.service.UnlikeTrack(c.Request.Context(), userID, trackID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "track unliked"})
}

func (h *Handler) ListLikedTracks(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	items, err := h.service.ListLikedTracks(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) LikeAlbum(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req LikeAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.LikeAlbum(c.Request.Context(), userID, req.AlbumID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "album liked"})
}

func (h *Handler) UnlikeAlbum(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	albumID, err := parseInt64Param(c, "albumId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid album id"})
		return
	}

	if err := h.service.UnlikeAlbum(c.Request.Context(), userID, albumID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "album unliked"})
}

func (h *Handler) ListLikedAlbums(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	items, err := h.service.ListLikedAlbums(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) FollowArtist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req FollowArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.FollowArtist(c.Request.Context(), userID, req.ArtistID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "artist followed"})
}

func (h *Handler) UnfollowArtist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	artistID, err := parseInt64Param(c, "artistId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid artist id"})
		return
	}

	if err := h.service.UnfollowArtist(c.Request.Context(), userID, artistID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "artist unfollowed"})
}

func (h *Handler) ListFollowedArtists(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	items, err := h.service.ListFollowedArtists(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) AddPlayHistory(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var req AddPlayHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	if err := h.service.AddPlayHistory(c.Request.Context(), userID, req.TrackID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "play history added"})
}

func (h *Handler) ListPlayHistory(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	limit := parseLimit(c, 50)

	items, err := h.service.ListPlayHistory(c.Request.Context(), userID, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) ListRecentlyPlayed(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	limit := parseLimit(c, 20)

	items, err := h.service.ListRecentlyPlayed(c.Request.Context(), userID, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrTrackNotFound),
		errors.Is(err, ErrAlbumNotFound),
		errors.Is(err, ErrArtistNotFound),
		errors.Is(err, ErrLikedTrackNotFound),
		errors.Is(err, ErrLikedAlbumNotFound),
		errors.Is(err, ErrFollowedArtistNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrTrackAlreadyLiked),
		errors.Is(err, ErrAlbumAlreadyLiked),
		errors.Is(err, ErrArtistAlreadyFollowed):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, ErrInvalidLimit):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}

func getUserIDFromGin(c *gin.Context) (int64, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := v.(int64)
	return id, ok
}

func parseInt64Param(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

func parseLimit(c *gin.Context, defaultValue int) int {
	v := c.Query("limit")
	if v == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}
