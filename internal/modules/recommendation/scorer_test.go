package recommendation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *mockCache) Exists(ctx context.Context, keys ...string) (bool, error) {
	args := m.Called(ctx, keys)
	return args.Bool(0), args.Error(1)
}

func TestScorer_EmptyCandidates(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)
	result := scorer.Score(context.Background(), "user-1", nil, DefaultScoringWeights())
	assert.Nil(t, result)
}

func TestScorer_BasicScoring(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)

	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).
		Return([]UserAffinity{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).
		Return([]string{}, nil)

	duration := 180
	score := 50.0
	candidates := []TrackItem{
		{ID: "t1", Title: "Pop Track", DurationSeconds: &duration, Score: &score, Genre: strPtr("Pop")},
		{ID: "t2", Title: "Unknown Track", DurationSeconds: nil, Score: nil, Genre: nil},
	}

	result := scorer.Score(context.Background(), "user-1", candidates, DefaultScoringWeights())

	assert.Len(t, result, 2)
	assert.Greater(t, result[0].Score, result[1].Score)
	m.AssertExpectations(t)
}

func TestScorer_LikedTrackBonus(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)

	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).
		Return([]UserAffinity{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).
		Return([]string{"t1"}, nil)

	score := 5.0
	candidates := []TrackItem{
		{ID: "t1", Title: "Liked Track", Score: &score},
		{ID: "t2", Title: "Not Liked", Score: &score},
	}

	result := scorer.Score(context.Background(), "user-1", candidates, DefaultScoringWeights())

	assert.Len(t, result, 2)
	assert.Equal(t, "t1", result[0].Track.ID)
	assert.Greater(t, result[0].Score, result[1].Score)
	m.AssertExpectations(t)
}

func TestScorer_AffinityBonus(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)

	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).
		Return([]UserAffinity{
			{TargetID: "t1", Score: 0.8},
		}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).
		Return([]string{}, nil)

	candidates := []TrackItem{
		{ID: "t1", Title: "Affinity Track"},
		{ID: "t2", Title: "No Affinity"},
	}

	result := scorer.Score(context.Background(), "user-1", candidates, DefaultScoringWeights())

	assert.Len(t, result, 2)
	assert.Equal(t, "t1", result[0].Track.ID)
	m.AssertExpectations(t)
}

func TestScorer_DurationScoring(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)

	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).
		Return([]UserAffinity{}, nil)
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).
		Return([]string{}, nil)

	shortDur := 30
	longDur := 360
	score := 10.0
	candidates := []TrackItem{
		{ID: "t1", Title: "Long Track", DurationSeconds: &longDur, Score: &score},
		{ID: "t2", Title: "Short Track", DurationSeconds: &shortDur, Score: &score},
	}

	result := scorer.Score(context.Background(), "user-1", candidates, DefaultScoringWeights())

	assert.Len(t, result, 2)
	assert.Equal(t, "t1", result[0].Track.ID)
	assert.Greater(t, result[0].Score, result[1].Score)
	m.AssertExpectations(t)
}

func TestScorer_DiversityRerank(t *testing.T) {
	scorer := NewScorer(nil, nil)

	tracks := []ScoredTrack{
		{Track: TrackItem{ID: "t1", ArtistID: strPtr("a1")}, Score: 10},
		{Track: TrackItem{ID: "t2", ArtistID: strPtr("a1")}, Score: 9},
		{Track: TrackItem{ID: "t3", ArtistID: strPtr("a2")}, Score: 8},
		{Track: TrackItem{ID: "t4", ArtistID: strPtr("a2")}, Score: 7},
	}

	result := scorer.DiversityRerank(tracks, 1)

	assert.Len(t, result, 2)
	assert.Equal(t, "t1", result[0].Track.ID)
	assert.Equal(t, "t3", result[1].Track.ID)
}

func TestScorer_DiversityRerank_UnderLimit(t *testing.T) {
	scorer := NewScorer(nil, nil)
	tracks := []ScoredTrack{
		{Track: TrackItem{ID: "t1", ArtistID: strPtr("a1")}},
	}
	result := scorer.DiversityRerank(tracks, 1)
	assert.Len(t, result, 1)
}

func TestScorer_DefaultWeights(t *testing.T) {
	w := DefaultScoringWeights()
	assert.Equal(t, 1.5, w.Affinity)
	assert.Equal(t, 2.0, w.Liked)
	assert.Equal(t, 1.0, w.Popularity)
	assert.Equal(t, 0.5, w.Duration)
	assert.Equal(t, 0.3, w.GenreBonus)
}

func TestScorer_AffinityLoadError(t *testing.T) {
	m := new(mockRepo)
	scorer := NewScorer(m, nil)

	m.On("GetUserAffinities", mock.Anything, "user-1", "track", 500).
		Return(nil, errors.New("db error"))
	m.On("GetLikedTrackIDs", mock.Anything, "user-1", 100).
		Return([]string{}, nil)

	score := 3.0
	candidates := []TrackItem{
		{ID: "t1", Title: "Resilient Track", Score: &score},
	}

	result := scorer.Score(context.Background(), "user-1", candidates, DefaultScoringWeights())
	assert.Len(t, result, 1)
	assert.Greater(t, result[0].Score, 0.0)
	m.AssertExpectations(t)
}

func strPtr(s string) *string { return &s }
