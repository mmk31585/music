package history

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrInvalidTrackID  = errors.New("invalid track id")
	ErrInvalidDuration = errors.New("invalid duration")
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetHistory(
	ctx context.Context,
	userID uuid.UUID,
	limitRaw string,
	offsetRaw string,
) ([]ListeningHistoryItem, PaginationResponse, error) {
	limit := parsePositiveInt(limitRaw, DefaultLimit)
	offset := parseNonNegativeInt(offsetRaw, 0)

	if limit > MaxLimit {
		limit = MaxLimit
	}

	// Fetch one extra item so we can determine has_more
	items, err := s.repo.ListByUser(ctx, userID, limit+1, offset)
	if err != nil {
		return nil, PaginationResponse{}, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	pagination := PaginationResponse{
		Limit:   limit,
		Offset:  offset,
		Count:   len(items),
		HasMore: hasMore,
	}

	return items, pagination, nil
}

// Record is not exposed through a public route yet.
// This is for playback/player integration later.
func (s *Service) Record(
	ctx context.Context,
	userID uuid.UUID,
	req RecordListeningRequest,
) (*ListeningHistoryItem, error) {
	trackID, err := uuid.Parse(req.TrackID)
	if err != nil {
		return nil, ErrInvalidTrackID
	}

	if req.Duration < 0 {
		return nil, ErrInvalidDuration
	}

	return s.repo.Record(ctx, userID, trackID, req.Duration, req.Completed)
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func parseNonNegativeInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return fallback
	}

	return value
}
