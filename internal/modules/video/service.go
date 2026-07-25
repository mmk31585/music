package video

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StalenessThreshold is the TTL for music status — if updated_at is older than
// this, the user is considered "not playing" even if current_track_id is set.
const StalenessThreshold = 30 * time.Minute

// Common errors.
var (
	ErrVideoNotFound = errors.New("video not found")
	ErrNotOwner      = errors.New("not the owner of this video")
	ErrRequiresAdmin = errors.New("admin role required")
	ErrRequiresAuth  = errors.New("authentication required")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyLiked  = errors.New("already liked")
	ErrNotLiked      = errors.New("not liked")
	ErrTrackNotFound = errors.New("track not found")
	ErrUserNotFound  = errors.New("user not found")
)

// RepositoryInterface defines the data access methods needed by the Service.
type RepositoryInterface interface {
	Create(ctx context.Context, v *Video) error
	GetByID(ctx context.Context, id uuid.UUID) (*Video, error)
	Update(ctx context.Context, v *Video) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTrack(ctx context.Context, trackID uuid.UUID) ([]Video, error)
	ListAll(ctx context.Context, limit, offset int) ([]Video, error)
	ListExplore(ctx context.Context, limit, offset int) ([]Video, error)
	ListByUser(ctx context.Context, userID, viewerID uuid.UUID, limit, offset int) ([]Video, error)
	CreateLike(ctx context.Context, like *VideoLike) error
	DeleteLike(ctx context.Context, videoID, userID uuid.UUID) error
	GetLike(ctx context.Context, videoID, userID uuid.UUID) (*VideoLike, error)
	IncrementLikeCount(ctx context.Context, videoID uuid.UUID) error
	DecrementLikeCount(ctx context.Context, videoID uuid.UUID) error
	IncrementViewCount(ctx context.Context, videoID uuid.UUID) error
	CheckTrackLikeExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error)
	UpsertTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID, visibility string) error
	GetTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID) (*TrackLikeVisibility, error)
	GetPublicLikedTracks(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TrackLikeVisibility, error)
	UpsertMusicStatus(ctx context.Context, status *UserMusicStatus) error
	DeleteMusicStatus(ctx context.Context, userID uuid.UUID) error
	GetMusicStatus(ctx context.Context, userID uuid.UUID) (*UserMusicStatus, error)
	GetTrackCoverURL(ctx context.Context, trackID uuid.UUID) (*string, error)
	GetTrackCovers(ctx context.Context, trackIDs []uuid.UUID) (map[uuid.UUID]*string, error)
}

type FollowService interface {
	IsFollowingUser(ctx context.Context, followerID uuid.UUID, targetUserID string) (bool, error)
}

type Service struct {
	repo          RepositoryInterface
	followService FollowService
}

func NewService(repo RepositoryInterface, followService FollowService) *Service {
	return &Service{
		repo:          repo,
		followService: followService,
	}
}

// ---------------------------------------------------------------------------
// Official MV (admin only)
// ---------------------------------------------------------------------------

func (s *Service) UploadOfficialMV(ctx context.Context, trackID, uploaderID uuid.UUID, req CreateVideoRequest) (*Video, error) {
	v := &Video{
		TrackID:       trackID,
		UploaderID:    uploaderID,
		Type:          VideoTypeOfficialMV,
		Status:        VideoStatusProcessing,
		Title:         req.Title,
		Description:   req.Description,
		RawVideoPath:  &req.RawVideoURL,
		DurationMs:    req.DurationMs,
		AspectRatio:   req.AspectRatio,
		FileSizeBytes: req.FileSizeBytes,
		TrackStartMs:  req.TrackStartMs,
		TrackEndMs:    req.TrackEndMs,
		IsPublic:      true,
		IsApproved:    true,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) DeleteVideo(ctx context.Context, id uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if v == nil {
		return ErrVideoNotFound
	}

	// Admin can delete any video; regular users can only delete their own user edits.
	if !isAdmin {
		if v.UploaderID != requesterID {
			return ErrNotOwner
		}
		if v.Type != VideoTypeUserEdit {
			return ErrRequiresAdmin
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) ApproveVideo(ctx context.Context, id uuid.UUID) (*Video, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrVideoNotFound
	}
	v.IsApproved = true
	v.Status = VideoStatusReady
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// User edits
// ---------------------------------------------------------------------------

func (s *Service) CreateUserEdit(ctx context.Context, uploaderID uuid.UUID, req CreateVideoRequest) (*Video, error) {
	trackID, err := uuid.Parse(req.TrackID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	v := &Video{
		TrackID:       trackID,
		UploaderID:    uploaderID,
		Type:          VideoTypeUserEdit,
		Status:        VideoStatusProcessing,
		Title:         req.Title,
		Description:   req.Description,
		RawVideoPath:  &req.RawVideoURL,
		DurationMs:    req.DurationMs,
		AspectRatio:   req.AspectRatio,
		FileSizeBytes: req.FileSizeBytes,
		TrackStartMs:  req.TrackStartMs,
		TrackEndMs:    req.TrackEndMs,
		IsPublic:      true,
		IsApproved:    false, // requires admin approval
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) GetVideo(ctx context.Context, id uuid.UUID) (*Video, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrVideoNotFound
	}
	return v, nil
}

func (s *Service) DeleteOwnVideo(ctx context.Context, id, userID uuid.UUID) error {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if v == nil {
		return ErrVideoNotFound
	}
	if v.UploaderID != userID {
		return ErrNotOwner
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListByTrack(ctx context.Context, trackID uuid.UUID) ([]Video, error) {
	return s.repo.ListByTrack(ctx, trackID)
}

func (s *Service) ListAll(ctx context.Context, limit, offset int) ([]Video, error) {
	return s.repo.ListAll(ctx, limit, offset)
}

func (s *Service) AdminUpdateVideo(ctx context.Context, id uuid.UUID, req AdminUpdateVideoRequest) (*Video, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrVideoNotFound
	}

	// Apply partial updates
	if req.Title != nil {
		v.Title = *req.Title
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.Type != nil {
		v.Type = VideoType(*req.Type)
	}
	if req.IsPublic != nil {
		v.IsPublic = *req.IsPublic
	}
	if req.IsApproved != nil {
		v.IsApproved = *req.IsApproved
	}
	if req.Status != nil {
		v.Status = VideoStatus(*req.Status)
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// Explore
// ---------------------------------------------------------------------------

func (s *Service) ListExplore(ctx context.Context, limit, offset int) ([]Video, error) {
	return s.repo.ListExplore(ctx, limit, offset)
}

func (s *Service) ListByUser(ctx context.Context, userID, viewerID uuid.UUID, limit, offset int) ([]Video, error) {
	return s.repo.ListByUser(ctx, userID, viewerID, limit, offset)
}

// ---------------------------------------------------------------------------
// Social — Video likes
// ---------------------------------------------------------------------------

func (s *Service) LikeVideo(ctx context.Context, videoID, userID uuid.UUID) error {
	v, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return err
	}
	if v == nil {
		return ErrVideoNotFound
	}

	existing, err := s.repo.GetLike(ctx, videoID, userID)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrAlreadyLiked
	}

	like := &VideoLike{
		VideoID: videoID,
		UserID:  userID,
	}
	if err := s.repo.CreateLike(ctx, like); err != nil {
		return err
	}
	return s.repo.IncrementLikeCount(ctx, videoID)
}

func (s *Service) UnlikeVideo(ctx context.Context, videoID, userID uuid.UUID) error {
	existing, err := s.repo.GetLike(ctx, videoID, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotLiked
	}

	if err := s.repo.DeleteLike(ctx, videoID, userID); err != nil {
		return err
	}
	return s.repo.DecrementLikeCount(ctx, videoID)
}

func (s *Service) ViewVideo(ctx context.Context, videoID uuid.UUID) error {
	v, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return err
	}
	if v == nil {
		return ErrVideoNotFound
	}
	return s.repo.IncrementViewCount(ctx, videoID)
}

// ---------------------------------------------------------------------------
// Video Processing Callback (from Python ML service)
// ---------------------------------------------------------------------------

// HandleVideoCallback processes the result of a video processing job.
// The ML service calls this after audio replacement completes or fails.
func (s *Service) HandleVideoCallback(ctx context.Context, payload VideoCallbackPayload) (*Video, error) {
	videoID, err := uuid.Parse(payload.VideoID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	v, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrVideoNotFound
	}

	switch payload.ProcessingStatus {
	case "completed":
		v.Status = VideoStatusReady
		if payload.FinalVideoPath != "" {
			v.FinalVideoPath = &payload.FinalVideoPath
		}
		if payload.ThumbnailPath != "" {
			v.ThumbnailPath = &payload.ThumbnailPath
		}
		if payload.DurationMs > 0 {
			v.DurationMs = payload.DurationMs
		}
		if payload.AspectRatio != "" {
			v.AspectRatio = payload.AspectRatio
		}

		// Official MVs are auto-approved; user edits require admin approval
		if v.Type == VideoTypeOfficialMV {
			v.IsApproved = true
		}
		// User edits stay IsApproved = false — admin must approve via ApproveVideo endpoint

	case "failed":
		v.Status = VideoStatusFailed
		if payload.ErrorMessage != nil {
			v.Description = "processing failed: " + *payload.ErrorMessage
		}

	default:
		return nil, ErrInvalidInput
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// Track Like Visibility
// ---------------------------------------------------------------------------

func (s *Service) SetTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID, visibility string) (*TrackLikeVisibility, error) {
	if visibility != "public" && visibility != "private" {
		return nil, ErrInvalidInput
	}

	// User must have actually liked the track before setting visibility
	liked, err := s.repo.CheckTrackLikeExists(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}
	if !liked {
		return nil, ErrNotLiked
	}

	if err := s.repo.UpsertTrackLikeVisibility(ctx, userID, trackID, visibility); err != nil {
		return nil, err
	}

	return s.repo.GetTrackLikeVisibility(ctx, userID, trackID)
}

func (s *Service) GetPublicLikedTracks(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TrackLikeVisibility, error) {
	return s.repo.GetPublicLikedTracks(ctx, userID, limit, offset)
}

// ---------------------------------------------------------------------------
// Music Status
// ---------------------------------------------------------------------------

func (s *Service) SetMusicStatus(ctx context.Context, userID uuid.UUID, trackID *string, visibility string) (*UserMusicStatus, error) {
	if visibility != "public" && visibility != "followers" && visibility != "private" {
		return nil, ErrInvalidInput
	}

	var parsedTrackID *uuid.UUID
	if trackID != nil && *trackID != "" {
		pid, err := uuid.Parse(*trackID)
		if err != nil {
			return nil, ErrInvalidInput
		}
		parsedTrackID = &pid
	}

	status := &UserMusicStatus{
		UserID:         userID,
		CurrentTrackID: parsedTrackID,
		Visibility:     visibility,
	}

	if err := s.repo.UpsertMusicStatus(ctx, status); err != nil {
		return nil, err
	}
	return s.repo.GetMusicStatus(ctx, userID)
}

func (s *Service) ClearMusicStatus(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteMusicStatus(ctx, userID)
}

func (s *Service) GetMusicStatus(ctx context.Context, targetUserID uuid.UUID, requestingUserID *uuid.UUID) (*UserMusicStatus, error) {
	status, err := s.repo.GetMusicStatus(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	// No row = not broadcasting (never 404, never forbidden)
	if status == nil {
		return nil, nil
	}

	// TTL staleness check: if updated_at is >30 minutes ago,
	// treat as "not currently playing"
	if time.Since(status.UpdatedAt) > StalenessThreshold {
		return nil, nil
	}

	// Privacy enforcement
	switch status.Visibility {
	case "public":
		// Return to anyone
		return status, nil

	case "followers":
		// Check if requesting user follows the target user
		if requestingUserID == nil {
			// Unauthenticated users cannot see followers-only status
			return nil, nil
		}
		isFollowing, err := s.followService.IsFollowingUser(ctx, *requestingUserID, targetUserID.String())
		if err != nil {
			return nil, err
		}
		if !isFollowing {
			return nil, nil
		}
		return status, nil

	case "private":
		// Only return to the user themselves
		if requestingUserID == nil || *requestingUserID != targetUserID {
			return nil, nil
		}
		return status, nil

	default:
		return nil, nil
	}
}

func (s *Service) GetTrackCoverURL(ctx context.Context, trackID uuid.UUID) (*string, error) {
	return s.repo.GetTrackCoverURL(ctx, trackID)
}

func (s *Service) GetTrackCovers(ctx context.Context, trackIDs []uuid.UUID) (map[uuid.UUID]*string, error) {
	return s.repo.GetTrackCovers(ctx, trackIDs)
}
