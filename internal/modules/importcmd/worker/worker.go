package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"music/internal/modules/importcmd/acquisition"
	"music/internal/modules/ingestion"
	"music/internal/platform/metrics"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	streamKey     = "stream:import:jobs"
	consumerGroup = "import-workers"
	consumerName  = "worker-1"
)

type Downloader interface {
	Download(ctx context.Context, url, dir string) (entry *DownloadResult, err error)
}

type DownloadResult struct {
	Title     string
	Uploader  string
	Duration  int
	LocalPath string
}

type Resolver interface {
	Resolve(ctx context.Context, q acquisition.ResolveQuery) (*acquisition.Candidate, error)
}

type ytDlpMeta struct {
	Title      string  `json:"title"`
	Duration   float64 `json:"duration"`
	Uploader   string  `json:"uploader"`
	Channel    string  `json:"channel"`
	WebpageURL string  `json:"webpage_url"`
}

type DownloadService struct {
	proxy string
}

func NewDownloadService(proxy string) *DownloadService {
	return &DownloadService{proxy: proxy}
}

func (s *DownloadService) proxyArgs() []string {
	if s.proxy == "" {
		return nil
	}
	return []string{"--proxy", s.proxy}
}

func (s *DownloadService) Download(ctx context.Context, url, dir string) (*DownloadResult, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)

	infoArgs := []string{"--js-runtimes", "node", "--dump-json", "--no-warnings", "--no-update", "--socket-timeout", "10"}
	infoArgs = append(infoArgs, s.proxyArgs()...)
	infoArgs = append(infoArgs, url)

	// Use explicit path to yt-dlp
	ytdlpPath := "/home/space/.local/bin/yt-dlp"
	infoCmd := exec.CommandContext(ctx, ytdlpPath, infoArgs...)
	var infoOut, infoErr bytes.Buffer
	infoCmd.Stdout = &infoOut
	infoCmd.Stderr = &infoErr

	// Clear HTTPS_PROXY/HTTP_PROXY env vars for yt-dlp subprocess
	// to prevent conflicts with --proxy flag. yt-dlp respects HTTPS_PROXY
	// and using both can cause issues in some environments.
	infoCmd.Env = append(infoCmd.Environ(),
		"HTTPS_PROXY=",
		"HTTP_PROXY=",
		"https_proxy=",
		"http_proxy=",
	)

	if err := infoCmd.Run(); err != nil {
		return nil, fmt.Errorf("metadata fetch: %w (stderr: %s)", err, strings.TrimSpace(infoErr.String()))
	}

	var meta ytDlpMeta
	if err := json.Unmarshal(infoOut.Bytes(), &meta); err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	outputTemplate := filepath.Join(dir, "%(title)s-"+timestamp+".%(ext)s")

	dlArgs := []string{
		"--js-runtimes", "node",
		"-x", "--audio-format", "mp3", "--audio-quality", "0",
		"--embed-thumbnail", "--add-metadata",
		"--output", outputTemplate, "--no-playlist",
		"--no-warnings", "--no-update", "--socket-timeout", "10",
		"--extractor-args", "youtube:skip=webpage",
	}
	dlArgs = append(dlArgs, s.proxyArgs()...)
	dlArgs = append(dlArgs, url)

	dlCmd := exec.CommandContext(ctx, "yt-dlp", dlArgs...)
	var dlErr bytes.Buffer
	dlCmd.Stdout = nil // redirected to file by --output
	dlCmd.Stderr = &dlErr

	if err := dlCmd.Run(); err != nil {
		return nil, fmt.Errorf("download: %w (stderr: %s)", err, strings.TrimSpace(dlErr.String()))
	}

	// Construct the output path from the template + metadata.
	// yt-dlp's %(title)s sanitizes filenames by removing \ / : * ? " < > | # chars.
	// We must match that sanitization here, otherwise the path won't match the file
	// yt-dlp actually wrote. See yt-dlp's sanitize_filename().
	sanitize := func(s string) string {
		return strings.NewReplacer(
			`\`, "", "/", "", ":", "", "*", "", "?", "",
			`"`, "", "<", "", ">", "", "|", "", "#", "",
		).Replace(s)
	}
	localPath := filepath.Join(dir, sanitize(meta.Title)+"-"+timestamp+".mp3")

	uploader := meta.Uploader
	if uploader == "" {
		uploader = meta.Channel
	}

	return &DownloadResult{
		Title:     meta.Title,
		Uploader:  uploader,
		Duration:  int(meta.Duration),
		LocalPath: localPath,
	}, nil
}

type ImportWorker struct {
	logger       *zap.Logger
	rdb          redis.UniversalClient
	downloader   Downloader
	ingestionSvc *ingestion.Service
	resolver     Resolver
	progress     *ProgressStore
	downloadDir  string
	pollInterval time.Duration
}

func NewImportWorker(
	logger *zap.Logger,
	rdb redis.UniversalClient,
	downloader Downloader,
	ingestionSvc *ingestion.Service,
	resolver Resolver,
	downloadDir string,
) *ImportWorker {
	return &ImportWorker{
		logger:       logger,
		rdb:          rdb,
		downloader:   downloader,
		ingestionSvc: ingestionSvc,
		resolver:     resolver,
		progress:     NewProgressStore(rdb),
		downloadDir:  downloadDir,
		pollInterval: 5 * time.Second,
	}
}

func (w *ImportWorker) Name() string { return "import_worker" }

func (w *ImportWorker) Run(ctx context.Context) error {
	if err := w.ensureStreamGroup(ctx); err != nil {
		return fmt.Errorf("ensure stream group: %w", err)
	}

	w.logger.Info("import worker started",
		zap.String("stream", streamKey),
		zap.String("group", consumerGroup),
		zap.Duration("poll_interval", w.pollInterval),
	)

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *ImportWorker) ensureStreamGroup(ctx context.Context) error {
	err := w.rdb.XGroupCreateMkStream(ctx, streamKey, consumerGroup, "0").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			// Group already exists — do NOT reset the cursor.
			// Pending messages (from crashed workers) are re-delivered
			// automatically via XReadGroup with ID "0" on the next poll.
			return nil
		}
		return err
	}
	return nil
}

func (w *ImportWorker) poll(ctx context.Context) {
	// Phase 1: drain pending messages (from crashed workers).
	// XReadGroup with ID "0" delivers messages in the PEL
	// (pending entries list) — messages that were delivered but not ACK'd.
	pending, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamKey, "0"},
		Count:    1,
		Block:    1 * time.Second,
	}).Result()

	if err != nil && err != redis.Nil {
		w.logger.Warn("import poll XReadGroup pending error", zap.Error(err))
		return
	}

	if len(pending) > 0 && len(pending[0].Messages) > 0 {
		msg := pending[0].Messages[0]
		w.processMessage(ctx, msg)
		return
	}

	// Phase 2: no pending messages — read new messages with ">".
	newMsgs, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"},
		Count:    1,
		Block:    1 * time.Second,
	}).Result()

	if err != nil {
		if err != redis.Nil {
			w.logger.Warn("import poll XReadGroup error", zap.Error(err))
		}
		return
	}
	if len(newMsgs) == 0 || len(newMsgs[0].Messages) == 0 {
		return
	}

	w.processMessage(ctx, newMsgs[0].Messages[0])
}

func (w *ImportWorker) processMessage(ctx context.Context, msg redis.XMessage) {
	msgID := msg.ID
	w.logger.Debug("import poll received message", zap.String("msg_id", msgID),
		zap.Any("values", msg.Values))

	payload, ok := msg.Values["payload"].(string)
	if !ok {
		w.logger.Warn("import job payload field is not a string",
			zap.String("msg_id", msgID),
			zap.Any("type", fmt.Sprintf("%T", msg.Values["payload"])),
			zap.Any("value", msg.Values["payload"]))
		w.rdb.XAck(ctx, streamKey, consumerGroup, msgID)
		return
	}

	w.logger.Debug("import poll payload string", zap.String("msg_id", msgID),
		zap.String("payload_len", strconv.Itoa(len(payload))),
		zap.String("payload_prefix", payload[:min(len(payload), 100)]))

	// The payload is double-wrapped: Enqueue stores {"payload":{job JSON}}
	// so we need to unwrap the outer payload field first.
	var job Job
	var wrapper struct {
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal([]byte(payload), &wrapper); err != nil {
		w.logger.Warn("invalid job wrapper", zap.Error(err))
		w.rdb.XAck(ctx, streamKey, consumerGroup, msgID)
		return
	}
	if err := json.Unmarshal(wrapper.Payload, &job); err != nil {
		w.logger.Warn("invalid job payload", zap.Error(err))
		w.rdb.XAck(ctx, streamKey, consumerGroup, msgID)
		return
	}

	metrics.SetImportJobsInFlight(1)
	w.processJob(ctx, &job)
	metrics.SetImportJobsInFlight(0)

	w.rdb.XAck(ctx, streamKey, consumerGroup, msgID)
}

func (w *ImportWorker) Enqueue(ctx context.Context, job *Job) error {
	payload, err := json.Marshal(map[string]any{"payload": job})
	if err != nil {
		return err
	}

	// Write job record to Redis hash
	job.Status = StatusQueued
	w.progress.Set(ctx, job.ID, StatusQueued, 0, "queued", "", "")

	_, err = w.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]any{
			"payload": string(payload),
		},
	}).Result()
	return err
}

func (w *ImportWorker) processJob(ctx context.Context, job *Job) {
	jobStart := time.Now()

	w.logger.Info("processing import job",
		zap.String("job_id", job.ID),
		zap.String("url", job.URL),
		zap.String("source", job.Source),
	)

	w.setProgress(job.ID, StatusResolving, 5, "resolving", "", "")

	// Resolve to downloadable URL if needed (e.g., Spotify URL needs yt-dlp resolve)
	downloadURL := job.URL

	// iTunes preview URLs are 30-second clips, not full tracks.
	// Force resolution via acquisition (SoundCloud/YouTube) to get the full song.
	if isPreviewURL(downloadURL) {
		w.logger.Info("iTunes preview URL detected, resolving full track",
			zap.String("preview_url", downloadURL),
		)
		downloadURL = ""
	}

	if downloadURL == "" {
		resolveQuery := acquisition.ResolveQuery{
			Title:  job.Title,
			Artist: job.Artist,
		}

		// Resolve via acquisition (SoundCloud/YouTube) when:
		//   1. The source is a known metadata provider (Spotify, Deezer, etc.), OR
		//   2. We have at least title+artist and no URL (manual/queue import without a source).
		if job.Source == "spotify" || job.Source == "deezer" || job.Source == "musicbrainz" || job.Source == "lastfm" || job.Source == "itunes" ||
			(job.Title != "" && job.Artist != "") {
			resolveStart := time.Now()
			candidate, err := w.resolver.Resolve(ctx, resolveQuery)
			metrics.ObserveImportJobStage("resolve", resolveStart)
			if err != nil {
				w.setProgress(job.ID, StatusFailed, 0, "resolve_failed", "", "Could not find a downloadable source for this track. Try importing with a direct YouTube or SoundCloud URL instead.")
				metrics.IncImportJob("failed")
				return
			}
			if candidate == nil {
				w.setProgress(job.ID, StatusFailed, 0, "resolve_failed", "", "No downloadable source found for this track. The track may not be available on YouTube or SoundCloud.")
				metrics.IncImportJob("failed")
				return
			}
			downloadURL = candidate.URL
			w.logger.Info("track resolved to downloadable URL",
				zap.String("source", candidate.Source),
				zap.String("quality", candidate.Quality),
				zap.String("url", downloadURL),
			)
		} else {
			// URL is empty and source doesn't need resolution — can't proceed
			w.setProgress(job.ID, StatusFailed, 0, "resolve_failed", "", "Provide a YouTube or SoundCloud URL to import this track, or search for it first.")
			metrics.IncImportJob("failed")
			return
		}
	} else {
		w.logger.Info("using provided download URL, skipping resolution",
			zap.String("url", downloadURL),
		)
	}

	w.setProgress(job.ID, StatusDownloading, 20, "downloading", "", "")

	if err := os.MkdirAll(w.downloadDir, 0755); err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "mkdir_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}

	// Debug: log the exact proxy and URL before download
	w.logger.Info("worker debug: starting download",
		zap.String("job_id", job.ID),
		zap.String("url", downloadURL),
		zap.String("download_dir", w.downloadDir),
		zap.String("IMPORT_PROXY", os.Getenv("IMPORT_PROXY")),
	)

	dlStart := time.Now()
	result, err := w.downloader.Download(ctx, downloadURL, w.downloadDir)
	metrics.ObserveImportJobStage("download", dlStart)
	if err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "download_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}

	// Override downloaded metadata with search result metadata when available.
	// This is important for sources like iTunes where yt-dlp uses a generic
	// filename instead of the actual track title.
	if job.Title != "" {
		result.Title = job.Title
	}
	if job.Artist != "" {
		result.Uploader = job.Artist
	}
	if job.Source != "" {
		// Use proper source attribution
	}

	w.setProgress(job.ID, StatusExtracting, 50, "extracting", "", "")
	extractStart := time.Now()

	f, err := os.Open(result.LocalPath)
	if err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "open_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		os.Remove(result.LocalPath)
		w.setProgress(job.ID, StatusFailed, 0, "stat_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}
	metrics.ObserveImportJobStage("extract", extractStart)

	w.setProgress(job.ID, StatusUploading, 75, "uploading", "", "")
	uploadStart := time.Now()

	// Use the actual file extension from the downloaded file
	// instead of hardcoding .mp3 (iTunes previews are .m4a, etc.)
	ext := ".mp3"
	if idx := strings.LastIndex(result.LocalPath, "."); idx != -1 {
		ext = result.LocalPath[idx:]
	}
	header := &multipart.FileHeader{
		Filename: result.Title + ext,
		Size:     stat.Size(),
	}

	uploadRes, err := w.ingestionSvc.Upload(ctx, f, header, job.UserID)
	f.Close()
	os.Remove(result.LocalPath)

	metrics.ObserveImportJobStage("upload", uploadStart)
	if err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "upload_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}

	job.DraftID = uploadRes.DraftID

	// Override extracted metadata with the job's title/artist/album if available.
	// This fixes the case where yt-dlp infers a garbage title from the URL
	// (e.g., iTunes previews with generic mzaf_ filenames).
	if job.Title != "" || job.Artist != "" || job.Album != "" {
		if err := w.ingestionSvc.UpdateExtractedMetadata(ctx, uploadRes.DraftID, job.Title, job.Artist, job.Album); err != nil {
			w.logger.Warn("failed to update extracted metadata",
				zap.String("draft_id", uploadRes.DraftID),
				zap.Error(err),
			)
		} else {
			w.logger.Info("overrode extracted metadata",
				zap.String("draft_id", uploadRes.DraftID),
				zap.String("title", job.Title),
				zap.String("artist", job.Artist),
			)
		}
	}

	w.setProgress(job.ID, StatusComplete, 100, "complete", uploadRes.DraftID, "")

	metrics.ObserveImportJobStage("total", jobStart)
	metrics.IncImportJob("complete")

	w.logger.Info("import job completed",
		zap.String("job_id", job.ID),
		zap.String("draft_id", uploadRes.DraftID),
		zap.String("title", uploadRes.OriginalFilename),
	)
}

func (w *ImportWorker) GetProgress(ctx context.Context, jobID string) (*ProgressUpdate, error) {
	return w.progress.Get(ctx, jobID)
}

func (w *ImportWorker) setProgress(jobID string, status Status, progress int, stage string, draftID string, errMsg string) {
	w.progress.Set(context.Background(), jobID, status, progress, stage, draftID, errMsg)

	w.logger.Info("import progress",
		zap.String("job_id", jobID),
		zap.String("status", string(status)),
		zap.Int("progress", progress),
		zap.String("stage", stage),
		zap.String("draft_id", draftID),
		zap.String("error", errMsg),
	)
}

// isPreviewURL detects iTunes preview URLs (30-second clips).
// These URLs contain "itunes.apple.com" and "AudioPreview" in the path
// and point to short preview clips, not full tracks.
func isPreviewURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}
	return strings.Contains(rawURL, "itunes.apple.com") && strings.Contains(rawURL, "AudioPreview")
}
