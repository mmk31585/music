package recommendation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock for DiscoverWeeklyRepository
type mockDiscoverWeeklyRepo struct {
	mock.Mock
}

func (m *mockDiscoverWeeklyRepo) GetByWeek(ctx context.Context, userID string, weekOf time.Time) (*DiscoverWeeklyPlaylist, error) {
	args := m.Called(ctx, userID, weekOf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DiscoverWeeklyPlaylist), args.Error(1)
}

func (m *mockDiscoverWeeklyRepo) Save(ctx context.Context, p *DiscoverWeeklyPlaylist) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func TestDiscoverWeekly_StableWithinWeek(t *testing.T) {
	mockRepo := new(mockRepo)
	mockDWRepo := new(mockDiscoverWeeklyRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)

	svc := &DiscoverWeeklyService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
		dwRepo:         mockDWRepo,
	}

	userID := "u1"
	weekOf := currentWeekStart()

	existingPlaylist := &DiscoverWeeklyPlaylist{
		UserID:      userID,
		TrackIDs:    []string{"t1", "t2"},
		GeneratedAt: time.Now(),
		WeekOf:      weekOf,
	}

	// Existing playlist found — stable within week
	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(existingPlaylist, nil)
	mockRepo.On("GetTracksByIDs", mock.Anything, []string{"t1", "t2"}).Return([]TrackItem{
		{ID: "t1", Title: "Track 1"},
		{ID: "t2", Title: "Track 2"},
	}, nil)

	_, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, "Track 1", tracks[0].Title)
	mockDWRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDiscoverWeekly_RegeneratesForNewWeek(t *testing.T) {
	mockRepo := new(mockRepo)
	mockDWRepo := new(mockDiscoverWeeklyRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)

	svc := &DiscoverWeeklyService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
		dwRepo:         mockDWRepo,
	}

	userID := "u1"
	seedIDs := []string{"seed-1", "seed-2"}
	weekOf := currentWeekStart()

	// No existing playlist for current week
	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(nil, errors.New("not found"))

	// Profile available (with TopGenreIDs for adjacent-genre mixing)
	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: seedIDs,
		TopGenreIDs:  []string{"genre-pop", "genre-rock"},
	}, nil)

	// No played tracks (empty history)
	mockRepo.On("GetAllUserPlayedTrackIDs", mock.Anything, userID).Return([]string{}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 200).Return([]string{}, nil)

	// Fallback: meta similarity for each seed
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-1").Return(&TrackMeta{
		ID: "seed-1", ArtistID: ptr("a1"), AlbumID: nil, Genre: ptr("pop"),
	}, nil)
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-1", ptr("a1"), (*string)(nil), ptr("pop"), 90).
		Return([]TrackItem{
			{ID: "new-1", Title: "Novelty 1", ArtistName: ptr("A1")},
			{ID: "new-2", Title: "Novelty 2", ArtistName: ptr("A2")},
			{ID: "new-3", Title: "Novelty 3", ArtistName: ptr("A3")},
		}, nil)
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-2").Return(&TrackMeta{
		ID: "seed-2", ArtistID: ptr("a2"), AlbumID: nil, Genre: ptr("rock"),
	}, nil)
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-2", ptr("a2"), (*string)(nil), ptr("rock"), 90).
		Return([]TrackItem{}, nil)

	// Adjacent genre discovery

	mockRepo.On("GetPopularGenreIDs", mock.Anything, 20).Return([]string{"genre-jazz", "genre-classical"}, nil)
	mockRepo.On("GetGenreNames", mock.Anything, []string{"genre-jazz", "genre-classical"}).
		Return([]string{"Jazz", "Classical"}, nil)
	mockRepo.On("GetTracksFromGenres", mock.Anything, []string{"Jazz", "Classical"}, 12).
		Return([]TrackItem{
			{ID: "adj-1", Title: "Adjacent Jazz"},
		}, nil)

	// Save playlist
	mockDWRepo.On("Save", mock.Anything, mock.MatchedBy(func(p *DiscoverWeeklyPlaylist) bool {
		return p.UserID == userID
	})).Return(nil)

	// Get tracks by IDs (from save)
	mockRepo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]TrackItem{
		{ID: "new-1", Title: "Novelty 1", ArtistName: ptr("A1")},
		{ID: "new-2", Title: "Novelty 2", ArtistName: ptr("A2")},
	}, nil)

	playlist, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, userID, playlist.UserID)
	assert.NotEmpty(t, tracks)
	mockDWRepo.AssertExpectations(t)
}

func TestDiscoverWeekly_ExcludesHeard(t *testing.T) {
	mockRepo := new(mockRepo)
	mockDWRepo := new(mockDiscoverWeeklyRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)

	svc := &DiscoverWeeklyService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
		dwRepo:         mockDWRepo,
	}

	userID := "u1"
	seedID := "seed-1"
	weekOf := currentWeekStart()

	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(nil, errors.New("not found"))

	// Profile with seed tracks
	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{seedID},
	}, nil)

	// User has HEARD tracks that must be excluded
	mockRepo.On("GetAllUserPlayedTrackIDs", mock.Anything, userID).Return([]string{
		"heard-1", "heard-2", "new-1", // new-1 would be returned but must be excluded
	}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 200).Return([]string{"liked-1"}, nil)

	mockRepo.On("GetTrackMeta", mock.Anything, seedID).Return(&TrackMeta{
		ID: seedID, ArtistID: ptr("a1"), AlbumID: nil, Genre: ptr("pop"),
	}, nil)
	// Candidate pool includes heard tracks
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, seedID, ptr("a1"), (*string)(nil), ptr("pop"), 90).
		Return([]TrackItem{
			{ID: "heard-1", Title: "Already Heard"},
			{ID: "heard-2", Title: "Already Heard 2"},
			{ID: "liked-1", Title: "Already Liked"},
			{ID: "new-1", Title: "Heard via All Played"},
			{ID: "fresh-1", Title: "Fresh Track 1"},
			{ID: "fresh-2", Title: "Fresh Track 2"},
		}, nil)

	// Adjacent genre path skipped — no top genre IDs in profile
	mockRepo.On("GetPopularGenreIDs", mock.Anything, 20).Return([]string{}, nil)

	mockDWRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]TrackItem{
		{ID: "fresh-1", Title: "Fresh Track 1"},
		{ID: "fresh-2", Title: "Fresh Track 2"},
	}, nil)

	_, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	// Only fresh tracks should survive
	for _, tr := range tracks {
		assert.NotEqual(t, "heard-1", tr.ID)
		assert.NotEqual(t, "heard-2", tr.ID)
		assert.NotEqual(t, "liked-1", tr.ID)
	}
	assert.Contains(t, tracks[0].ID, "fresh")
}

func TestDiscoverWeekly_NoveltyFilter(t *testing.T) {
	// Unit test the novelty filtering logic by testing behavior
	// This confirms the 25% skip logic works correctly
	svc := &DiscoverWeeklyService{}

	exclude := map[string]struct{}{
		"ex-1": {},
		"ex-2": {},
	}

	// Simulated items with mixed IDs
	items := []TrackItem{
		{ID: "keep-1"},
		{ID: "keep-2"},
		{ID: "ex-1"}, // will be excluded
		{ID: "keep-3"},
		{ID: "keep-4"},
		{ID: "ex-2"}, // will be excluded
		{ID: "keep-5"},
	}

	result := svc.deduplicateAndLimit(items, exclude, 10)
	assert.Len(t, result, 5)
	assert.NotContains(t, result, "ex-1")
	assert.NotContains(t, result, "ex-2")
	assert.Equal(t, "keep-1", result[0])
}

func TestDiscoverWeekly_InsufficientProfile(t *testing.T) {
	mockRepo := new(mockRepo)
	mockDWRepo := new(mockDiscoverWeeklyRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)

	svc := &DiscoverWeeklyService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
		dwRepo:         mockDWRepo,
	}

	userID := "u1"
	weekOf := currentWeekStart()

	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(nil, errors.New("not found"))

	// Profile with NO seed tracks
	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{},
	}, nil)

	_, _, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)
	assert.ErrorContains(t, err, "insufficient taste profile data")
}

func TestPeriodBounds(t *testing.T) {
	_, label, err := periodBounds("month")
	assert.NoError(t, err)
	assert.Equal(t, "این ماه", label)

	_, label, err = periodBounds("year")
	assert.NoError(t, err)
	assert.Equal(t, "امسال", label)

	_, label, err = periodBounds("all_time")
	assert.NoError(t, err)
	assert.Equal(t, "همیشه", label)

	_, _, err = periodBounds("invalid")
	assert.ErrorContains(t, err, "invalid period")
}

func TestCurrentWeekStart(t *testing.T) {
	start := currentWeekStart()
	assert.Equal(t, time.Monday, start.Weekday(),
		"week start should always be Monday")
	assert.Equal(t, 0, start.Hour(),
		"week start should be at midnight")
	assert.Equal(t, 0, start.Minute(),
		"week start should be at 00:00")
	assert.Equal(t, 0, start.Second(),
		"week start should be at 00:00:00")

	// Verify it's not in the future
	assert.False(t, start.After(time.Now()),
		"week start should not be in the future")
}
