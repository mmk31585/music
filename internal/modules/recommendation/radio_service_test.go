package recommendation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRadioRepo struct {
	mock.Mock
}

func (m *mockRadioRepo) CreateSession(ctx context.Context, session *RadioSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *mockRadioRepo) GetSession(ctx context.Context, sessionID string) (*RadioSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RadioSession), args.Error(1)
}

func (m *mockRadioRepo) GetActiveSession(ctx context.Context, userID string) (*RadioSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RadioSession), args.Error(1)
}

func (m *mockRadioRepo) AppendPlayedTrack(ctx context.Context, sessionID string, trackID string) error {
	args := m.Called(ctx, sessionID, trackID)
	return args.Error(0)
}

func (m *mockRadioRepo) UpdateLastActive(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *mockRadioRepo) GetPlayedTrackIDs(ctx context.Context, sessionID string) ([]string, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRadioRepo) GetTracksByIDs(ctx context.Context, trackIDs []string) ([]RadioTrackItem, error) {
	args := m.Called(ctx, trackIDs)
	return args.Get(0).([]RadioTrackItem), args.Error(1)
}

func (m *mockRadioRepo) GetPopularTracks(ctx context.Context, genre string, limit int) ([]RadioTrackItem, error) {
	args := m.Called(ctx, genre, limit)
	return args.Get(0).([]RadioTrackItem), args.Error(1)
}

func TestRadioService_StartRadio_CreatesSessionAndReturnsBatch(t *testing.T) {
	repo := new(mockRadioRepo)
	svc := NewRadioService(repo, nil)

	repo.On("GetActiveSession", mock.Anything, "user-1").Return(nil, nil)
	repo.On("CreateSession", mock.Anything, mock.MatchedBy(func(s *RadioSession) bool {
		return s.UserID == "user-1" && s.SeedTrackID == "track-1"
	})).Return(nil)
	repo.On("AppendPlayedTrack", mock.Anything, mock.Anything, "track-1").Return(nil)
	repo.On("GetPopularTracks", mock.Anything, "", 10).Return([]RadioTrackItem{
		{ID: "t1", Title: "Track 1", ArtistName: "Artist 1"},
		{ID: "t2", Title: "Track 2", ArtistName: "Artist 2"},
	}, nil)
	repo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]RadioTrackItem{
		{ID: "t1", Title: "Track 1", ArtistName: "Artist 1"},
		{ID: "t2", Title: "Track 2", ArtistName: "Artist 2"},
	}, nil)

	session, tracks, err := svc.StartRadio(context.Background(), "user-1", "track-1", "track")

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "user-1", session.UserID)
	assert.Equal(t, "track-1", session.SeedTrackID)
	assert.Len(t, tracks, 2)
	repo.AssertExpectations(t)
}

func TestRadioService_GetNextBatch_ExcludesAlreadyPlayed(t *testing.T) {
	repo := new(mockRadioRepo)
	svc := NewRadioService(repo, nil)

	session := &RadioSession{
		ID:             "session-1",
		UserID:         "user-1",
		SeedTrackID:    "track-1",
		PlayedTrackIDs: []string{"track-1", "t1"},
	}

	repo.On("GetSession", mock.Anything, "session-1").Return(session, nil)
	repo.On("GetPopularTracks", mock.Anything, "", 10).Return([]RadioTrackItem{
		{ID: "t2", Title: "Track 2", ArtistName: "Artist 2"},
		{ID: "t3", Title: "Track 3", ArtistName: "Artist 3"},
		{ID: "t1", Title: "Track 1", ArtistName: "Artist 1"},
	}, nil)
	repo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]RadioTrackItem{
		{ID: "t2", Title: "Track 2", ArtistName: "Artist 2"},
		{ID: "t3", Title: "Track 3", ArtistName: "Artist 3"},
	}, nil)
	repo.On("UpdateLastActive", mock.Anything, "session-1").Return(nil)

	tracks, err := svc.GetNextBatch(context.Background(), "session-1", 10)

	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	for _, tr := range tracks {
		assert.NotEqual(t, "t1", tr.ID, "played track should be excluded")
	}
	repo.AssertExpectations(t)
}

func TestRadioService_GetNextBatch_FallsBackWhenPoolExhausted(t *testing.T) {
	repo := new(mockRadioRepo)
	svc := NewRadioService(repo, nil)

	session := &RadioSession{
		ID:             "session-1",
		UserID:         "user-1",
		SeedTrackID:    "niche-track",
		PlayedTrackIDs: []string{"niche-track"},
	}

	repo.On("GetSession", mock.Anything, "session-1").Return(session, nil)
	repo.On("GetPopularTracks", mock.Anything, "", 10).Return([]RadioTrackItem{
		{ID: "popular-1", Title: "Popular 1", ArtistName: "Artist"},
		{ID: "popular-2", Title: "Popular 2", ArtistName: "Artist"},
	}, nil)
	repo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]RadioTrackItem{
		{ID: "popular-1", Title: "Popular 1", ArtistName: "Artist"},
		{ID: "popular-2", Title: "Popular 2", ArtistName: "Artist"},
	}, nil)
	repo.On("UpdateLastActive", mock.Anything, "session-1").Return(nil)

	tracks, err := svc.GetNextBatch(context.Background(), "session-1", 10)

	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, "popular-1", tracks[0].ID)
	assert.Equal(t, "popular-2", tracks[1].ID)
	repo.AssertExpectations(t)
}

func TestRadioService_EndRadio_UpdatesLastActive(t *testing.T) {
	repo := new(mockRadioRepo)
	svc := NewRadioService(repo, nil)

	repo.On("UpdateLastActive", mock.Anything, "session-1").Return(nil)

	err := svc.EndRadio(context.Background(), "session-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRadioService_StartRadio_EndsExistingSession(t *testing.T) {
	repo := new(mockRadioRepo)
	svc := NewRadioService(repo, nil)

	existingSession := &RadioSession{
		ID:          "old-session",
		UserID:      "user-1",
		SeedTrackID: "old-track",
	}

	repo.On("GetActiveSession", mock.Anything, "user-1").Return(existingSession, nil)
	repo.On("UpdateLastActive", mock.Anything, "old-session").Return(nil)
	repo.On("CreateSession", mock.Anything, mock.MatchedBy(func(s *RadioSession) bool {
		return s.UserID == "user-1" && s.SeedTrackID == "track-2"
	})).Return(nil)
	repo.On("AppendPlayedTrack", mock.Anything, mock.Anything, "track-2").Return(nil)
	repo.On("GetPopularTracks", mock.Anything, "", 10).Return([]RadioTrackItem{}, nil)

	session, tracks, err := svc.StartRadio(context.Background(), "user-1", "track-2", "track")

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "track-2", session.SeedTrackID)
	assert.Empty(t, tracks)
	repo.AssertExpectations(t)
}
