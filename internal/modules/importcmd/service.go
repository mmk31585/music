package importcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	proxy  string
}

func NewService(logger *zap.Logger) *Service {
	proxy := os.Getenv("IMPORT_PROXY")
	if proxy != "" {
		user := os.Getenv("IMPORT_PROXY_USER")
		pass := os.Getenv("IMPORT_PROXY_PASSWORD")
		if user != "" && pass != "" {
			proxy = fmt.Sprintf("http://%s:%s@%s", user, pass, proxy)
		}
		logger.Info("yt-dlp proxy configured", zap.String("proxy", maskProxy(proxy)))
	}
	return &Service{logger: logger, proxy: proxy}
}

func maskProxy(proxy string) string {
	if idx := strings.LastIndex(proxy, "@"); idx > 0 {
		return "http://***:***@" + proxy[idx+1:]
	}
	return proxy
}

func (s *Service) proxyArgs() []string {
	if s.proxy == "" {
		return nil
	}
	return []string{"--proxy", s.proxy}
}

type YtDlpEntry struct {
	Title      string  `json:"title"`
	Duration   float64 `json:"duration"`
	WebpageURL string  `json:"webpage_url"`
	Thumbnail  string  `json:"thumbnail"`
	Channel    string  `json:"channel"`
	Uploader   string  `json:"uploader"`
}

func (s *Service) Search(ctx context.Context, query string) ([]SearchResult, error) {
	type sourceSearch struct {
		prefix string
		source string
	}

	// Priority ladder: best quality first, stop at first source that returns results.
	// Note: yt-dlp search prefixes differ by version. Current format (2026+):
	//   ytsearch:  - YouTube search (e.g. "ytsearch5:query" for 5 results)
	//   scsearch:  - SoundCloud search
	//   bandcamp and archive.org have no generic search prefix in this version.
	ladder := []sourceSearch{
		{"scsearch", "soundcloud"},
		{"ytsearch5", "youtube"},
	}

	// Cap total search time so the user isn't waiting forever.
	totalCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	start := time.Now()
	var lastErr error

	for _, src := range ladder {
		if totalCtx.Err() != nil {
			s.logger.Warn("search cancelled, total deadline reached",
				zap.String("query", query),
				zap.Duration("elapsed", time.Since(start)),
			)
			break
		}

		s.logger.Info("trying search source",
			zap.String("source", src.source),
			zap.String("query", query),
		)

		results, err := s.searchSource(totalCtx, src.prefix, src.source, query)
		if err != nil {
			s.logger.Warn("source search failed, trying next",
				zap.String("source", src.source),
				zap.Error(err),
			)
			lastErr = err
			continue
		}
		if len(results) == 0 {
			s.logger.Info("source returned no results, trying next",
				zap.String("source", src.source),
			)
			continue
		}

		s.logger.Info("search succeeded",
			zap.String("source", src.source),
			zap.Int("results", len(results)),
			zap.Duration("duration", time.Since(start)),
		)
		return results, nil
	}

	elapsed := time.Since(start)
	s.logger.Warn("all search sources exhausted",
		zap.String("query", query),
		zap.Duration("duration", elapsed),
		zap.Error(lastErr),
	)
	// Return empty slice, not error — the frontend handles "no results" gracefully.
	return []SearchResult{}, nil
}

func (s *Service) searchSource(ctx context.Context, searchPrefix, source, query string) ([]SearchResult, error) {
	args := []string{
		"--flat-playlist",
		"--dump-json",
		"--no-warnings",
		"--no-update",
		"--default-search", searchPrefix,
		"--socket-timeout", "10",
		query,
	}

	// Use parent context (may have a total deadline) but cap per-source at 5s.
	searchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(searchCtx, "yt-dlp", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		elapsed := time.Since(start)
		stderrStr := strings.TrimSpace(stderr.String())
		s.logger.Debug("yt-dlp command failed",
			zap.String("source", source),
			zap.String("query", query),
			zap.Duration("elapsed", elapsed),
			zap.String("stderr", stderrStr),
			zap.Bool("is_timeout", errors.Is(err, context.DeadlineExceeded)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%s search failed after %v: %s (stderr: %s)", source, elapsed.Round(time.Second), err, stderrStr)
	}

	elapsed := time.Since(start)
	output := strings.TrimSpace(stdout.String())

	if output == "" {
		s.logger.Debug("yt-dlp returned empty output",
			zap.String("source", source),
			zap.String("query", query),
			zap.Duration("elapsed", elapsed),
		)
		return nil, nil
	}

	lines := strings.Split(output, "\n")
	results := make([]SearchResult, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry YtDlpEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			s.logger.Warn("failed to parse yt-dlp output line",
				zap.String("source", source),
				zap.String("line", truncate(line, 200)),
				zap.Error(err),
			)
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
			Source:    source,
		})
	}

	s.logger.Info("source search parsed",
		zap.String("source", source),
		zap.Int("results", len(results)),
		zap.Int("raw_lines", len(lines)),
		zap.Duration("elapsed", elapsed),
	)

	return results, nil
}

func (s *Service) Download(ctx context.Context, url, dir string) (*YtDlpEntry, string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)

	infoArgs := []string{
		"--dump-json",
		"--no-warnings",
		"--no-update",
		"--socket-timeout", "10",
	}
	infoArgs = append(infoArgs, s.proxyArgs()...)
	infoArgs = append(infoArgs, url)

	s.logger.Info("fetching metadata", zap.String("url", url))

	infoCmd := exec.CommandContext(ctx, "yt-dlp", infoArgs...)
	var infoOut, infoErr bytes.Buffer
	infoCmd.Stdout = &infoOut
	infoCmd.Stderr = &infoErr

	if err := infoCmd.Run(); err != nil {
		s.logger.Warn("metadata fetch failed", zap.String("stderr", infoErr.String()), zap.Error(err))
		return nil, "", fmt.Errorf("metadata fetch failed: %w", err)
	}

	var entry YtDlpEntry
	if err := json.Unmarshal(infoOut.Bytes(), &entry); err != nil {
		return nil, "", fmt.Errorf("parse metadata: %w", err)
	}

	outputTemplate := filepath.Join(dir, "%(title)s-"+timestamp+".%(ext)s")

	dlArgs := []string{
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--embed-thumbnail",
		"--add-metadata",
		"--output", outputTemplate,
		"--no-playlist",
		"--no-warnings",
		"--no-update",
		"--socket-timeout", "10",
		"--print", "filename",
		"--extractor-args", "youtube:skip=webpage",
	}
	dlArgs = append(dlArgs, s.proxyArgs()...)
	dlArgs = append(dlArgs, url)

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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
