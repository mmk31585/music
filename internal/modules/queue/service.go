package queue

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidTrackID     = errors.New("invalid track id")
	ErrInvalidQueueItemID = errors.New("invalid queue item id")
	ErrInvalidQueueMode   = errors.New("invalid queue mode")
	ErrInvalidPosition    = errors.New("invalid position")
	ErrEmptyReorderList   = errors.New("reorder list cannot be empty")
)

const (
	ModeNext  = "next"
	ModeLater = "later"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]QueueItem, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) AddTrack(ctx context.Context, userID uuid.UUID, req AddTrackRequest) (*QueueItem, error) {
	trackID, err := uuid.Parse(req.TrackID)
	if err != nil {
		return nil, ErrInvalidTrackID
	}

	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode == "" {
		mode = ModeLater
	}

	switch mode {
	case ModeNext:
		return s.repo.AddNext(ctx, userID, trackID)

	case ModeLater:
		return s.repo.AddLater(ctx, userID, trackID)

	default:
		return nil, ErrInvalidQueueMode
	}
}

func (s *Service) Remove(ctx context.Context, userID uuid.UUID, itemIDRaw string) error {
	itemID, err := uuid.Parse(itemIDRaw)
	if err != nil {
		return ErrInvalidQueueItemID
	}

	return s.repo.Remove(ctx, userID, itemID)
}

func (s *Service) Clear(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Clear(ctx, userID)
}

func (s *Service) Reorder(ctx context.Context, userID uuid.UUID, req ReorderQueueRequest) error {
	if len(req.Items) == 0 {
		return ErrEmptyReorderList
	}

	items := make([]ReorderItem, 0, len(req.Items))

	for _, item := range req.Items {
		id, err := uuid.Parse(item.ID)
		if err != nil {
			return ErrInvalidQueueItemID
		}

		if item.Position < 1 {
			return ErrInvalidPosition
		}

		items = append(items, ReorderItem{
			ID:       id,
			Position: item.Position,
		})
	}

	return s.repo.Reorder(ctx, userID, items)
}

func (s *Service) Move(ctx context.Context, userID uuid.UUID, itemIDRaw string, position int) error {
	itemID, err := uuid.Parse(itemIDRaw)
	if err != nil {
		return ErrInvalidQueueItemID
	}

	if position < 1 {
		return ErrInvalidPosition
	}

	return s.repo.Move(ctx, userID, itemID, position)
}
