package library

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"music/internal/platform/events"
)

var ErrInvalidLimit = errors.New("invalid limit")

type Service struct {
	repo      *Repository
	publisher events.Publisher
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithPublisher(repo *Repository, publisher events.Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) LikeTrack(ctx context.Context, userID string, trackID string) error {
	exists, err := s.repo.TrackExists(ctx, trackID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFound
	}

	if err := s.repo.LikeTrack(ctx, userID, trackID); err != nil {
		return err
	}

	// Auto-add to Liked Songs playlist (best-effort, don't fail the like)
	_ = s.repo.AddTrackToFavorites(ctx, userID, trackID)

	// Publish async event for taste profile signal enrichment
	if s.publisher != nil {
		parsedUserID, _ := uuid.Parse(userID)
		parsedTrackID, _ := uuid.Parse(trackID)
		_ = s.publisher.Publish(ctx, events.TrackLikedEvent{
			BaseEvent: events.NewBaseEvent(),
			UserID:    parsedUserID,
			TrackID:   parsedTrackID,
		})
	}

	return nil
}

func (s *Service) UnlikeTrack(ctx context.Context, userID string, trackID string) error {
	if err := s.repo.UnlikeTrack(ctx, userID, trackID); err != nil {
		return err
	}

	// Auto-remove from Liked Songs playlist (best-effort)
	_ = s.repo.RemoveTrackFromFavorites(ctx, userID, trackID)

	return nil
}

func (s *Service) LikeAlbum(ctx context.Context, userID string, albumID string) error {
	exists, err := s.repo.AlbumExists(ctx, albumID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAlbumNotFound
	}
	return s.repo.LikeAlbum(ctx, userID, albumID)
}

func (s *Service) UnlikeAlbum(ctx context.Context, userID string, albumID string) error {
	return s.repo.UnlikeAlbum(ctx, userID, albumID)
}

func (s *Service) FollowArtist(ctx context.Context, userID string, artistID string) error {
	exists, err := s.repo.ArtistExists(ctx, artistID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrArtistNotFound
	}
	return s.repo.FollowArtist(ctx, userID, artistID)
}

func (s *Service) UnfollowArtist(ctx context.Context, userID string, artistID string) error {
	return s.repo.UnfollowArtist(ctx, userID, artistID)
}

type AddPlayHistoryInput struct {
	UserID    string
	TrackID   string
	Duration  *int
	Completed *bool
}

func (s *Service) AddPlayHistory(ctx context.Context, input AddPlayHistoryInput) error {
	exists, err := s.repo.TrackExists(ctx, input.TrackID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFound
	}
	return s.repo.AddPlayHistory(ctx, input.UserID, input.TrackID, input.Duration, input.Completed)
}

func (s *Service) ListLikedTracks(ctx context.Context, userID string) ([]LibraryTrackItem, error) {
	return s.repo.ListLikedTracks(ctx, userID)
}

func (s *Service) ListLikedAlbums(ctx context.Context, userID string) ([]LibraryAlbumItem, error) {
	return s.repo.ListLikedAlbums(ctx, userID)
}

func (s *Service) ListFollowedArtists(ctx context.Context, userID string) ([]LibraryArtistItem, error) {
	return s.repo.ListFollowedArtists(ctx, userID)
}

func (s *Service) ListPlayHistory(ctx context.Context, userID string, limit int) ([]LibraryTrackItem, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		return nil, ErrInvalidLimit
	}
	return s.repo.ListPlayHistory(ctx, userID, limit)
}

func (s *Service) ListRecentlyPlayed(ctx context.Context, userID string, limit int) ([]LibraryTrackItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		return nil, ErrInvalidLimit
	}
	return s.repo.ListRecentlyPlayed(ctx, userID, limit)
}
