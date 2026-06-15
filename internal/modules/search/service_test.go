package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch_EmptyQuery(t *testing.T) {
	svc := NewService(nil)
	resp, err := svc.Search(context.Background(), "", 10)
	assert.NoError(t, err)
	assert.Equal(t, "", resp.Query)
	assert.Empty(t, resp.Tracks)
	assert.Empty(t, resp.Albums)
	assert.Empty(t, resp.Artists)
	assert.Empty(t, resp.Playlists)
}
