package lyrics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"music/internal/modules/lyrics/lrc"
)

var (
	ErrForbiddenLyricsAccess = errors.New("forbidden lyrics access")
	ErrInvalidLyricsLanguage = errors.New("invalid lyrics language")
)

var lrcPattern = regexp.MustCompile(`(?m)^\[\d{1,2}:\d{2}[\.:]\d{2,3}\]`)

func isLRCLyrics(content string) bool {
	return lrcPattern.MatchString(content)
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLyrics(ctx context.Context, req CreateLyricsRequest) (Lyrics, error) {
	// Validate language
	req.Language = strings.TrimSpace(req.Language)
	if req.Language == "" {
		return Lyrics{}, ErrInvalidLyricsLanguage
	}

	// Auto-detect LRC from content
	req.Type = strings.TrimSpace(req.Type)
	if isLRCLyrics(req.Content) {
		req.Type = "lrc"
	}

	if req.Type != "plain" && req.Type != "lrc" {
		return Lyrics{}, ErrInvalidLyricsType
	}

	return s.repo.CreateLyrics(ctx, req)
}

func (s *Service) UpdateLyrics(ctx context.Context, lyricsID string, req UpdateLyricsRequest) (Lyrics, error) {
	// Validate language if provided
	if req.Language != nil {
		*req.Language = strings.TrimSpace(*req.Language)
		if *req.Language == "" {
			return Lyrics{}, ErrInvalidLyricsLanguage
		}
	}

	// Auto-detect LRC from content
	if req.Content != nil && req.Type != nil && *req.Type == "plain" && isLRCLyrics(*req.Content) {
		*req.Type = "lrc"
	}

	// Validate type if provided
	if req.Type != nil {
		*req.Type = strings.TrimSpace(*req.Type)
		if *req.Type != "" && *req.Type != "plain" && *req.Type != "lrc" {
			return Lyrics{}, ErrInvalidLyricsType
		}
	}

	return s.repo.UpdateLyrics(ctx, lyricsID, req)
}

func (s *Service) DeleteLyrics(ctx context.Context, lyricsID string) error {
	return s.repo.DeleteLyrics(ctx, lyricsID)
}

func (s *Service) GetLyricsByTrackID(ctx context.Context, trackID string) ([]Lyrics, error) {
	return s.repo.GetLyricsByTrackID(ctx, trackID)
}

func (s *Service) GetLyricsByTrackAndLanguage(ctx context.Context, trackID, language string) (LyricsByTrackResponse, error) {
	lyrics, err := s.repo.GetLyricsByTrackAndLanguage(ctx, trackID, language)
	if err != nil {
		return LyricsByTrackResponse{}, err
	}

	var lines []LyricsLine // This is lyrics.LyricsLine
	if lyrics.Type == "lrc" {
		parsedLines, err := lrc.ParseLRC(lyrics.Content)
		if err != nil {
			return LyricsByTrackResponse{}, err
		}

		// Map lrc.LyricsLine to lyrics.LyricsLine
		lines = make([]LyricsLine, len(parsedLines))
		for i, v := range parsedLines {
			lines[i] = LyricsLine{
				TimeSeconds: v.TimeSeconds,
				Text:        v.Text,
			}
		}
	}

	return ToLyricsByTrackResponse(lyrics, lines), nil
}

// FetchFromLRC queries LRCLIB for synced lyrics and creates/updates a lyrics record.
func (s *Service) FetchFromLRC(ctx context.Context, trackID string) (Lyrics, error) {
	// 1. Get track info
	info, err := s.repo.GetTrackInfo(ctx, trackID)
	if err != nil {
		return Lyrics{}, err
	}
	if info.Title == "" {
		return Lyrics{}, ErrTrackNoTitle
	}

	// 2. Call LRCLIB
	result, err := fetchLRCLib(ctx, info.Title, info.ArtistName, info.DurationSeconds)
	if err != nil {
		return Lyrics{}, ErrLRCLibNoLyrics
	}
	if result == nil {
		return Lyrics{}, ErrLRCLibNoLyrics
	}

	// Prefer synced lyrics
	content := result.SyncedLyrics
	lrcType := "lrc"
	if content == "" {
		content = result.PlainLyrics
		lrcType = "plain"
	}
	if content == "" {
		return Lyrics{}, ErrLRCLibEmptyResponse
	}

	// 3. Remove existing lyrics for this track + "en" if they exist
	existing, lookupErr := s.repo.GetLyricsByTrackAndLanguage(ctx, trackID, "en")
	if lookupErr == nil {
		_ = s.repo.DeleteLyrics(ctx, existing.ID)
	}

	// 4. Create new lyrics record
	req := CreateLyricsRequest{
		TrackID:  trackID,
		Language: "en",
		Type:     lrcType,
		Content:  content,
	}
	return s.CreateLyrics(ctx, req)
}

// fetchLRCLib calls LRCLIB API: first exact /get, then fallback /search.
func fetchLRCLib(ctx context.Context, trackName, artistName string, duration int) (*lrclibResult, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseURL := "https://lrclib.net/api"

	// Step 1: Try exact /get
	params := url.Values{}
	if artistName != "" {
		params.Set("artist_name", artistName)
	}
	if trackName != "" {
		params.Set("track_name", trackName)
	}
	if duration > 0 {
		params.Set("duration", fmt.Sprintf("%d", duration))
	}

	if len(params) > 0 {
		getURL := fmt.Sprintf("%s/get?%s", baseURL, params.Encode())
		req, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
		if err == nil {
			req.Header.Set("User-Agent", "MojaMusic/1.0 (music-admin)")
			req.Header.Set("Accept", "application/json")

			resp, err := httpClient.Do(req)
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var result lrclibResult
					if json.NewDecoder(resp.Body).Decode(&result) == nil {
						resp.Body.Close()
						return &result, nil
					}
					resp.Body.Close()
				} else {
					resp.Body.Close()
				}
			}
		}
	}

	// Step 2: Fallback to /search
	if trackName == "" && artistName == "" {
		return nil, nil
	}
	searchQ := trackName
	if artistName != "" {
		searchQ = artistName + " " + searchQ
	}

	searchURL := fmt.Sprintf("%s/search?q=%s", baseURL, url.QueryEscape(searchQ))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MojaMusic/1.0 (music-admin)")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil // silent fallback
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var results []lrclibResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, nil
	}

	// Prefer synced lyrics
	var best *lrclibResult
	for i := range results {
		r := &results[i]
		if r.SyncedLyrics != "" {
			return r, nil
		}
		if best == nil && r.PlainLyrics != "" {
			best = r
		}
	}
	return best, nil
}

// lrclibResult mirrors the LRCLIB API response.
type lrclibResult struct {
	ID           int    `json:"id"`
	TrackName    string `json:"trackName"`
	ArtistName   string `json:"artistName"`
	AlbumName    string `json:"albumName"`
	Duration     int    `json:"duration"`
	Synced       bool   `json:"synced"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

func (s *Service) GetTrackLyrics(ctx context.Context, trackID, language string) (LyricsByTrackResponse, error) {
	var lyrics Lyrics
	var err error

	if language == "" {
		lyricsList, err := s.repo.GetLyricsByTrackID(ctx, trackID)
		if err != nil || len(lyricsList) == 0 {
			return LyricsByTrackResponse{}, ErrLyricsNotFound
		}
		lyrics = lyricsList[0]
	} else {
		lyrics, err = s.repo.GetLyricsByTrackAndLanguage(ctx, trackID, language)
		if err != nil {
			return LyricsByTrackResponse{}, err
		}
	}

	var lines []LyricsLine
	if lyrics.Type == "lrc" {
		parsedLines, err := lrc.ParseLRC(lyrics.Content)
		if err != nil {
			return LyricsByTrackResponse{}, err
		}

		// Conversion loop
		lines = make([]LyricsLine, len(parsedLines))
		for i, v := range parsedLines {
			lines[i] = LyricsLine{
				TimeSeconds: v.TimeSeconds,
				Text:        v.Text,
			}
		}
	}

	return ToLyricsByTrackResponse(lyrics, lines), nil
}

func (s *Service) LyricsExists(ctx context.Context, trackID, language string) (bool, error) {
	return s.repo.LyricsExists(ctx, trackID, language)
}

// GetTrackInfo fetches track metadata needed for lyrics operations.
func (s *Service) GetTrackInfo(ctx context.Context, trackID string) (*TrackInfo, error) {
	return s.repo.GetTrackInfo(ctx, trackID)
}

// FetchOrGenerateResult describes the outcome of a fetch-or-generate attempt.
type FetchOrGenerateResult struct {
	Source  string  // "lrclib" | "ai" | "none"
	Lyrics  *Lyrics // non-nil when source == "lrclib"
	JobID   string  // non-empty when source == "ai"
	TrackID string
}
