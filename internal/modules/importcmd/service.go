package importcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{logger: logger}
}

type ytDlpEntry struct {
	Title      string  `json:"title"`
	Duration   float64 `json:"duration"`
	WebpageURL string  `json:"webpage_url"`
	Thumbnail  string  `json:"thumbnail"`
	Channel    string  `json:"channel"`
	Uploader   string  `json:"uploader"`
}

func (s *Service) Search(ctx context.Context, query string) ([]SearchResult, error) {
	args := []string{
		"--flat-playlist",
		"--dump-json",
		"--no-warnings",
		"--default-search", "ytsearch5",
		query,
	}

	s.logger.Info("searching with yt-dlp", zap.String("query", query))

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		s.logger.Warn("yt-dlp search failed", zap.String("stderr", stderr.String()), zap.Error(err))
		return nil, fmt.Errorf("search failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	results := make([]SearchResult, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry ytDlpEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			s.logger.Warn("failed to parse yt-dlp output line", zap.String("line", line), zap.Error(err))
			continue
		}

		if entry.Title == "" {
			continue
		}

		artist := entry.Uploader
		if artist == "" {
			artist = entry.Channel
		}

		results = append(results, SearchResult{
			Title:     entry.Title,
			Artist:    artist,
			URL:       entry.WebpageURL,
			Duration:  int(entry.Duration),
			Thumbnail: entry.Thumbnail,
			Source:    "youtube",
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for query: %s", query)
	}

	return results, nil
}

func (s *Service) Download(ctx context.Context, url, dir string) (*ytDlpEntry, string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)

	infoArgs := []string{
		"--dump-json",
		"--no-warnings",
		url,
	}

	s.logger.Info("fetching metadata", zap.String("url", url))

	infoCmd := exec.CommandContext(ctx, "yt-dlp", infoArgs...)
	var infoOut, infoErr bytes.Buffer
	infoCmd.Stdout = &infoOut
	infoCmd.Stderr = &infoErr

	if err := infoCmd.Run(); err != nil {
		s.logger.Warn("metadata fetch failed", zap.String("stderr", infoErr.String()), zap.Error(err))
		return nil, "", fmt.Errorf("metadata fetch failed: %w", err)
	}

	var entry ytDlpEntry
	if err := json.Unmarshal(infoOut.Bytes(), &entry); err != nil {
		return nil, "", fmt.Errorf("parse metadata: %w", err)
	}

	outputTemplate := filepath.Join(dir, "%(title)s-"+timestamp+".%(ext)s")

	dlArgs := []string{
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--output", outputTemplate,
		"--no-warnings",
		"--print", "filename",
		url,
	}

	s.logger.Info("downloading audio", zap.String("url", url))

	dlCmd := exec.CommandContext(ctx, "yt-dlp", dlArgs...)
	var dlOut, dlErr bytes.Buffer
	dlCmd.Stdout = &dlOut
	dlCmd.Stderr = &dlErr

	if err := dlCmd.Run(); err != nil {
		s.logger.Warn("download failed", zap.String("stderr", dlErr.String()), zap.Error(err))
		return nil, "", fmt.Errorf("download failed: %w", err)
	}

	filepath := strings.TrimSpace(dlOut.String())
	if filepath == "" {
		return nil, "", fmt.Errorf("no output file from yt-dlp")
	}

	return &entry, filepath, nil
}
