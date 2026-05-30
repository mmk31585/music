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

// LikeTrack godoc
// @Summary Like a track
// @Description Adds a track to the authenticated user's liked tracks.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body LikeTrackRequest true "Like track request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/tracks/like [post]
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

// UnlikeTrack godoc
// @Summary Unlike a track
// @Description Removes a track from the authenticated user's liked tracks.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param trackId path int true "Track ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/tracks/{trackId}/like [delete]
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

// ListLikedTracks godoc
// @Summary List liked tracks
// @Description Returns all liked tracks for the authenticated user.
// @Tags library
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/tracks/liked [get]
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

// LikeAlbum godoc
// @Summary Like an album
// @Description Adds an album to the authenticated user's liked albums.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body LikeAlbumRequest true "Like album request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/albums/like [post]
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

// UnlikeAlbum godoc
// @Summary Unlike an album
// @Description Removes an album from the authenticated user's liked albums.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param albumId path int true "Album ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/albums/{albumId}/like [delete]
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

// ListLikedAlbums godoc
// @Summary List liked albums
// @Description Returns all liked albums for the authenticated user.
// @Tags library
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/albums/liked [get]
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

// FollowArtist godoc
// @Summary Follow an artist from library
// @Description Adds an artist to the authenticated user's followed artists.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body FollowArtistRequest true "Follow artist request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/artists/follow [post]
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

// UnfollowArtist godoc
// @Summary Unfollow an artist from library
// @Description Removes an artist from the authenticated user's followed artists.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param artistId path int true "Artist ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/artists/{artistId}/follow [delete]
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

// ListFollowedArtists godoc
// @Summary List followed artists
// @Description Returns all followed artists for the authenticated user.
// @Tags library
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/artists/followed [get]
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

// AddPlayHistory godoc
// @Summary Add play history entry
// @Description Adds a track to the authenticated user's play history.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body AddPlayHistoryRequest true "Add play history request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/history [post]
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

// ListPlayHistory godoc
// @Summary List play history
// @Description Returns the authenticated user's play history.
// @Tags library
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/history [get]
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

// ListRecentlyPlayed godoc
// @Summary List recently played tracks
// @Description Returns recently played tracks for the authenticated user.
// @Tags library
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum number of items to return" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /library/history/recent [get]
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
