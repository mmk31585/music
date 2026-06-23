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
	ml          MLClient
	logger      *zap.Logger
}

func NewEnricher(mb MusicBrainzClient, lfm LastFMClient, spot SpotifyClient, lrc LRCLibClient, ml MLClient, logger *zap.Logger) *Enricher {
	return &Enricher{
		musicBrainz: mb,
		lastFM:      lfm,
		spotify:     spot,
		lrclib:      lrc,
		ml:          ml,
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
		mlResult   *MLResult
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
				e.logger.Warn("last.fm enrichment failed", zap.Error(err))
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
				e.logger.Warn("spotify enrichment failed", zap.Error(err))
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

	if runSource(SourceML) {
		wg.Add(2)
		go func() {
			defer wg.Done()
			var err error
			mlResult, err = e.ml.FetchLyrics(ctx, query)
			if err != nil {
				e.logger.Warn("ml lyrics enrichment failed", zap.Error(err))
			}
		}()
		go func() {
			defer wg.Done()
			if result, err := e.ml.FetchCover(ctx, query); err != nil {
				e.logger.Warn("ml cover enrichment failed", zap.Error(err))
			} else if result != nil {
				if result.AlbumCoverURL != "" {
					if mlResult == nil {
						mlResult = result
					} else {
						mlResult.AlbumCoverURL = result.AlbumCoverURL
					}
				}
			}
		}()
	}

	wg.Wait()

	result := MergeResults(mbResult, lfmResult, spotResult, lrcResult)
	result.ML = mlResult
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
