package follow

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidArtistID  = errors.New("invalid artist id")
	ErrCannotSelfFollow = errors.New("cannot follow yourself")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FollowUser(ctx context.Context, followerID uuid.UUID, targetUserIDRaw string) error {
	targetUserID, err := uuid.Parse(targetUserIDRaw)
	if err != nil {
		return ErrInvalidUserID
	}

	if followerID == targetUserID {
		return ErrCannotSelfFollow
	}

	return s.repo.FollowUser(ctx, followerID, targetUserID)
}

func (s *Service) UnfollowUser(ctx context.Context, followerID uuid.UUID, targetUserIDRaw string) error {
	targetUserID, err := uuid.Parse(targetUserIDRaw)
	if err != nil {
		return ErrInvalidUserID
	}

	if followerID == targetUserID {
		return ErrCannotSelfFollow
	}

	return s.repo.UnfollowUser(ctx, followerID, targetUserID)
}

func (s *Service) FollowArtist(ctx context.Context, userID uuid.UUID, artistIDRaw string) error {
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return ErrInvalidArtistID
	}

	return s.repo.FollowArtist(ctx, userID, artistID)
}

func (s *Service) UnfollowArtist(ctx context.Context, userID uuid.UUID, artistIDRaw string) error {
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return ErrInvalidArtistID
	}

	return s.repo.UnfollowArtist(ctx, userID, artistID)
}

func (s *Service) IsFollowingUser(ctx context.Context, followerID uuid.UUID, targetUserIDRaw string) (bool, error) {
	targetUserID, err := uuid.Parse(targetUserIDRaw)
	if err != nil {
		return false, ErrInvalidUserID
	}

	return s.repo.IsFollowingUser(ctx, followerID, targetUserID)
}

func (s *Service) IsFollowingArtist(ctx context.Context, userID uuid.UUID, artistIDRaw string) (bool, error) {
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return false, ErrInvalidArtistID
	}

	return s.repo.IsFollowingArtist(ctx, userID, artistID)
}
