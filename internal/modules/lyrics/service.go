package lyrics

import (
	"context"
	"errors"
	"regexp"
	"strings"

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
