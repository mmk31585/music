package history

// ── Audit of existing history tracking ─────────────────────────────────
// What the existing system captures:
//   1. Event: TrackPlayedEvent (via events.EventTrackPlayed bus) — fires
//      from the player module when a track finishes playing.
//   2. Fields stored in listening_history:
//      - user_id, track_id, played_at (timestamp)
//      - duration (seconds played, INTEGER, min=0)
//      - completed (BOOLEAN, FALSE by default — marks natural end vs abort)
//   3. Does NOT capture:
//      - session_id (no way to group plays in a listening session)
//      - completion_percent (only raw duration in seconds)
//      - signal_type (no skip/completion/replay classification)
//      - track_duration_ms (no denominator to compute completion % from)
//   4. The Record endpoint is not exposed publicly yet — only consumed via
//      the event handler from player events.
//   5. There is a separate play_history table in the library module that
//      tracks play counts via an explicit REST endpoint. The taste profile
//      system uses listening_history as its primary source.
// What this Phase adds:
//   - session_id, track_duration_ms, completion_percent, signal_type
//     columns to listening_history
//   - Signal classification (skip_negative, complete_positive,
//     replay_strong, partial_neutral) at write time
//   - Explicit like (is_explicit_like) integration
//   - Taste profile aggregation in the recommendation module,
//     driven asynchronously by a "track.signal_recorded" event
// ────────────────────────────────────────────────────────────────────────

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"music/internal/platform/events"
)

var (
	ErrInvalidTrackID  = errors.New("invalid track id")
	ErrInvalidDuration = errors.New("invalid duration")
	ErrInvalidSignal   = errors.New("invalid signal type")
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

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
// It now computes signal type and publishes an async event.
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

	// Basic insert (keep old path for backward compat)
	item, err := s.repo.Record(ctx, userID, trackID, req.Duration, req.Completed)
	if err != nil {
		return nil, err
	}

	// Compute and store signal (non-blocking, async)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("history recording panicked", zap.Any("recover", r))
			}
		}()

		sigCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		playedMs := int64(req.Duration) * 1000
		trackMs := req.TrackDurationMs
		if trackMs <= 0 {
			trackMs = playedMs // fallback if unknown
		}

		completion := float64(0)
		if trackMs > 0 {
			completion = float64(playedMs) / float64(trackMs)
		}

		var sessionUUID uuid.UUID
		if req.SessionID != "" {
			sessionUUID, _ = uuid.Parse(req.SessionID)
		}

		// Determine if this is a replay within the session
		isReplay := false
		if sessionUUID != uuid.Nil {
			signals, err := s.repo.ListSignalsByUser(sigCtx, userID, 1)
			if err == nil {
				for _, sig := range signals {
					if sig.TrackID == trackID && sig.SessionID == sessionUUID {
						isReplay = true
						break
					}
				}
			}
		}

		signalType := ClassifySignal(playedMs, trackMs, isReplay)

		_ = s.repo.RecordSignal(sigCtx, userID, trackID, playedMs, trackMs, completion, signalType, false, sessionUUID)

		// Publish event for async profile recomputation
		if s.publisher != nil {
			_ = s.publisher.Publish(sigCtx, events.PlaybackSignalRecordedEvent{
				BaseEvent:         events.NewBaseEvent(),
				UserID:            userID,
				TrackID:           trackID,
				SessionID:         sessionUUID,
				PlayedDurationMs:  playedMs,
				TrackDurationMs:   trackMs,
				CompletionPercent: completion,
				SignalType:        signalType,
			})
		}
	}()

	return item, nil
}

// RecordSignal is a synchronous version that records a classified signal
// and publishes the event. Used by the like integration and tests.
func (s *Service) RecordSignal(
	ctx context.Context,
	userID, trackID uuid.UUID,
	playedDurationMs, trackDurationMs int64,
	completionPercent float64,
	signalType string,
	isExplicitLike bool,
	sessionID uuid.UUID,
) error {
	if err := s.repo.RecordSignal(ctx, userID, trackID, playedDurationMs, trackDurationMs, completionPercent, signalType, isExplicitLike, sessionID); err != nil {
		return err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, events.PlaybackSignalRecordedEvent{
			BaseEvent:         events.NewBaseEvent(),
			UserID:            userID,
			TrackID:           trackID,
			SessionID:         sessionID,
			PlayedDurationMs:  playedDurationMs,
			TrackDurationMs:   trackDurationMs,
			CompletionPercent: completionPercent,
			SignalType:        signalType,
		})
	}

	return nil
}

// ListSignalsByUser returns signal-enriched history for taste profile recomputation.
func (s *Service) ListSignalsByUser(ctx context.Context, userID uuid.UUID, days int) ([]PlaybackSignal, error) {
	return s.repo.ListSignalsByUser(ctx, userID, days)
}

// DeleteItem removes a single history entry by ID.
func (s *Service) DeleteItem(ctx context.Context, userID, historyID uuid.UUID) error {
	return s.repo.DeleteByID(ctx, userID, historyID)
}

// ClearAll removes all listening history for the user.
func (s *Service) ClearAll(ctx context.Context, userID uuid.UUID) error {
	return s.repo.ClearAll(ctx, userID)
}

// GetRepo exposes the repository for cross-module usage (taste profile).
func (s *Service) GetRepo() *Repository {
	return s.repo
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
