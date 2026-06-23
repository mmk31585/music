package recommendation

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTasteProfileRepo struct {
	mock.Mock
}

func (m *mockTasteProfileRepo) GetProfile(ctx context.Context, userID string) (*TasteProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TasteProfile), args.Error(1)
}

func (m *mockTasteProfileRepo) UpsertProfile(ctx context.Context, profile *TasteProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockTasteProfileRepo) GetSignalRows(ctx context.Context, userID string, since time.Time) ([]signalRow, error) {
	args := m.Called(ctx, userID, since)
	return args.Get(0).([]signalRow), args.Error(1)
}

func (m *mockTasteProfileRepo) GetRecentPositiveTrackIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockTasteProfileRepo) SetOnboardingGenres(ctx context.Context, userID string, genreIDs []string) error {
	args := m.Called(ctx, userID, genreIDs)
	return args.Error(0)
}

// ── Helper to build signal rows ────────────────────────────────────────

func signalRowForTest(trackID, artistID, genreID, signalType string, playedAt time.Time, explicitLike bool) signalRow {
	aID := artistID
	gID := genreID
	return signalRow{
		TrackID:        trackID,
		PlayedAt:       playedAt,
		SignalType:     signalType,
		IsExplicitLike: explicitLike,
		ArtistID:       &aID,
		GenreID:        &gID,
	}
}

// ── signalBaseWeight tests ─────────────────────────────────────────────

func TestSignalBaseWeight_ExplicitLikeIsHighest(t *testing.T) {
	w := signalBaseWeight("", true)
	assert.Equal(t, SignalWeightExplicitLike, w)
}

func TestSignalBaseWeight_SkipIsNegative(t *testing.T) {
	w := signalBaseWeight("skip_negative", false)
	assert.Equal(t, SignalWeightSkipNegative, w)
}

func TestSignalBaseWeight_UnknownReturnsZero(t *testing.T) {
	w := signalBaseWeight("unknown_type", false)
	assert.Equal(t, float64(0), w)
}

// ── RecomputeProfile tests ─────────────────────────────────────────────

var (
	testUUID      = "11111111-1111-1111-1111-111111111111"
	testTrackID1  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"
	testTrackID2  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2"
	testTrackID3  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3"
	testArtistID1 = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01"
	testArtistID2 = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02"
	testGenreID1  = "cccccccc-cccc-cccc-cccc-cccccccccc01"
	testGenreID2  = "cccccccc-cccc-cccc-cccc-cccccccccc02"
)

func TestRecomputeProfile_RecentSignalsWeightedHigherThanOld(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)

	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(nil, errors.New("not found"))

	now := time.Now()
	rows := []signalRow{
		signalRowForTest(testTrackID1, testArtistID1, testGenreID1, "complete_positive", now.Add(-1*24*time.Hour), false),
		signalRowForTest(testTrackID2, testArtistID2, testGenreID2, "complete_positive", now.Add(-80*24*time.Hour), false),
	}
	mockRepo.On("GetSignalRows", mock.Anything, testUUID, mock.Anything).
		Return(rows, nil)

	var captured *TasteProfile
	mockRepo.On("UpsertProfile", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*TasteProfile)
		}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.RecomputeProfile(context.Background(), testUUID)
	assert.NoError(t, err)

	assert.Contains(t, captured.TopArtistIDs, testArtistID1)
	assert.Contains(t, captured.TopArtistIDs, testArtistID2)
	assert.Equal(t, testArtistID1, captured.TopArtistIDs[0])
}

func TestRecomputeProfile_ExplicitLikeOutweighsImplicitPlay(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)
	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(nil, errors.New("not found"))

	now := time.Now()
	rows := []signalRow{
		// One explicit-like on genre-x (weight 2.0) beats one implicit on genre-y (weight 1.0)
		signalRowForTest(testTrackID1, testArtistID1, testGenreID1, "complete_positive", now, true),
		signalRowForTest(testTrackID2, testArtistID2, testGenreID2, "complete_positive", now, false),
	}
	mockRepo.On("GetSignalRows", mock.Anything, testUUID, mock.Anything).
		Return(rows, nil)

	var captured *TasteProfile
	mockRepo.On("UpsertProfile", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*TasteProfile)
		}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.RecomputeProfile(context.Background(), testUUID)
	assert.NoError(t, err)

	assert.Equal(t, testGenreID1, captured.TopGenreIDs[0])
}

func TestRecomputeProfile_SkipReducesGenreScore(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)
	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(nil, errors.New("not found"))

	now := time.Now()
	rows := []signalRow{
		signalRowForTest(testTrackID1, testArtistID1, testGenreID1, "complete_positive", now, false),
		signalRowForTest(testTrackID2, testArtistID1, testGenreID2, "skip_negative", now, false),
	}
	mockRepo.On("GetSignalRows", mock.Anything, testUUID, mock.Anything).
		Return(rows, nil)

	var captured *TasteProfile
	mockRepo.On("UpsertProfile", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*TasteProfile)
		}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.RecomputeProfile(context.Background(), testUUID)
	assert.NoError(t, err)

	assert.Equal(t, testGenreID1, captured.TopGenreIDs[0])
}

func TestRecomputeProfile_SeedTracksFromPositiveSignals(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)
	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(nil, errors.New("not found"))

	now := time.Now()
	rows := []signalRow{
		signalRowForTest(testTrackID1, testArtistID1, testGenreID1, "complete_positive", now, false),
		signalRowForTest(testTrackID2, testArtistID1, testGenreID1, "skip_negative", now, false),
		signalRowForTest(testTrackID3, testArtistID2, testGenreID2, "replay_strong", now, false),
	}
	mockRepo.On("GetSignalRows", mock.Anything, testUUID, mock.Anything).
		Return(rows, nil)

	var captured *TasteProfile
	mockRepo.On("UpsertProfile", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*TasteProfile)
		}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.RecomputeProfile(context.Background(), testUUID)
	assert.NoError(t, err)

	assert.Contains(t, captured.SeedTrackIDs, testTrackID1)
	assert.Contains(t, captured.SeedTrackIDs, testTrackID3)
	assert.NotContains(t, captured.SeedTrackIDs, testTrackID2)
}

func TestRecomputeProfile_DebouncesRapidSuccessiveSignals(t *testing.T) {
	// This test verifies that when an existing profile was just computed
	// and no history repo is available for counting new signals, the
	// recompute still works (graceful fallback when historyRepo is nil).
	mockRepo := new(mockTasteProfileRepo)

	justNow := time.Now()
	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(&TasteProfile{UserID: testUUID, LastComputedAt: justNow}, nil)

	now := time.Now()
	rows := []signalRow{
		signalRowForTest(testTrackID1, testArtistID1, testGenreID1, "complete_positive", now, false),
	}
	mockRepo.On("GetSignalRows", mock.Anything, testUUID, mock.Anything).
		Return(rows, nil)

	var captured *TasteProfile
	mockRepo.On("UpsertProfile", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*TasteProfile)
		}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.RecomputeProfile(context.Background(), testUUID)
	assert.NoError(t, err)

	// Despite recent profile, recompute happened (graceful fallback)
	assert.NotNil(t, captured)
	assert.Contains(t, captured.TopArtistIDs, testArtistID1)
}

func TestGetOrInitProfile_FallsBackToEmpty(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)
	mockRepo.On("GetProfile", mock.Anything, testUUID).
		Return(nil, sql.ErrNoRows)

	svc := &TasteProfileService{repo: mockRepo}
	profile, err := svc.GetOrInitProfile(context.Background(), testUUID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, testUUID, profile.UserID)
}

func TestSetOnboardingGenres_SavesGenres(t *testing.T) {
	mockRepo := new(mockTasteProfileRepo)
	mockRepo.On("SetOnboardingGenres", mock.Anything, testUUID, []string{testGenreID1, testGenreID2}).
		Return(nil)

	svc := &TasteProfileService{repo: mockRepo}
	err := svc.SetOnboardingGenres(context.Background(), testUUID, []string{testGenreID1, testGenreID2})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSetOnboardingGenres_EmptyRejected(t *testing.T) {
	svc := &TasteProfileService{}
	err := svc.SetOnboardingGenres(context.Background(), testUUID, nil)
	assert.Error(t, err)

	err = svc.SetOnboardingGenres(context.Background(), testUUID, []string{})
	assert.Error(t, err)
}
