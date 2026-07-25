package reactions

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) React(ctx context.Context, userID, targetID, targetType, reactionType string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	return s.repo.Upsert(ctx, uid, targetID, targetType, reactionType)
}

func (s *Service) RemoveReaction(ctx context.Context, userID, targetID, targetType string) error {
	uid, _ := uuid.Parse(userID)
	return s.repo.Remove(ctx, uid, targetID, targetType)
}

func (s *Service) GetUserReaction(ctx context.Context, userID, targetID, targetType string) (*Reaction, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetUserReaction(ctx, uid, targetID, targetType)
}

func (s *Service) GetCounts(ctx context.Context, targetID, targetType string) (CountsResponse, error) {
	counts, err := s.repo.GetCounts(ctx, targetID, targetType)
	if err != nil {
		return CountsResponse{}, err
	}

	result := CountsResponse{}
	for _, c := range counts {
		switch c.Type {
		case "like":
			result.Like = c.Count
		case "love":
			result.Love = c.Count
		case "dislike":
			result.Dislike = c.Count
		}
		result.Total += c.Count
	}
	return result, nil
}

func (s *Service) GetUserLikedTracks(ctx context.Context, userID string, limit, offset int) ([]Reaction, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetUserReactions(ctx, uid, "track", limit, offset)
}

func (s *Service) GetUserLikedAlbums(ctx context.Context, userID string, limit, offset int) ([]Reaction, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetUserReactions(ctx, uid, "album", limit, offset)
}
