package video

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/modules/auth"
	"music/internal/modules/catalog/common"
)

type Handler struct {
	service     *Service
	commentRepo *CommentRepository
}

func NewHandler(service *Service, commentRepo *CommentRepository) *Handler {
	return &Handler{service: service, commentRepo: commentRepo}
}

// ---------------------------------------------------------------------------
// Official MV (admin only)
// ---------------------------------------------------------------------------

// UploadOfficialMV godoc
// @Summary Upload official music video
// @Description Upload an official music video for a track. Admin only.
// @Tags videos
// @Accept json
// @Produce json
// @Security Bearer
// @Param trackID path string true "Track ID"
// @Param request body CreateVideoRequest true "Video upload payload"
// @Success 201 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/tracks/{trackId}/mv [post]
func (h *Handler) UploadOfficialMV(c *gin.Context) {
	trackIDStr := c.Param("trackId")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid track id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("invalid user", nil))
		return
	}

	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", err))
		return
	}

	// Override track_id from the path
	req.TrackID = trackIDStr

	video, err := h.service.UploadOfficialMV(c.Request.Context(), trackID, userID, req)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to upload official MV", err))
		return
	}

	response.Success[any](c, http.StatusCreated, "official MV uploaded", toVideoResponse(video))
}

// DeleteVideo godoc
// @Summary Delete a video
// @Description Delete a video by ID. Admin can delete any; users can delete own edits.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/videos/{id} [delete]
func (h *Handler) DeleteVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, _ := uuid.Parse(userIDStr)
	role := auth.UserRoleFromContext(c)
	isAdmin := role == "admin"

	if err := h.service.DeleteVideo(c.Request.Context(), id, userID, isAdmin); err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		case errors.Is(err, ErrNotOwner), errors.Is(err, ErrRequiresAdmin):
			response.Error(c, apperrors.Forbidden("not authorized to delete this video", nil))
		default:
			response.Error(c, apperrors.Internal("failed to delete video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video deleted", nil)
}

// ApproveVideo godoc
// @Summary Approve a user edit video
// @Description Approve a user-uploaded edit so it appears in explore. Admin only.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/videos/{id}/approve [post]
func (h *Handler) ApproveVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	video, err := h.service.ApproveVideo(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		default:
			response.Error(c, apperrors.Internal("failed to approve video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video approved", toVideoResponse(video))
}

// AdminListVideos godoc
// @Summary List all videos (admin)
// @Description Paginated list of ALL videos — approved, pending, public, private. Admin only.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param limit query int false "Maximum items (default 20, max 100)"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} response.SuccessResponseData
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/videos [get]
func (h *Handler) AdminListVideos(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListAll(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to list videos", err))
		return
	}

	resp := make([]VideoResponse, 0, len(items))
	for _, v := range items {
		resp = append(resp, toVideoResponse(&v))
	}

	response.Success[any](c, http.StatusOK, "videos list", ExploreResponse{
		Items:  resp,
		Limit:  p.Limit,
		Offset: p.Offset,
		Count:  len(resp),
	})
}

// AdminUpdateVideo godoc
// @Summary Update a video (admin)
// @Description Partially update video metadata. Admin only.
// @Tags videos
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Param request body AdminUpdateVideoRequest true "Video update payload"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/videos/{id} [put]
func (h *Handler) AdminUpdateVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	var req AdminUpdateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", err))
		return
	}

	video, err := h.service.AdminUpdateVideo(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		default:
			response.Error(c, apperrors.Internal("failed to update video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video updated", toVideoResponse(video))
}

// ---------------------------------------------------------------------------
// User edits
// ---------------------------------------------------------------------------

// CreateUserEdit godoc
// @Summary Upload a user edit video
// @Description Upload a user-created video edit for a track. Requires authentication.
// @Tags videos
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateVideoRequest true "Video upload payload"
// @Success 201 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos [post]
func (h *Handler) CreateUserEdit(c *gin.Context) {
	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", err))
		return
	}

	video, err := h.service.CreateUserEdit(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, apperrors.BadRequest("invalid input", nil))
		default:
			response.Error(c, apperrors.Internal("failed to create user edit", err))
		}
		return
	}

	response.Success[any](c, http.StatusCreated, "user edit uploaded", toVideoResponse(video))
}

// GetVideo godoc
// @Summary Get a single video
// @Description Get full details of a video by ID.
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos/{id} [get]
func (h *Handler) GetVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	video, err := h.service.GetVideo(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		default:
			response.Error(c, apperrors.Internal("failed to get video", err))
		}
		return
	}

	// Look up track cover for poster fallback when video has no dedicated thumbnail
	coverURL, _ := h.service.GetTrackCoverURL(c.Request.Context(), video.TrackID)

	response.Success[any](c, http.StatusOK, "video found", toVideoResponse(video, coverURL))
}

// StreamVideo godoc
// @Summary Stream video file
// @Description Stream the video file for playback. Returns the raw video file.
// @Tags videos
// @Produce octet-stream
// @Param id path string true "Video ID"
// @Success 200 {file} binary "Video file stream"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /videos/{id}/stream [get]
func (h *Handler) StreamVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	video, err := h.service.GetVideo(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		default:
			response.Error(c, apperrors.Internal("failed to get video", err))
		}
		return
	}

	// Determine which file to serve
	// Priority: final processed video > raw uploaded video
	if video.FinalVideoPath != nil && *video.FinalVideoPath != "" {
		c.Header("Accept-Ranges", "bytes")
		c.Header("Cache-Control", "public, max-age=3600, immutable")
		c.Header("Cross-Origin-Resource-Policy", "cross-origin")
		c.File(*video.FinalVideoPath)
		return
	}

	if video.RawVideoPath != nil && *video.RawVideoPath != "" {
		// Raw video is a URL (e.g. /uploads/video/uuid-filename.mp4) served by the media handler.
		// Redirect the browser so it fetches through the static file server.
		c.Redirect(http.StatusFound, *video.RawVideoPath)
		return
	}

	response.Error(c, apperrors.NotFound("video file not available", nil))
}

// DeleteOwnVideo godoc
// @Summary Delete own video
// @Description Delete a video uploaded by the authenticated user.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos/{id} [delete]
func (h *Handler) DeleteOwnVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	if err := h.service.DeleteOwnVideo(c.Request.Context(), id, userID); err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		case errors.Is(err, ErrNotOwner):
			response.Error(c, apperrors.Forbidden("not the owner of this video", nil))
		default:
			response.Error(c, apperrors.Internal("failed to delete video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video deleted", nil)
}

// ListByTrack godoc
// @Summary List videos for a track
// @Description Returns all videos (MV + edits) for a track.
// @Tags videos
// @Produce json
// @Param trackId path string true "Track ID"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tracks/{trackId}/videos [get]
func (h *Handler) ListByTrack(c *gin.Context) {
	trackIDStr := c.Param("trackId")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid track id", err))
		return
	}

	items, err := h.service.ListByTrack(c.Request.Context(), trackID)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to list videos", err))
		return
	}

	// Single track cover lookup (all videos share the same track)
	coverURL, _ := h.service.GetTrackCoverURL(c.Request.Context(), trackID)

	resp := make([]VideoResponse, 0, len(items))
	for _, v := range items {
		resp = append(resp, toVideoResponse(&v, coverURL))
	}

	response.Success[any](c, http.StatusOK, "videos found", resp)
}

// ---------------------------------------------------------------------------
// Explore
// ---------------------------------------------------------------------------

// ListExplore godoc
// @Summary Explore videos
// @Description Paginated list of approved + public videos. Public endpoint.
// @Tags videos
// @Produce json
// @Param limit query int false "Maximum items (default 20, max 100)"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} response.SuccessResponseData
// @Failure 500 {object} response.ErrorResponse
// @Router /explore/videos [get]
func (h *Handler) ListExplore(c *gin.Context) {
	p := common.ParsePagination(c)

	items, err := h.service.ListExplore(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to list explore videos", err))
		return
	}

	// Batch-load track covers for poster fallback
	trackIDs := make([]uuid.UUID, 0, len(items))
	for _, v := range items {
		trackIDs = append(trackIDs, v.TrackID)
	}
	covers, _ := h.service.GetTrackCovers(c.Request.Context(), trackIDs)

	resp := make([]VideoResponse, 0, len(items))
	for _, v := range items {
		resp = append(resp, toVideoResponse(&v, covers[v.TrackID]))
	}

	response.Success[any](c, http.StatusOK, "explore videos", ExploreResponse{
		Items:  resp,
		Limit:  p.Limit,
		Offset: p.Offset,
		Count:  len(resp),
	})
}

// ---------------------------------------------------------------------------
// User videos
// ---------------------------------------------------------------------------

// GetUserVideos godoc
// @Summary List user's uploaded videos
// @Description Returns paginated videos uploaded by a specific user.
// @Tags videos
// @Produce json
// @Param id path string true "User ID"
// @Param limit query int false "Maximum items (default 20, max 100)"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} response.SuccessResponseData
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id}/videos [get]
func (h *Handler) GetUserVideos(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid user id", nil))
		return
	}

	p := common.ParsePagination(c)

	// Determine viewer (optional, for private video access)
	var viewerID uuid.UUID
	if viewerStr := auth.UserIDFromContext(c); viewerStr != "" {
		viewerID, _ = uuid.Parse(viewerStr)
	}

	items, err := h.service.ListByUser(c.Request.Context(), userID, viewerID, p.Limit, p.Offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to list user videos", err))
		return
	}

	// Batch-load track covers for poster fallback
	trackIDs := make([]uuid.UUID, 0, len(items))
	for _, v := range items {
		trackIDs = append(trackIDs, v.TrackID)
	}
	covers, _ := h.service.GetTrackCovers(c.Request.Context(), trackIDs)

	resp := make([]VideoResponse, 0, len(items))
	for _, v := range items {
		resp = append(resp, toVideoResponse(&v, covers[v.TrackID]))
	}

	response.Success[any](c, http.StatusOK, "user videos", ExploreResponse{
		Items:  resp,
		Limit:  p.Limit,
		Offset: p.Offset,
		Count:  len(resp),
	})
}

// ---------------------------------------------------------------------------
// Social — Video likes
// ---------------------------------------------------------------------------

// LikeVideo godoc
// @Summary Like a video
// @Description Like a video by ID.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos/{id}/like [post]
func (h *Handler) LikeVideo(c *gin.Context) {
	idStr := c.Param("id")
	videoID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	if err := h.service.LikeVideo(c.Request.Context(), videoID, userID); err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		case errors.Is(err, ErrAlreadyLiked):
			response.Success[any](c, http.StatusOK, "already liked", nil)
		default:
			response.Error(c, apperrors.Internal("failed to like video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video liked", nil)
}

// UnlikeVideo godoc
// @Summary Unlike a video
// @Description Remove like from a video by ID.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos/{id}/like [delete]
func (h *Handler) UnlikeVideo(c *gin.Context) {
	idStr := c.Param("id")
	videoID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	if err := h.service.UnlikeVideo(c.Request.Context(), videoID, userID); err != nil {
		switch {
		case errors.Is(err, ErrNotLiked):
			response.Success[any](c, http.StatusOK, "not liked", nil)
		default:
			response.Error(c, apperrors.Internal("failed to unlike video", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "video unliked", nil)
}

// ViewVideo godoc
// @Summary Increment video view count
// @Description Increment the view count for a video.
// @Tags videos
// @Produce json
// @Security Bearer
// @Param id path string true "Video ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /videos/{id}/view [post]
func (h *Handler) ViewVideo(c *gin.Context) {
	idStr := c.Param("id")
	videoID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	if err := h.service.ViewVideo(c.Request.Context(), videoID); err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, apperrors.NotFound("video not found", nil))
		default:
			response.Error(c, apperrors.Internal("failed to record view", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "view recorded", nil)
}

// ---------------------------------------------------------------------------
// Track Like Visibility
// ---------------------------------------------------------------------------

// SetTrackLikeVisibility godoc
// @Summary Set track like visibility
// @Description Set whether a user's like on a track is public or private.
// @Tags library
// @Accept json
// @Produce json
// @Security Bearer
// @Param trackId path string true "Track ID"
// @Param request body VideoVisibilityRequest true "Visibility setting"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /tracks/{trackId}/like/visibility [put]
func (h *Handler) SetTrackLikeVisibility(c *gin.Context) {
	trackIDStr := c.Param("trackId")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid track id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	var req VideoVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", err))
		return
	}

	tlv, err := h.service.SetTrackLikeVisibility(c.Request.Context(), userID, trackID, req.Visibility)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, apperrors.BadRequest("invalid visibility value", nil))
		default:
			response.Error(c, apperrors.Internal("failed to set visibility", err))
		}
		return
	}

	response.Success[any](c, http.StatusOK, "visibility updated", gin.H{
		"track_id":   tlv.TrackID.String(),
		"visibility": tlv.Visibility,
	})
}

// GetPublicLikedTracks godoc
// @Summary Get public liked tracks
// @Description Returns tracks a user has publicly liked (visible on their profile).
// @Tags users
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Maximum items (default 20, max 100)"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/liked-tracks [get]
func (h *Handler) GetPublicLikedTracks(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid user id", err))
		return
	}

	p := common.ParsePagination(c)

	items, err := h.service.GetPublicLikedTracks(c.Request.Context(), userID, p.Limit, p.Offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to get liked tracks", err))
		return
	}

	resp := make([]LikedTrackResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, LikedTrackResponse{
			TrackID:    item.TrackID.String(),
			Visibility: item.Visibility,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		})
	}

	response.Success[any](c, http.StatusOK, "public liked tracks", gin.H{
		"items":  resp,
		"limit":  p.Limit,
		"offset": p.Offset,
		"count":  len(resp),
	})
}

// ---------------------------------------------------------------------------
// Music Status
// ---------------------------------------------------------------------------

// SetMusicStatus godoc
// @Summary Set current music status
// @Description Set what track you are currently listening to and its visibility.
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateMusicStatusRequest true "Music status payload"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/me/music-status [put]
func (h *Handler) SetMusicStatus(c *gin.Context) {
	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	var req UpdateMusicStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", err))
		return
	}

	// Default visibility to private if not set
	if req.Visibility == "" {
		req.Visibility = "private"
	}

	status, err := h.service.SetMusicStatus(c.Request.Context(), userID, req.TrackID, req.Visibility)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, apperrors.BadRequest("invalid input", nil))
		default:
			response.Error(c, apperrors.Internal("failed to set music status", err))
		}
		return
	}

	var trackIDStr *string
	if status.CurrentTrackID != nil {
		s := status.CurrentTrackID.String()
		trackIDStr = &s
	}

	response.Success[any](c, http.StatusOK, "music status updated", gin.H{
		"current_track_id": trackIDStr,
		"visibility":       status.Visibility,
		"updated_at":       status.UpdatedAt.Format(time.RFC3339),
	})
}

// ClearMusicStatus godoc
// @Summary Clear current music status
// @Description Stop broadcasting what you are listening to.
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/me/music-status [delete]
func (h *Handler) ClearMusicStatus(c *gin.Context) {
	userIDStr := auth.UserIDFromContext(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	if err := h.service.ClearMusicStatus(c.Request.Context(), userID); err != nil {
		response.Error(c, apperrors.Internal("failed to clear music status", err))
		return
	}

	response.Success[any](c, http.StatusOK, "music status cleared", nil)
}

// GetMusicStatus godoc
// @Summary Get user music status
// @Description Get what track a user is currently listening to. Privacy-aware.
// @Tags users
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} response.SuccessResponseData
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/music-status [get]
func (h *Handler) GetMusicStatus(c *gin.Context) {
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid user id", err))
		return
	}

	// Get requesting user (optional — might not be authenticated)
	var requestingUserID *uuid.UUID
	if rawID := auth.UserIDFromContext(c); rawID != "" {
		if pid, err := uuid.Parse(rawID); err == nil {
			requestingUserID = &pid
		}
	}

	status, err := h.service.GetMusicStatus(c.Request.Context(), targetUserID, requestingUserID)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to get music status", err))
		return
	}

	if status == nil {
		response.Success[any](c, http.StatusOK, "music status", MusicStatusResponse{Playing: false})
		return
	}

	resp := MusicStatusResponse{
		Playing:   true,
		UpdatedAt: strPtr(status.UpdatedAt.Format(time.RFC3339)),
	}
	if status.CurrentTrackID != nil {
		s := status.CurrentTrackID.String()
		resp.CurrentTrackID = &s
	}

	response.Success[any](c, http.StatusOK, "music status", resp)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func toVideoResponse(v *Video, trackCoverURL ...*string) VideoResponse {
	var finalPath, thumbPath *string
	if v.FinalVideoPath != nil {
		finalPath = v.FinalVideoPath
	}
	if v.ThumbnailPath != nil {
		thumbPath = v.ThumbnailPath
	}

	var coverURL *string
	if len(trackCoverURL) > 0 {
		coverURL = trackCoverURL[0]
	}

	return VideoResponse{
		ID:             v.ID.String(),
		TrackID:        v.TrackID.String(),
		UploaderID:     v.UploaderID.String(),
		Type:           string(v.Type),
		Status:         string(v.Status),
		Title:          v.Title,
		Description:    v.Description,
		FinalVideoPath: finalPath,
		ThumbnailPath:  thumbPath,
		TrackCoverURL:  coverURL,
		DurationMs:     v.DurationMs,
		AspectRatio:    v.AspectRatio,
		ViewCount:      v.ViewCount,
		LikeCount:      v.LikeCount,
		IsPublic:       v.IsPublic,
		IsApproved:     v.IsApproved,
		TrackStartMs:   v.TrackStartMs,
		TrackEndMs:     v.TrackEndMs,
		CreatedAt:      v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      v.UpdatedAt.Format(time.RFC3339),
	}
}

func strPtr(s string) *string {
	return &s
}

// ---------------------------------------------------------------------------
// Internal: Video processing webhook callback (from Python ML service)
// ---------------------------------------------------------------------------

// VideoCallbackPayload is the JSON body sent by the Python ML service
// after processing a video job (audio replacement).
type VideoCallbackPayload struct {
	VideoID          string  `json:"video_id" binding:"required"`
	FinalVideoPath   string  `json:"final_video_path"`
	ThumbnailPath    string  `json:"thumbnail_path"`
	DurationMs       int64   `json:"duration_ms"`
	AspectRatio      string  `json:"aspect_ratio"`
	ProcessingStatus string  `json:"processing_status" binding:"required"` // "completed" | "failed"
	ErrorMessage     *string `json:"error_message"`
}

// HandleVideoCallback receives the result of a video processing job
// from the Python ML service (audio replacement).
// This endpoint is protected by HMAC signature verification, not JWT auth.
func (h *Handler) HandleVideoCallback(c *gin.Context) {
	var payload VideoCallbackPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	video, err := h.service.HandleVideoCallback(c.Request.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process video callback"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "received",
		"video_id": video.ID.String(),
	})
}
