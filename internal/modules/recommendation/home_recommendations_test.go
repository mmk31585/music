package recommendation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPersonalizedTracks_UsesSeedTracksWhenAvailable(t *testing.T) {
	mockRepo := new(mockRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)
	svc := &HomeRecommendationService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
	}

	userID := "u1"
	seedID := "track-seed-1"

	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{seedID},
	}, nil)

	mockRepo.On("GetRecentTracks", mock.Anything, userID, 20).Return([]TrackItem{}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 100).Return([]string{}, nil)
	mockRepo.On("GetTrackMeta", mock.Anything, seedID).Return(&TrackMeta{
		ID: seedID, ArtistID: ptr("artist-1"), AlbumID: nil, Genre: ptr("rock"),
	}, nil)

	expected := []TrackItem{
		{ID: "similar-1", Title: "Similar Track 1", ArtistName: ptr("A1")},
		{ID: "similar-2", Title: "Similar Track 2", ArtistName: ptr("A2")},
	}
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, seedID, ptr("artist-1"), (*string)(nil), ptr("rock"), 10).
		Return(expected, nil)

	items, err := svc.GetPersonalizedTracks(context.Background(), userID, 10)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "Similar Track 1", items[0].Title)
	mockRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
}

func TestGetPersonalizedTracks_FallsBackToOnboardingGenres(t *testing.T) {
	mockRepo := new(mockRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)
	svc := &HomeRecommendationService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
	}

	userID := "u1"
	genreIDs := []string{"genre-jazz-id", "genre-blues-id"}

	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:             userID,
		SeedTrackIDs:       []string{},
		OnboardingGenreIDs: genreIDs,
	}, nil)

	mockRepo.On("GetRecentTracks", mock.Anything, userID, 20).Return([]TrackItem{}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 100).Return([]string{}, nil)
	mockRepo.On("GetGenreNames", mock.Anything, genreIDs).Return([]string{"Jazz", "Blues"}, nil)

	expected := []TrackItem{
		{ID: "g1", Title: "Jazz Track", Genre: ptr("Jazz")},
		{ID: "g2", Title: "Blues Track", Genre: ptr("Blues")},
	}
	mockRepo.On("GetTracksFromGenres", mock.Anything, []string{"Jazz", "Blues"}, 20).Return(expected, nil)

	items, err := svc.GetPersonalizedTracks(context.Background(), userID, 10)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	mockRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
}

func TestGetPersonalizedTracks_ExcludesRecentlyPlayedAndLiked(t *testing.T) {
	mockRepo := new(mockRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)
	svc := &HomeRecommendationService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
	}

	userID := "u1"
	seedID := "track-seed-1"

	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{seedID},
	}, nil)

	mockRepo.On("GetRecentTracks", mock.Anything, userID, 20).Return([]TrackItem{
		{ID: "recent-1"},
	}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 100).Return([]string{
		"liked-1",
	}, nil)

	mockRepo.On("GetTrackMeta", mock.Anything, seedID).Return(&TrackMeta{
		ID: seedID, ArtistID: ptr("artist-1"), AlbumID: nil, Genre: ptr("rock"),
	}, nil)

	similarTracks := []TrackItem{
		{ID: "recent-1", Title: "Recently Played"},
		{ID: "liked-1", Title: "Already Liked"},
		{ID: "fresh-1", Title: "Fresh Track"},
	}
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, seedID, ptr("artist-1"), (*string)(nil), ptr("rock"), 10).
		Return(similarTracks, nil)

	items, err := svc.GetPersonalizedTracks(context.Background(), userID, 10)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Fresh Track", items[0].Title)
	mockRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
}

func TestGetPersonalizedTracks_ExcludesLibrarySaved(t *testing.T) {
	mockRepo := new(mockRepo)
	mockTasteRepo := new(mockTasteProfileRepo)
	profileSvc := NewTasteProfileService(mockTasteRepo, nil, nil)
	svc := &HomeRecommendationService{
		profileService: profileSvc,
		mlClient:       nil,
		repo:           mockRepo,
	}

	userID := "u1"
	seedID := "track-seed-1"

	mockTasteRepo.On("GetProfile", mock.Anything, userID).Return(&TasteProfile{
		UserID:       userID,
		SeedTrackIDs: []string{seedID},
	}, nil)

	mockRepo.On("GetRecentTracks", mock.Anything, userID, 20).Return([]TrackItem{}, nil)
	mockRepo.On("GetLikedTrackIDs", mock.Anything, userID, 100).Return([]string{}, nil)

	mockRepo.On("GetTrackMeta", mock.Anything, seedID).Return(&TrackMeta{
		ID: seedID, ArtistID: ptr("artist-1"), AlbumID: nil, Genre: ptr("rock"),
	}, nil)

	expected := []TrackItem{
		{ID: "fresh-1", Title: "Fresh Track", ArtistName: ptr("A1")},
	}
	mockRepo.On("GetSimilarTracksByMeta", mock.Anything, seedID, ptr("artist-1"), (*string)(nil), ptr("rock"), 10).
		Return(expected, nil)

	items, err := svc.GetPersonalizedTracks(context.Background(), userID, 10)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Fresh Track", items[0].Title)
	mockRepo.AssertExpectations(t)
	mockTasteRepo.AssertExpectations(t)
}
