package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRankedQuery_HasQuery(t *testing.T) {
	body := buildRankedQuery(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights())
	_, ok := body["query"]
	assert.True(t, ok)
	_, ok = body["size"]
	assert.True(t, ok)
}

func TestBuildRankedQuery_FunctionScore(t *testing.T) {
	body := buildRankedQuery(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights())
	q := body["query"].(map[string]any)
	fs, ok := q["function_score"]
	assert.True(t, ok)
	fm := fs.(map[string]any)
	assert.Contains(t, fm, "functions")
	assert.Contains(t, fm, "boost_mode")
}

func TestBuildRankedQuery_MinShouldMatch(t *testing.T) {
	body := buildRankedQuery(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights())
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	bq := fs["query"].(map[string]any)
	b := bq["bool"].(map[string]any)
	assert.Equal(t, 1, b["minimum_should_match"])
}

func TestBuildRankedQuery_DefaultWeights(t *testing.T) {
	w := DefaultWeights()
	assert.Equal(t, 0.4, w.Popularity)
	assert.Equal(t, 0.2, w.Freshness)
	assert.Equal(t, 0.2, w.GenreMatch)
	assert.Equal(t, 0.15, w.Collaborative)
	assert.Equal(t, 0.05, w.Personalization)
}

func TestBuildRankedQuery_PersonalizationNil_NothingAdded(t *testing.T) {
	body := buildRankedQueryWithPersonalization(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights(), nil)
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	functions := fs["functions"].([]any)
	assert.Len(t, functions, 2)
}

func TestBuildRankedQuery_PersonalizationEmptyGenres(t *testing.T) {
	body := buildRankedQueryWithPersonalization(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights(), &PersonalizationTerms{})
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	functions := fs["functions"].([]any)
	assert.Len(t, functions, 2)
}

func TestBuildRankedQuery_PersonalizationWithGenres_BoostAdded(t *testing.T) {
	terms := &PersonalizationTerms{
		Genres:    []string{"Pop", "Rock"},
		ArtistIDs: nil,
	}
	body := buildRankedQueryWithPersonalization(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights(), terms)
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	functions := fs["functions"].([]any)
	assert.Len(t, functions, 3)

	genreBoost := functions[2].(map[string]any)
	assert.Equal(t, 0.05, genreBoost["weight"])
	filter := genreBoost["filter"].(map[string]any)
	termsFilter := filter["terms"].(map[string]any)
	assert.Equal(t, []string{"Pop", "Rock"}, termsFilter["genre"])
}

func TestBuildRankedQuery_PersonalizationWithArtists_BoostAdded(t *testing.T) {
	terms := &PersonalizationTerms{
		Genres:    nil,
		ArtistIDs: []string{"artist-1", "artist-2"},
	}
	body := buildRankedQueryWithPersonalization(context.Background(), "test", "test", []string{"title"}, "", DefaultWeights(), terms)
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	functions := fs["functions"].([]any)
	assert.Len(t, functions, 3)

	artistBoost := functions[2].(map[string]any)
	filter := artistBoost["filter"].(map[string]any)
	termsFilter := filter["terms"].(map[string]any)
	assert.Equal(t, []string{"artist-1", "artist-2"}, termsFilter["artist_id"])
}

func TestBuildRankedQuery_PersonalizationZeroWeight_NothingAdded(t *testing.T) {
	w := DefaultWeights()
	w.Personalization = 0
	terms := &PersonalizationTerms{
		Genres:    []string{"Pop"},
		ArtistIDs: []string{"artist-1"},
	}
	body := buildRankedQueryWithPersonalization(context.Background(), "test", "test", []string{"title"}, "", w, terms)
	q := body["query"].(map[string]any)
	fs := q["function_score"].(map[string]any)
	functions := fs["functions"].([]any)
	assert.Len(t, functions, 2)
}

func TestBuildMultiMatchClauses_IncludesFuzzyMatch(t *testing.T) {
	should := buildMultiMatchClauses("test", "test", []string{"title"})
	assert.GreaterOrEqual(t, len(should), 1)

	first := should[0].(map[string]any)
	mm := first["multi_match"].(map[string]any)
	assert.Equal(t, "test", mm["query"])
	assert.Equal(t, "AUTO", mm["fuzziness"])
}

func TestBuildMultiMatchClauses_IncludesPersianFields(t *testing.T) {
	should := buildMultiMatchClauses("test", "test", []string{"title"})
	assert.GreaterOrEqual(t, len(should), 2)

	second := should[1].(map[string]any)
	mm := second["multi_match"].(map[string]any)
	fields := mm["fields"].([]string)
	assert.Contains(t, fields[0], "persian_title")
}

func TestBuildMultiMatchClauses_IncludesAutocomplete(t *testing.T) {
	should := buildMultiMatchClauses("test", "test", []string{"title"})
	assert.GreaterOrEqual(t, len(should), 3)

	third := should[2].(map[string]any)
	mm := third["multi_match"].(map[string]any)
	assert.Equal(t, "bool_prefix", mm["type"])
}
