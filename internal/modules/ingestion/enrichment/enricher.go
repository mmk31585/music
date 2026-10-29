package enrichment

import (
	"context"
	"encoding/json"
	"errors"
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

func (e *Enricher) Enrich(ctx context.Context, title, artist, album string, durationSeconds int, scope ...EnrichScope) (*EnrichmentResult, error) {
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

	// Determine which sources to run
	runSource := func(s Source) bool {
		if len(scope) == 0 || len(scope[0]) == 0 {
			return true // no scope = run all
		}
		for _, allowed := range scope[0] {
			if allowed == s {
				return true
			}
		}
		return false
	}

	var (
		mbResult   *MusicBrainzResult
		lfmResult  *LastFMResult
		spotResult *SpotifyResult
		lrcResult  *LRCLibResult
	)

	var wg sync.WaitGroup

	if runSource(SourceMusicBrainz) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			mbResult, err = e.musicBrainz.SearchRecording(ctx, query)
			if err != nil {
				e.logger.Warn("musicbrainz enrichment failed", zap.Error(err))
			}
		}()
	}

	if runSource(SourceLastFM) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			lfmResult, err = e.lastFM.SearchTrack(ctx, query)
			if err != nil {
				if errors.Is(err, ErrLastFMNotConfigured) {
					e.logger.Debug("last.fm enrichment skipped: not configured")
				} else {
					e.logger.Warn("last.fm enrichment failed", zap.Error(err))
				}
			}
		}()
	}

	if runSource(SourceSpotify) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			spotResult, err = e.spotify.SearchTrack(ctx, query)
			if err != nil {
				if errors.Is(err, ErrSpotifyNotConfigured) {
					e.logger.Debug("spotify enrichment skipped: not configured")
				} else {
					e.logger.Warn("spotify enrichment failed", zap.Error(err))
				}
			}
		}()
	}

	if runSource(SourceLRCLib) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			lrcResult, err = e.lrclib.SearchLyrics(ctx, query)
			if err != nil {
				e.logger.Warn("lrclib enrichment failed", zap.Error(err))
			}
		}()
	}

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
