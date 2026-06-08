package library

import (
	"context"
	"errors"
)

var ErrInvalidLimit = errors.New("invalid limit")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LikeTrack(ctx context.Context, userID, trackID int64) error {
	exists, err := s.repo.TrackExists(ctx, trackID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFound
	}
	return s.repo.LikeTrack(ctx, userID, trackID)
}

func (s *Service) UnlikeTrack(ctx context.Context, userID, trackID int64) error {
	return s.repo.UnlikeTrack(ctx, userID, trackID)
}

func (s *Service) LikeAlbum(ctx context.Context, userID, albumID int64) error {
	exists, err := s.repo.AlbumExists(ctx, albumID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAlbumNotFound
	}
	return s.repo.LikeAlbum(ctx, userID, albumID)
}

func (s *Service) UnlikeAlbum(ctx context.Context, userID, albumID int64) error {
	return s.repo.UnlikeAlbum(ctx, userID, albumID)
}

func (s *Service) FollowArtist(ctx context.Context, userID, artistID int64) error {
	exists, err := s.repo.ArtistExists(ctx, artistID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrArtistNotFound
	}
	return s.repo.FollowArtist(ctx, userID, artistID)
}

func (s *Service) UnfollowArtist(ctx context.Context, userID, artistID int64) error {
	return s.repo.UnfollowArtist(ctx, userID, artistID)
}

func (s *Service) AddPlayHistory(ctx context.Context, userID, trackID int64) error {
	exists, err := s.repo.TrackExists(ctx, trackID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTrackNotFound
	}
	return s.repo.AddPlayHistory(ctx, userID, trackID)
}

func (s *Service) ListLikedTracks(ctx context.Context, userID int64) ([]LibraryTrackItem, error) {
	return s.repo.ListLikedTracks(ctx, userID)
}

func (s *Service) ListLikedAlbums(ctx context.Context, userID int64) ([]LibraryAlbumItem, error) {
	return s.repo.ListLikedAlbums(ctx, userID)
}

func (s *Service) ListFollowedArtists(ctx context.Context, userID int64) ([]LibraryArtistItem, error) {
	return s.repo.ListFollowedArtists(ctx, userID)
}

func (s *Service) ListPlayHistory(ctx context.Context, userID int64, limit int) ([]LibraryTrackItem, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		return nil, ErrInvalidLimit
	}
	return s.repo.ListPlayHistory(ctx, userID, limit)
}

func (s *Service) ListRecentlyPlayed(ctx context.Context, userID int64, limit int) ([]LibraryTrackItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		return nil, ErrInvalidLimit
	}
	return s.repo.ListRecentlyPlayed(ctx, userID, limit)
}
