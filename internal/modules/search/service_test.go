package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch_EmptyQuery(t *testing.T) {
	svc := NewService(nil, nil)
	resp, err := svc.Search(context.Background(), "", "all", 10)
	assert.NoError(t, err)
	assert.Equal(t, "", resp.Query)
	assert.Empty(t, resp.Tracks)
	assert.Empty(t, resp.Albums)
	assert.Empty(t, resp.Artists)
	assert.Empty(t, resp.Playlists)
}

func TestSearchWithUser_EmptyQuery(t *testing.T) {
	svc := NewService(nil, nil)
	resp, err := svc.SearchWithUser(context.Background(), "", "all", 10, "user-1")
	assert.NoError(t, err)
	assert.Empty(t, resp.Tracks)
}

func TestSuggestions_EmptyPrefix(t *testing.T) {
	svc := NewService(nil, nil)
	suggestions, err := svc.Suggestions(context.Background(), "", 5)
	assert.NoError(t, err)
	assert.Nil(t, suggestions)
}

func TestSuggestions_NilOS(t *testing.T) {
	svc := NewService(nil, nil)
	suggestions, err := svc.Suggestions(context.Background(), "tes", 5)
	assert.NoError(t, err)
	assert.Nil(t, suggestions)
}
