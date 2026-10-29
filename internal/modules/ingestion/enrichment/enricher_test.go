package enrichment

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestEnricher_AllAPIsReturnData(t *testing.T) {
	mb := &mockMusicBrainz{result: fullMusicBrainzResult()}
	lfm := &mockLastFM{result: fullLastFMResult()}
	spot := &mockSpotify{result: fullSpotifyResult()}
	lrc := &mockLRCLib{result: fullLRCLibResult()}

	enricher := NewEnricher(mb, lfm, spot, lrc, zap.NewNop())
	result, err := enricher.Enrich(context.Background(), "Test Song", "Test Artist", "Test Album", 240)
	if err != nil {
		t.Fatalf("Enrich failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.Attempted {
		t.Error("expected Attempted = true")
	}
	if len(result.Suggestions) == 0 {
		t.Fatal("expected at least one suggestion")
	}

	titleFilled := false
	artistMBIDFilled := false
	spotifyIDFilled := false
	for _, s := range result.Suggestions {
		if s.Field == "title" && s.Value == "Test Song" && s.Source == SourceMusicBrainz {
			titleFilled = true
		}
		if s.Field == "artist_mbid" && s.Value == "artist-mb-id-456" {
			artistMBIDFilled = true
		}
		if s.Field == "spotify_id" && s.Value == "spotify-id-abc" {
			spotifyIDFilled = true
		}
	}

	if !titleFilled {
		t.Error("title suggestion missing or incorrect")
	}
	if !artistMBIDFilled {
		t.Error("artist_mbid suggestion missing")
	}
	if !spotifyIDFilled {
		t.Error("spotify_id suggestion missing")
	}
}

func TestEnricher_PartialAPIFailure(t *testing.T) {
	mb := &mockMusicBrainz{result: fullMusicBrainzResult()}
	lfm := &mockLastFM{err: errors.New("rate limited")}
	spot := &mockSpotify{result: fullSpotifyResult()}
	lrc := &mockLRCLib{result: fullLRCLibResult()}

	enricher := NewEnricher(mb, lfm, spot, lrc, zap.NewNop())
	result, err := enricher.Enrich(context.Background(), "Test Song", "Test Artist", "", 0)
	if err != nil {
		t.Fatalf("Enrich failed: %v", err)
	}

	if !result.Attempted {
		t.Error("expected Attempted = true")
	}

	if result.MusicBrainz == nil {
		t.Error("expected MusicBrainz result despite Last.fm failure")
	}
	if result.Spotify == nil {
		t.Error("expected Spotify result despite Last.fm failure")
	}
	if result.LastFM != nil {
		t.Error("expected no Last.fm result on error")
	}
}

func TestEnricher_AllAPIsFail(t *testing.T) {
	mb := &mockMusicBrainz{err: errors.New("timeout")}
	lfm := &mockLastFM{err: errors.New("not found")}
	spot := &mockSpotify{err: errors.New("unauthorized")}
	lrc := &mockLRCLib{result: nil}

	enricher := NewEnricher(mb, lfm, spot, lrc, zap.NewNop())
	result, err := enricher.Enrich(context.Background(), "Test Song", "Test Artist", "", 0)
	if err != nil {
		t.Fatalf("Enrich failed: %v", err)
	}

	if !result.Attempted {
		t.Error("expected Attempted = true")
	}
	if len(result.Suggestions) != 0 {
		t.Error("expected no suggestions when all APIs fail")
	}
}

func TestEnricher_NoMatch(t *testing.T) {
	mb := &mockMusicBrainz{result: nil}
	lfm := &mockLastFM{result: nil}
	spot := &mockSpotify{result: nil}
	lrc := &mockLRCLib{result: nil}

	enricher := NewEnricher(mb, lfm, spot, lrc, zap.NewNop())
	result, err := enricher.Enrich(context.Background(), "Unknown Song", "Unknown Artist", "", 0)
	if err != nil {
		t.Fatalf("Enrich failed: %v", err)
	}

	if result.Attempted != true {
		t.Error("expected Attempted = true")
	}
	if result.MusicBrainz != nil {
		t.Error("expected nil MusicBrainz result on no match")
	}
	if result.LastFM != nil {
		t.Error("expected nil LastFM result on no match")
	}
	if result.Spotify != nil {
		t.Error("expected nil Spotify result on no match")
	}
}

func TestEnricher_EmptyQuery(t *testing.T) {
	mb := &mockMusicBrainz{result: fullMusicBrainzResult()}
	lfm := &mockLastFM{result: fullLastFMResult()}
	spot := &mockSpotify{result: fullSpotifyResult()}
	lrc := &mockLRCLib{result: nil}

	enricher := NewEnricher(mb, lfm, spot, lrc, zap.NewNop())
	result, err := enricher.Enrich(context.Background(), "", "", "", 0)
	if err != nil {
		t.Fatalf("Enrich failed: %v", err)
	}

	if !result.Attempted {
		t.Error("expected Attempted = true")
	}
}

func TestSerializationRoundTrip(t *testing.T) {
	original := &EnrichmentResult{
		MusicBrainz: fullMusicBrainzResult(),
		LastFM:      fullLastFMResult(),
		Spotify:     fullSpotifyResult(),
		Attempted:   true,
		Suggestions: []EnrichedSuggestion{
			{Field: "title", Value: "Test Song", Source: SourceMusicBrainz, Confidence: ConfidenceExact},
		},
	}

	jsonStr, err := SerializeResult(original)
	if err != nil {
		t.Fatalf("SerializeResult failed: %v", err)
	}

	restored, err := DeserializeResult(jsonStr)
	if err != nil {
		t.Fatalf("DeserializeResult failed: %v", err)
	}

	if restored.Attempted != true {
		t.Error("expected Attempted = true after round-trip")
	}
	if restored.MusicBrainz == nil || restored.MusicBrainz.Title != "Test Song" {
		t.Errorf("MusicBrainz title = %v, want 'Test Song'", restored.MusicBrainz)
	}
	if len(restored.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(restored.Suggestions))
	}
}

func TestMergeResults_FullData(t *testing.T) {
	result := MergeResults(fullMusicBrainzResult(), fullLastFMResult(), fullSpotifyResult(), nil)

	fields := map[string]bool{}
	for _, s := range result.Suggestions {
		fields[s.Field] = true
	}

	expected := []string{"title", "artist", "artist_mbid", "album", "album_mbid", "year", "duration_seconds", "genre", "play_count", "spotify_id", "popularity", "preview_url", "album_cover_url", "artist_image_url", "artist_bio", "similar_artists"}
	for _, f := range expected {
		if !fields[f] {
			t.Errorf("missing suggestion field: %s", f)
		}
	}
}

func TestMergeResults_PartialData(t *testing.T) {
	result := MergeResults(nil, fullLastFMResult(), nil, nil)

	if result.MusicBrainz != nil {
		t.Error("expected nil MusicBrainz")
	}
	if result.LastFM == nil {
		t.Error("expected non-nil LastFM")
	}
	if result.Spotify != nil {
		t.Error("expected nil Spotify")
	}

	fields := map[string]bool{}
	for _, s := range result.Suggestions {
		fields[s.Field] = true
	}

	if fields["title"] {
		t.Error("title should not be present without MB")
	}
	if !fields["play_count"] {
		t.Error("play_count should be present from LastFM")
	}
}

func TestMergeResult_Empty(t *testing.T) {
	result := MergeResults(nil, nil, nil, nil)
	if len(result.Suggestions) != 0 {
		t.Error("expected no suggestions for empty merge")
	}
}
