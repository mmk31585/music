package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Enricher struct {
	musicBrainz MusicBrainzClient
	lastFM      LastFMClient
	spotify     SpotifyClient
	lrclib      LRCLibClient
	logger      *zap.Logger
}

func NewEnricher(mb MusicBrainzClient, lfm LastFMClient, spot SpotifyClient, lrc LRCLibClient, logger *zap.Logger) *Enricher {
	return &Enricher{
		musicBrainz: mb,
		lastFM:      lfm,
		spotify:     spot,
		lrclib:      lrc,
		logger:      logger,
	}
}

func (e *Enricher) Enrich(ctx context.Context, title, artist, album string, durationSeconds int) (*EnrichmentResult, error) {
	query := TrackQuery{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: durationSeconds,
	}

	if query.Title == "" && query.Artist == "" {
		return &EnrichmentResult{Attempted: true}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var (
		mbResult   *MusicBrainzResult
		lfmResult  *LastFMResult
		spotResult *SpotifyResult
		lrcResult  *LRCLibResult
	)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		var err error
		mbResult, err = e.musicBrainz.SearchRecording(ctx, query)
		if err != nil {
			e.logger.Warn("musicbrainz enrichment failed", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		lfmResult, err = e.lastFM.SearchTrack(ctx, query)
		if err != nil {
			e.logger.Warn("last.fm enrichment failed", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		spotResult, err = e.spotify.SearchTrack(ctx, query)
		if err != nil {
			e.logger.Warn("spotify enrichment failed", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		lrcResult, err = e.lrclib.SearchLyrics(ctx, query)
		if err != nil {
			e.logger.Warn("lrclib enrichment failed", zap.Error(err))
		}
	}()

	wg.Wait()

	result := MergeResults(mbResult, lfmResult, spotResult, lrcResult)
	return &result, nil
}

func SerializeResult(result *EnrichmentResult) (string, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("marshal enrichment result: %w", err)
	}
	return string(data), nil
}

func DeserializeResult(raw string) (*EnrichmentResult, error) {
	var result EnrichmentResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("unmarshal enrichment result: %w", err)
	}
	if result.Suggestions == nil {
		result.Suggestions = []EnrichedSuggestion{}
	}
	return &result, nil
}
