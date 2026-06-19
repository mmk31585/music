package social

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type QueueEngine struct {
	repo            Repository
	roomBroadcaster *RoomBroadcaster
	logger          *zap.Logger
}

func NewQueueEngine(repo Repository, roomBroadcaster *RoomBroadcaster, logger *zap.Logger) *QueueEngine {
	return &QueueEngine{
		repo:            repo,
		roomBroadcaster: roomBroadcaster,
		logger:          logger,
	}
}

func (e *QueueEngine) SuggestTrack(ctx context.Context, roomID, trackID, userID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	tid, err := uuid.Parse(trackID)
	if err != nil {
		return fmt.Errorf("invalid track id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	candidate := &QueueCandidate{
		RoomID:      rid,
		TrackID:     tid,
		SuggestedBy: uid,
		VoteCount:   0,
		CreatedAt:   time.Now().UTC(),
	}

	// ON CONFLICT DO NOTHING so duplicate suggests are idempotent
	if err := e.repo.SuggestTrack(ctx, candidate); err != nil {
		return fmt.Errorf("suggest track: %w", err)
	}

	// Auto-cast a vote from the suggester.
	// For a new candidate, candidate.ID was set by the repo.
	// For duplicates (ON CONFLICT DO NOTHING), candidate.ID doesn't exist
	// in the DB, so look up the actual candidate row first.
	if candidate.ID != uuid.Nil {
		if err := e.repo.CastVoteTx(ctx, candidate.ID, uid); err != nil {
			// Candidate may have already existed (ON CONFLICT DO NOTHING) —
			// look up the real candidate and cast the vote there.
			existing, lookupErr := e.repo.GetCandidates(ctx, rid)
			if lookupErr != nil {
				return lookupErr
			}
			for _, c := range existing {
				if c.TrackID == tid {
					candidate = &c
					break
				}
			}
			if candidate.ID != uuid.Nil {
				_ = e.repo.CastVoteTx(ctx, candidate.ID, uid)
			}
		}
	}

	// If nothing is currently playing, auto-advance so this (or the top-voted) track
	// starts playing immediately. This is the expected UX for parties and rooms:
	// suggest a track → it plays.
	np, _ := e.repo.GetNowPlaying(ctx, rid)
	if np == nil {
		e.logger.Info("no track playing — auto-advancing queue after suggest",
			zap.String("room_id", rid.String()))
		if _, err := e.AdvanceQueue(ctx, roomID); err != nil {
			e.logger.Warn("auto-advance after suggest failed",
				zap.String("room_id", rid.String()), zap.Error(err))
		}
	} else {
		e.broadcastQueueState(ctx, rid)
	}

	return nil
}

func (e *QueueEngine) CastVote(ctx context.Context, candidateID, userID string) error {
	cid, err := uuid.Parse(candidateID)
	if err != nil {
		return fmt.Errorf("invalid candidate id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	candidate, err := e.repo.GetCandidateByID(ctx, cid)
	if err != nil {
		return fmt.Errorf("candidate not found: %w", err)
	}

	if err := e.repo.CastVoteTx(ctx, cid, uid); err != nil {
		return fmt.Errorf("cast vote: %w", err)
	}

	e.broadcastQueueState(ctx, candidate.RoomID)
	return nil
}

func (e *QueueEngine) RemoveVote(ctx context.Context, candidateID, userID string) error {
	cid, err := uuid.Parse(candidateID)
	if err != nil {
		return fmt.Errorf("invalid candidate id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	candidate, err := e.repo.GetCandidateByID(ctx, cid)
	if err != nil {
		return fmt.Errorf("candidate not found: %w", err)
	}

	if err := e.repo.RemoveVoteTx(ctx, cid, uid); err != nil {
		return fmt.Errorf("remove vote: %w", err)
	}

	e.broadcastQueueState(ctx, candidate.RoomID)
	return nil
}

func (e *QueueEngine) AdvanceQueue(ctx context.Context, roomIDStr string) (*RoomNowPlaying, error) {
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid room id: %w", err)
	}

	// Advisory lock to prevent concurrent advances for the same room
	if _, err := e.repo.LockRoomQueue(ctx, roomID); err != nil {
		return nil, fmt.Errorf("lock room queue: %w", err)
	}

	candidates, err := e.repo.GetCandidates(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("get candidates: %w", err)
	}

	var np *RoomNowPlaying

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			if candidates[i].VoteCount != candidates[j].VoteCount {
				return candidates[i].VoteCount > candidates[j].VoteCount
			}
			return candidates[i].CreatedAt.Before(candidates[j].CreatedAt)
		})
		winner := candidates[0]
		now := time.Now().UTC()
		np = &RoomNowPlaying{
			RoomID:      roomID,
			TrackID:     winner.TrackID,
			StartedAt:   now,
			SuggestedBy: &winner.SuggestedBy,
			Source:      "vote",
		}
		if err := e.repo.SetNowPlaying(ctx, np); err != nil {
			return nil, fmt.Errorf("set now playing: %w", err)
		}
		if err := e.repo.RemoveCandidate(ctx, winner.ID); err != nil {
			return nil, fmt.Errorf("remove candidate: %w", err)
		}
	} else {
		np, err = e.autoFill(ctx, roomID)
		if err != nil {
			return nil, fmt.Errorf("autofill: %w", err)
		}
	}

	e.broadcastTrackChanged(roomID, np.TrackID, np.Source, np.SuggestedBy)
	e.broadcastQueueState(ctx, roomID)

	return np, nil
}

func (e *QueueEngine) autoFill(ctx context.Context, roomID uuid.UUID) (*RoomNowPlaying, error) {
	trackID, err := e.repo.PickRandomTrack(ctx)
	if err != nil {
		return nil, fmt.Errorf("pick random track: %w", err)
	}

	now := time.Now().UTC()
	np := &RoomNowPlaying{
		RoomID:    roomID,
		TrackID:   *trackID,
		StartedAt: now,
		Source:    "autofill",
	}

	if err := e.repo.SetNowPlaying(ctx, np); err != nil {
		return nil, fmt.Errorf("set now playing: %w", err)
	}

	return np, nil
}

func (e *QueueEngine) GetQueueState(ctx context.Context, roomID, requestingUserID string) (*QueueStateResponse, error) {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid room id: %w", err)
	}
	uid, err := uuid.Parse(requestingUserID)
	if err != nil {
		uid = uuid.Nil
	}

	nowPlaying, _ := e.repo.GetNowPlayingWithTrack(ctx, rid)
	candidates, err := e.repo.GetCandidatesWithTrackAndUser(ctx, rid, uid)
	if err != nil {
		return nil, fmt.Errorf("get candidates: %w", err)
	}
	if candidates == nil {
		candidates = []CandidateResponse{}
	}

	return &QueueStateResponse{
		NowPlaying: nowPlaying,
		Candidates: candidates,
	}, nil
}

func (e *QueueEngine) broadcastTrackChanged(roomID uuid.UUID, trackID uuid.UUID, source string, suggestedBy *uuid.UUID) {
	if e.roomBroadcaster == nil {
		return
	}
	e.roomBroadcaster.TrackChangedDetailed(roomID, trackID, source, suggestedBy)
}

func (e *QueueEngine) broadcastQueueState(ctx context.Context, roomID uuid.UUID) {
	if e.roomBroadcaster == nil {
		return
	}
	state, err := e.GetQueueState(ctx, roomID.String(), "")
	if err != nil {
		e.logger.Warn("broadcast queue state failed", zap.String("room_id", roomID.String()), zap.Error(err))
		return
	}
	e.roomBroadcaster.QueueUpdatedDetailed(roomID, state.Candidates, state.NowPlaying)
}

// TrackEnded is called when a track finishes playing in a room context.
// It validates that the ended track matches what's currently playing, then advances.
func (e *QueueEngine) TrackEnded(ctx context.Context, roomID, trackID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	tid, err := uuid.Parse(trackID)
	if err != nil {
		return fmt.Errorf("invalid track id: %w", err)
	}

	np, err := e.repo.GetNowPlaying(ctx, rid)
	if err != nil {
		return fmt.Errorf("no track currently playing in room: %w", err)
	}

	if np.TrackID != tid {
		return fmt.Errorf("track id mismatch: ended %s != playing %s", tid, np.TrackID)
	}

	_, err = e.AdvanceQueue(ctx, roomID)
	return err
}
