package recommendation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetPopularTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetBestTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetRecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetTrackMeta(ctx context.Context, trackID string) (*TrackMeta, error) {
	args := m.Called(ctx, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TrackMeta), args.Error(1)
}

func (m *mockRepo) GetSimilarTracksByMeta(ctx context.Context, trackID string, artistID, albumID, genre *string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, trackID, artistID, albumID, genre, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetTracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, artistID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetTracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, genre, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetFollowedArtistIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetTopArtistIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetTopGenres(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetLikedTrackIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetTracksFromArtists(ctx context.Context, artistIDs []string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, artistIDs, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetTracksFromGenres(ctx context.Context, genres []string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, genres, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetTracksByIDs(ctx context.Context, trackIDs []string) ([]TrackItem, error) {
	args := m.Called(ctx, trackIDs)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetPopularTrackIDs(ctx context.Context, limit int) ([]string, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetCoListenTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetRecentlyPlayedIDs(ctx context.Context, userID string, hours int) ([]string, error) {
	args := m.Called(ctx, userID, hours)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) GetSimilarTracksByCooccurrence(ctx context.Context, trackID string, limit int) ([]TrackItem, error) {
	args := m.Called(ctx, trackID, limit)
	return args.Get(0).([]TrackItem), args.Error(1)
}

func (m *mockRepo) GetUserAffinities(ctx context.Context, userID string, targetType string, limit int) ([]UserAffinity, error) {
	args := m.Called(ctx, userID, targetType, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]UserAffinity), args.Error(1)
}

func ptr(s string) *string { return &s }

func TestPopularTracks_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	expected := []TrackItem{
		{ID: "t1", Title: "Track 1", ArtistName: ptr("A1")},
		{ID: "t2", Title: "Track 2", ArtistName: ptr("A2")},
	}
	m.On("GetPopularTracks", mock.Anything, 10).Return(expected, nil)

	items, err := svc.PopularTracks(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "Track 1", items[0].Title)
	m.AssertExpectations(t)
}

func TestPopularTracks_DefaultLimit(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPopularTracks", mock.Anything, 20).Return([]TrackItem{}, nil)

	items, err := svc.PopularTracks(context.Background(), 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestPopularTracks_InvalidLimit(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	_, err := svc.PopularTracks(context.Background(), 200)
	assert.ErrorIs(t, err, ErrInvalidLimit)
}

func TestSimilarTracks_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	artistID := "artist-1"
	genre := "Pop"
	m.On("GetTrackMeta", mock.Anything, "t1").Return(&TrackMeta{
		ID: "t1", ArtistID: &artistID, AlbumID: nil, Genre: &genre,
	}, nil)
	m.On("GetSimilarTracksByMeta", mock.Anything, "t1", &artistID, (*string)(nil), &genre, 5).
		Return([]TrackItem{
			{ID: "t2", Title: "Similar 1"},
			{ID: "t3", Title: "Similar 2"},
		}, nil)

	items, err := svc.SimilarTracks(context.Background(), "t1", 5)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	m.AssertExpectations(t)
}

func TestSimilarTracks_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetTrackMeta", mock.Anything, "missing").Return(nil, ErrTrackNotFound)

	_, err := svc.SimilarTracks(context.Background(), "missing", 5)
	assert.ErrorIs(t, err, ErrTrackNotFound)
	m.AssertExpectations(t)
}

func TestTracksByArtist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetTracksByArtist", mock.Anything, "artist-1", 10).Return([]TrackItem{
		{ID: "t1", Title: "A-Track"},
	}, nil)

	items, err := svc.TracksByArtist(context.Background(), "artist-1", 10)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestTracksByGenre_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetTracksByGenre", mock.Anything, "Pop", 10).Return([]TrackItem{
		{ID: "t1", Title: "Pop Track"},
	}, nil)

	items, err := svc.TracksByGenre(context.Background(), "Pop", 10)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestForYou_Basic(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	coListen := []TrackItem{
		{ID: "t1", Title: "Co-Listen 1", ArtistID: ptr("a1")},
		{ID: "t2", Title: "Co-Listen 2", ArtistID: ptr("a2")},
	}
	m.On("GetCoListenTracks", mock.Anything, "user-1", 100).Return(coListen, nil)
	m.On("GetFollowedArtistIDs", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetTopArtistIDs", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetTopGenres", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).Return([]string{}, nil)
	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).Return([]UserAffinity{}, nil)
	m.On("GetRecentlyPlayedIDs", mock.Anything, "user-1", 24).Return([]string{}, nil)
	m.On("GetPopularTrackIDs", mock.Anything, 24).Return([]string{}, nil)

	items, err := svc.ForYou(context.Background(), "user-1", 10)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	m.AssertExpectations(t)
}

func TestForYou_WithExclusions(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	coListen := []TrackItem{
		{ID: "t1", Title: "Liked Track", ArtistID: ptr("a1")},
		{ID: "t2", Title: "Recent Track", ArtistID: ptr("a2")},
		{ID: "t3", Title: "Fresh Track", ArtistID: ptr("a3")},
	}
	m.On("GetCoListenTracks", mock.Anything, "user-1", 100).Return(coListen, nil)
	m.On("GetFollowedArtistIDs", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetTopArtistIDs", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetTopGenres", mock.Anything, "user-1", 10).Return([]string{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 10).Return([]string{"t1"}, nil)
	m.On("GetTrackMeta", mock.Anything, "t1").Return(&TrackMeta{
		ID: "t1", ArtistID: ptr("a1"), AlbumID: nil, Genre: ptr("Pop"),
	}, nil)
	m.On("GetSimilarTracksByMeta", mock.Anything, "t1", ptr("a1"), (*string)(nil), ptr("Pop"), 20).
		Return([]TrackItem{}, nil)
	m.On("GetSimilarTracksByCooccurrence", mock.Anything, "t1", 10).
		Return([]TrackItem{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).Return([]string{"t1"}, nil)
	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).Return([]UserAffinity{}, nil)
	m.On("GetRecentlyPlayedIDs", mock.Anything, "user-1", 24).Return([]string{"t2"}, nil)
	m.On("GetPopularTrackIDs", mock.Anything, mock.Anything).Return([]string{}, nil)

	items, err := svc.ForYou(context.Background(), "user-1", 10)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "t3", items[0].ID)
	m.AssertExpectations(t)
}

func TestForYou_RepoError(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetCoListenTracks", mock.Anything, "user-1", 100).Return(nil, errors.New("db down"))

	_, err := svc.ForYou(context.Background(), "user-1", 10)
	assert.ErrorContains(t, err, "db down")
}

func TestBestTracks_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetBestTracks", mock.Anything, 10).Return([]TrackItem{
		{ID: "t1", Title: "Best 1"},
	}, nil)

	items, err := svc.BestTracks(context.Background(), 10)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestRecentTracks_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetRecentTracks", mock.Anything, "user-1", 10).Return([]TrackItem{
		{ID: "t1", Title: "Recent 1"},
	}, nil)

	items, err := svc.RecentTracks(context.Background(), "user-1", 10)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}
