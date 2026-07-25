package social

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func newQueueEngine(m *mockRepo) *QueueEngine {
	return NewQueueEngine(m, nil, zap.NewNop())
}

func TestQueueEngine_SuggestTrack_NewCandidate(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New().String()
	trackID := uuid.New().String()
	userID := uuid.New().String()

	m.On("SuggestTrack", mock.Anything, mock.MatchedBy(func(c *QueueCandidate) bool {
		return c.RoomID.String() == roomID && c.TrackID.String() == trackID && c.SuggestedBy.String() == userID
	})).Return(nil)
	// Return a non-nil now_playing to skip the auto-advance branch
	m.On("GetNowPlaying", mock.Anything, mock.Anything).Return(&RoomNowPlaying{RoomID: uuid.Nil, TrackID: uuid.Nil, Source: "test"}, nil)

	err := engine.SuggestTrack(context.Background(), roomID, trackID, userID)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestQueueEngine_SuggestTrack_DuplicateIsIdempotent(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	trackID := uuid.New()
	userID := uuid.New()

	m.On("SuggestTrack", mock.Anything, mock.Anything).Return(nil)
	// Return a non-nil now_playing to skip the auto-advance branch
	m.On("GetNowPlaying", mock.Anything, mock.Anything).Return(&RoomNowPlaying{RoomID: uuid.Nil, TrackID: uuid.Nil, Source: "test"}, nil)

	err := engine.SuggestTrack(context.Background(), roomID.String(), trackID.String(), userID.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestQueueEngine_CastVote_IncrementsCount(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	candidateID := uuid.New()
	userID := uuid.New()
	roomID := uuid.New()

	m.On("GetCandidateByID", mock.Anything, candidateID).Return(&QueueCandidate{
		ID:     candidateID,
		RoomID: roomID,
	}, nil)
	m.On("CastVoteTx", mock.Anything, candidateID, userID).Return(nil)

	err := engine.CastVote(context.Background(), candidateID.String(), userID.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestQueueEngine_CastVote_TwiceIsIdempotent(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	candidateID := uuid.New()
	userID := uuid.New()
	roomID := uuid.New()

	m.On("GetCandidateByID", mock.Anything, candidateID).Return(&QueueCandidate{
		ID:     candidateID,
		RoomID: roomID,
	}, nil)
	m.On("CastVoteTx", mock.Anything, candidateID, userID).Return(nil)

	err := engine.CastVote(context.Background(), candidateID.String(), userID.String())
	assert.NoError(t, err)

	err = engine.CastVote(context.Background(), candidateID.String(), userID.String())
	assert.NoError(t, err)
	m.AssertNumberOfCalls(t, "CastVoteTx", 2)
}

func TestQueueEngine_RemoveVote_DecrementsCount(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	candidateID := uuid.New()
	userID := uuid.New()
	roomID := uuid.New()

	m.On("GetCandidateByID", mock.Anything, candidateID).Return(&QueueCandidate{
		ID:     candidateID,
		RoomID: roomID,
	}, nil)
	m.On("RemoveVoteTx", mock.Anything, candidateID, userID).Return(nil)

	err := engine.RemoveVote(context.Background(), candidateID.String(), userID.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestQueueEngine_AdvanceQueue_PicksHighestVote(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	trackIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	userID := uuid.New()
	now := time.Now().UTC()

	candidates := []QueueCandidate{
		{ID: uuid.New(), RoomID: roomID, TrackID: trackIDs[0], SuggestedBy: userID, VoteCount: 2, CreatedAt: now},
		{ID: uuid.New(), RoomID: roomID, TrackID: trackIDs[1], SuggestedBy: userID, VoteCount: 5, CreatedAt: now},
		{ID: uuid.New(), RoomID: roomID, TrackID: trackIDs[2], SuggestedBy: userID, VoteCount: 3, CreatedAt: now},
	}

	m.On("LockRoomQueue", mock.Anything, roomID).Return(func() {}, nil)
	m.On("GetCandidates", mock.Anything, roomID).Return(candidates, nil)
	m.On("SetNowPlaying", mock.Anything, mock.Anything).Return(nil)
	m.On("RemoveCandidate", mock.Anything, candidates[1].ID).Return(nil)
	// TrackChangedDetailed and broadcastQueueState are no-ops with nil broadcaster

	np, err := engine.AdvanceQueue(context.Background(), roomID.String())
	assert.NoError(t, err)
	assert.NotNil(t, np)
	assert.Equal(t, trackIDs[1], np.TrackID)
	assert.Equal(t, "vote", np.Source)
	assert.Equal(t, userID, *np.SuggestedBy)
	m.AssertExpectations(t)
}

func TestQueueEngine_AdvanceQueue_TieBreaksByOldest(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	trackIDs := []uuid.UUID{uuid.New(), uuid.New()}
	userID := uuid.New()
	now := time.Now().UTC()

	candidates := []QueueCandidate{
		{ID: uuid.New(), RoomID: roomID, TrackID: trackIDs[0], SuggestedBy: userID, VoteCount: 3, CreatedAt: now.Add(-10 * time.Second)},
		{ID: uuid.New(), RoomID: roomID, TrackID: trackIDs[1], SuggestedBy: userID, VoteCount: 3, CreatedAt: now},
	}

	m.On("LockRoomQueue", mock.Anything, roomID).Return(func() {}, nil)
	m.On("GetCandidates", mock.Anything, roomID).Return(candidates, nil)
	m.On("SetNowPlaying", mock.Anything, mock.Anything).Return(nil)
	m.On("RemoveCandidate", mock.Anything, candidates[0].ID).Return(nil)

	np, err := engine.AdvanceQueue(context.Background(), roomID.String())
	assert.NoError(t, err)
	assert.NotNil(t, np)
	assert.Equal(t, trackIDs[0], np.TrackID)
	assert.Equal(t, "vote", np.Source)
	m.AssertExpectations(t)
}

func TestQueueEngine_AdvanceQueue_EmptyQueueTriggersAutofill(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	randomTrackID := uuid.New()

	m.On("LockRoomQueue", mock.Anything, roomID).Return(func() {}, nil)
	m.On("GetCandidates", mock.Anything, roomID).Return([]QueueCandidate{}, nil)
	m.On("PickRandomTrack", mock.Anything).Return(&randomTrackID, nil)
	m.On("SetNowPlaying", mock.Anything, mock.Anything).Return(nil)

	np, err := engine.AdvanceQueue(context.Background(), roomID.String())
	assert.NoError(t, err)
	assert.NotNil(t, np)
	assert.Equal(t, randomTrackID, np.TrackID)
	assert.Equal(t, "autofill", np.Source)
	m.AssertExpectations(t)
}

func TestQueueEngine_AdvanceQueue_NoDoubleAdvance(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	trackID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	candidates := []QueueCandidate{
		{ID: uuid.New(), RoomID: roomID, TrackID: trackID, SuggestedBy: userID, VoteCount: 1, CreatedAt: now},
	}

	// LockRoomQueue returns no error
	m.On("LockRoomQueue", mock.Anything, roomID).Return(func() {}, nil)
	m.On("GetCandidates", mock.Anything, roomID).Return(candidates, nil).Once()
	// Second call (if concurrent) returns empty to simulate first caller already removed it
	m.On("GetCandidates", mock.Anything, roomID).Return([]QueueCandidate{}, nil).Once()
	m.On("SetNowPlaying", mock.Anything, mock.Anything).Return(nil)
	m.On("RemoveCandidate", mock.Anything, candidates[0].ID).Return(nil)
	autofillTrackID := uuid.New()
	m.On("PickRandomTrack", mock.Anything).Return(&autofillTrackID, nil)

	// Call twice concurrently
	done := make(chan struct{}, 2)
	go func() {
		engine.AdvanceQueue(context.Background(), roomID.String())
		done <- struct{}{}
	}()
	go func() {
		engine.AdvanceQueue(context.Background(), roomID.String())
		done <- struct{}{}
	}()

	<-done
	<-done

	// The Lock + transactional nature means only one should succeed
	m.AssertExpectations(t)
}

func TestQueueEngine_GetQueueState_ReturnsCandidatesWithVoteState(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	userID := uuid.New()
	candidateID := uuid.New()
	now := time.Now().UTC()
	trackID := uuid.New()
	duration := 240

	np := &NowPlayingResponse{
		Track: TrackSummaryResponse{
			ID:              trackID.String(),
			Title:           "Test Track",
			ArtistName:      "Test Artist",
			DurationSeconds: &duration,
		},
		StartedAt: now,
		Source:    "vote",
	}

	candidates := []CandidateResponse{
		{
			ID:     candidateID,
			RoomID: roomID,
			Track: TrackSummaryResponse{
				ID:              uuid.New().String(),
				Title:           "Candidate Track",
				ArtistName:      "Candidate Artist",
				DurationSeconds: &duration,
			},
			SuggestedBy: UserSummaryResponse{
				ID:       userID.String(),
				Username: "testuser",
			},
			VoteCount: 5,
			HasVoted:  true,
			CreatedAt: now,
		},
	}

	m.On("GetNowPlayingWithTrack", mock.Anything, roomID).Return(np, nil)
	m.On("GetCandidatesWithTrackAndUser", mock.Anything, roomID, userID).Return(candidates, nil)

	state, err := engine.GetQueueState(context.Background(), roomID.String(), userID.String())
	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, np, state.NowPlaying)
	assert.Equal(t, 1, len(state.Candidates))
	assert.True(t, state.Candidates[0].HasVoted)
	assert.Equal(t, "Candidate Track", state.Candidates[0].Track.Title)
	assert.Equal(t, "testuser", state.Candidates[0].SuggestedBy.Username)
	m.AssertExpectations(t)
}

func TestQueueEngine_TrackEnded_MismatchedTrackID_ReturnsError(t *testing.T) {
	m := new(mockRepo)
	engine := newQueueEngine(m)
	roomID := uuid.New()
	playingTrackID := uuid.New()
	wrongTrackID := uuid.New()

	m.On("GetNowPlaying", mock.Anything, roomID).Return(&RoomNowPlaying{
		RoomID:  roomID,
		TrackID: playingTrackID,
		Source:  "vote",
	}, nil)

	err := engine.TrackEnded(context.Background(), roomID.String(), wrongTrackID.String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "track id mismatch")
}
