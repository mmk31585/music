package recommendation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errNotFound = errors.New("not found")

// mockService implements the Service interface for testing HomeFeedService
type mockService struct {
	mock.Mock
}

func (m *mockService) PopularTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) BestTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) RecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) SimilarTracks(ctx context.Context, trackID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, trackID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) TracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, artistID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) TracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, genre, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockService) ForYou(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

// =============================================================================
// Real-User Scenario Tests
// Simulates the full flow a user experiences with Phases 5 & 6:
//  1. User visits Home → sees 6 personalized sections
//  2. User opens Discover Weekly → gets a weekly playlist
//  3. User revisits Discover Weekly same week → stable playlist
//  4. Discover Weekly endpoint response format matches frontend expectations
// =============================================================================

// --- Scenario 1: User visits Home — HomeFeed returns 6 sections ---

func TestRealUser_HomeFeedReturnsAllSections(t *testing.T) {
	mockRepo := new(mockRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)

	// Mock the Service interface for ForYou
	mockSvc := new(mockService)
	svc := &HomeFeedService{
		repo:       mockRepo,
		profileSvc: profileSvc,
		mlClient:   nil,
		svc:        mockSvc,
	}

	userID := "user-real-1"

	// 1. Recently played
	mockRepo.On("GetRecentTracks", mock.Anything, userID, 10).
		Return([]TrackItem{
			{ID: "recent-1", Title: "Recent Track", ArtistName: ptr("A1")},
		}, nil)

	// 2. Trending
	mockRepo.On("GetPopularTracks", mock.Anything, 10).
		Return([]TrackItem{
			{ID: "popular-1", Title: "Popular Track", ArtistName: ptr("A2")},
		}, nil)

	// 3. For You
	mockSvc.On("ForYou", mock.Anything, userID, 10).
		Return([]TrackItem{
			{ID: "foryou-1", Title: "Personalized Track", ArtistName: ptr("A3")},
		}, nil)

	// 4. Taste profile (for artists + seeds + genres)
	mockTasteRepo.On("GetProfile", mock.Anything, userID).
		Return(&TasteProfile{
			UserID:       userID,
			TopArtistIDs: []string{"a1"},
			SeedTrackIDs: []string{"seed-1"},
			TopGenreIDs:  []string{"g1"},
		}, nil)

	// 5. From your favorite artists
	mockRepo.On("GetTracksFromArtists", mock.Anything, []string{"a1"}, 10).
		Return([]TrackItem{
			{ID: "artist-1", Title: "Artist Track", ArtistName: ptr("A1")},
		}, nil)

	// 6. Because you listened to X (seed track)
	mockRepo.On("GetTracksByIDs", mock.Anything, []string{"seed-1"}).
		Return([]TrackItem{
			{ID: "seed-1", Title: "Seed Track", ArtistName: ptr("A1")},
		}, nil)
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-1").
		Return(&TrackMeta{ID: "seed-1", ArtistID: ptr("a1"), AlbumID: nil, Genre: ptr("Pop")}, nil)
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-1", ptr("a1"), (*string)(nil), ptr("Pop"), 10).
		Return([]TrackItem{
			{ID: "similar-1", Title: "Similar Track", ArtistName: ptr("A3")},
		}, nil)

	// 7. Popular in your genres
	mockRepo.On("GetTopGenres", mock.Anything, userID, 5).
		Return([]string{"Pop"}, nil)
	mockRepo.On("GetTracksFromGenres", mock.Anything, []string{"Pop"}, 10).
		Return([]TrackItem{
			{ID: "genre-1", Title: "Genre Track", ArtistName: ptr("A4")},
		}, nil)

	feed, err := svc.BuildFeed(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	assert.Len(t, feed.Sections, 6, "Should have 6 sections")

	// Verify each section has items and correct IDs
	sectionIDs := make([]string, len(feed.Sections))
	for i, s := range feed.Sections {
		sectionIDs[i] = s.ID
		assert.Greater(t, len(s.Items), 0, "Section %s should have items", s.ID)
	}
	assert.Contains(t, sectionIDs, "recently_played")
	assert.Contains(t, sectionIDs, "trending")
	assert.Contains(t, sectionIDs, "for_you")
	assert.Contains(t, sectionIDs, "from_your_artists")
	assert.Contains(t, sectionIDs, "your_genres")

	// The "because_of" section ID contains a truncated seed ID
	var becauseSectionID string
	for _, id := range sectionIDs {
		if len(id) > 8 && id[:8] == "because_" {
			becauseSectionID = id
			break
		}
	}
	assert.Contains(t, becauseSectionID, "because_of_")

	mockRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
	mockSvc.AssertExpectations(t)
}

// --- Scenario 2: User opens Discover Weekly — full generation flow ---

func TestRealUser_DiscoverWeeklyGeneration(t *testing.T) {
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

	userID := "user-real-1"
	weekOf := currentWeekStart()

	// First visit: no existing playlist
	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(nil, errNotFound)

	// Profile: user has seed tracks + genres for adjacent mixing
	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{"seed-1", "seed-2"},
		TopGenreIDs:  []string{"g-pop", "g-rock"},
	}, nil)

	// User has heard some tracks (must be excluded from results)
	mockRepo.On("GetAllUserPlayedTrackIDs", mock.Anything, userID).
		Return([]string{"heard-1", "heard-2"}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 200).
		Return([]string{}, nil)

	// Fallback: similarToAnySeed iterates seeds, skipping top 25%
	// With 2 seeds, skipSeeds = 2/4 = 0, so both seeds are used
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-1").
		Return(&TrackMeta{ID: "seed-1", ArtistID: ptr("a1"), Genre: ptr("Pop")}, nil)
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-2").
		Return(&TrackMeta{ID: "seed-2", ArtistID: ptr("a2"), Genre: ptr("Rock")}, nil)

	candidateSize := 90
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-1", ptr("a1"), (*string)(nil), ptr("Pop"), candidateSize).
		Return([]TrackItem{
			{ID: "cand-1", Title: "Candidate 1"},
			{ID: "cand-2", Title: "Candidate 2"},
			{ID: "heard-1", Title: "Already Heard"}, // must be excluded
		}, nil)
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-2", ptr("a2"), (*string)(nil), ptr("Rock"), candidateSize).
		Return([]TrackItem{}, nil)

	// Adjacent genres: find genres NOT in the user's top genres
	mockRepo.On("GetPopularGenreIDs", mock.Anything, 20).
		Return([]string{"g-jazz", "g-classical"}, nil)
	mockRepo.On("GetGenreNames", mock.Anything, []string{"g-jazz", "g-classical"}).
		Return([]string{"Jazz", "Classical"}, nil)
	mockRepo.On("GetTracksFromGenres", mock.Anything, []string{"Jazz", "Classical"}, 12).
		Return([]TrackItem{
			{ID: "adj-1", Title: "Jazz Track"},
		}, nil)

	mockDWRepo.On("Save", mock.Anything, mock.MatchedBy(func(p *DiscoverWeeklyPlaylist) bool {
		return p.UserID == userID && len(p.TrackIDs) > 0
	})).Return(nil)

	mockRepo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]TrackItem{
		{ID: "cand-1", Title: "Candidate 1"},
		{ID: "cand-2", Title: "Candidate 2"},
	}, nil)

	// ACT: Generate Discover Weekly
	playlist, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, userID, playlist.UserID)
	assert.Equal(t, weekOf, playlist.WeekOf)

	// Verify heard tracks excluded from results
	for _, tr := range tracks {
		assert.NotEqual(t, "heard-1", tr.ID)
		assert.NotEqual(t, "heard-2", tr.ID)
	}
	assert.Greater(t, len(tracks), 0, "Should have at least some tracks")

	mockDWRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// --- Scenario 3: User revisits Discover Weekly same week → stable ---

func TestRealUser_DiscoverWeeklyStableWithinWeek(t *testing.T) {
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

	userID := "user-real-1"
	weekOf := currentWeekStart()

	// Simulate: user already has a playlist from this week
	existingPlaylist := &DiscoverWeeklyPlaylist{
		UserID:   userID,
		TrackIDs: []string{"dw-1", "dw-2", "dw-3"},
		WeekOf:   weekOf,
	}

	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(existingPlaylist, nil)
	mockRepo.On("GetTracksByIDs", mock.Anything, []string{"dw-1", "dw-2", "dw-3"}).
		Return([]TrackItem{
			{ID: "dw-1", Title: "Weekly Pick 1"},
			{ID: "dw-2", Title: "Weekly Pick 2"},
			{ID: "dw-3", Title: "Weekly Pick 3"},
		}, nil)

	// ACT: Get playground for the same week
	playlist, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, existingPlaylist.TrackIDs, playlist.TrackIDs,
		"Same week should return identical playlist")
	assert.Len(t, tracks, 3, "Should have the same 3 tracks")
	assert.Equal(t, "Weekly Pick 1", tracks[0].Title)

	mockDWRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// --- Scenario 4: Discover Weekly response format matches frontend ---

func TestRealUser_DiscoverWeeklyResponseShape(t *testing.T) {
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

	userID := "user-real-1"
	weekOf := currentWeekStart()

	mockDWRepo.On("GetByWeek", mock.Anything, userID, weekOf).Return(nil, errNotFound)
	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{"seed-1"},
		TopGenreIDs:  []string{},
	}, nil)
	mockRepo.On("GetAllUserPlayedTrackIDs", mock.Anything, userID).Return([]string{}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 200).Return([]string{}, nil)
	mockRepo.On("GetTrackMeta", mock.Anything, "seed-1").
		Return(&TrackMeta{ID: "seed-1", ArtistID: ptr("a1"), Genre: ptr("Pop")}, nil)
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, "seed-1", ptr("a1"), (*string)(nil), ptr("Pop"), 90).
		Return([]TrackItem{
			{ID: "dw-1", Title: "Weekly Pick 1", ArtistName: ptr("A1")},
			{ID: "dw-2", Title: "Weekly Pick 2", ArtistName: ptr("A2")},
		}, nil)
	mockRepo.On("GetPopularGenreIDs", mock.Anything, 20).Return([]string{}, nil)
	mockDWRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetTracksByIDs", mock.Anything, mock.Anything).Return([]TrackItem{
		{ID: "dw-1", Title: "Weekly Pick 1", ArtistName: ptr("A1")},
		{ID: "dw-2", Title: "Weekly Pick 2", ArtistName: ptr("A2")},
	}, nil)

	playlist, tracks, err := svc.GetOrGenerateWeeklyPlaylist(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.NotNil(t, tracks)

	// The handler (handler.go:386-394) wraps this as:
	// { "success": true, "data": { "playlist": { "generated_at", "week_of", "track_count" }, "tracks": [...] } }
	// Verify the inner fields that the frontend schema parses
	assert.NotZero(t, playlist.GeneratedAt, "generated_at must be set")
	assert.Equal(t, weekOf, playlist.WeekOf, "week_of must match current week")

	// Verify tracks have the fields the frontend RecommmendationTrack type expects
	for _, tr := range tracks {
		assert.NotEmpty(t, tr.ID, "Track must have an ID")
		assert.NotEmpty(t, tr.Title, "Track must have a title")
	}
	assert.Len(t, tracks, 2, "Should have exactly 2 tracks")
}
