package acquisition

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ytDlpEntry struct {
	Title      string  `json:"title"`
	Duration   float64 `json:"duration"`
	WebpageURL string  `json:"webpage_url"`
	Thumbnail  string  `json:"thumbnail"`
	Channel    string  `json:"channel"`
	Uploader   string  `json:"uploader"`
}

type ytDlpResolver struct {
	name     string
	prefix   string
	quality  string
	priority int
	proxy    string
}

func newSoundCloudResolver(proxy string) *ytDlpResolver {
	return &ytDlpResolver{
		name:     "soundcloud",
		prefix:   "scsearch",
		quality:  QualityV0,
		priority: 1,
		proxy:    proxy,
	}
}

func newYTMusicResolver(proxy string) *ytDlpResolver {
	return &ytDlpResolver{
		name:     "youtube",
		prefix:   "ytsearch5",
		quality:  Quality128,
		priority: 2,
		proxy:    proxy,
	}
}

func (r *ytDlpResolver) Name() string    { return r.name }
func (r *ytDlpResolver) Priority() int   { return r.priority }
func (r *ytDlpResolver) Quality() string { return r.quality }

func (r *ytDlpResolver) buildArgs(query string) []string {
	args := []string{
		"--flat-playlist",
		"--dump-json",
		"--no-warnings",
		"--no-update",
		"--default-search", r.prefix,
		"--socket-timeout", "10",
	}
	if r.proxy != "" {
		args = append(args, "--proxy", r.proxy)
	}
	args = append(args, query)
	return args
}

func (r *ytDlpResolver) Resolve(ctx context.Context, q ResolveQuery) (*Candidate, error) {
	query := fmt.Sprintf("%s - %s", q.Artist, q.Title)

	searchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(searchCtx, "yt-dlp", r.buildArgs(query)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s resolve: %w (stderr: %s)", r.name, err, strings.TrimSpace(stderr.String()))
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry ytDlpEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.Title == "" || entry.WebpageURL == "" {
			continue
		}

		return &Candidate{
			URL:         entry.WebpageURL,
			Source:      r.name,
			Quality:     r.quality,
			Confidence:  0.9,
			AudioFormat: "mp3",
		}, nil
	}

	return nil, nil
}
