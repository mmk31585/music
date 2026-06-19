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

	infoArgs := []string{"--dump-json", "--no-warnings", "--no-update", "--socket-timeout", "10"}
	infoArgs = append(infoArgs, s.proxyArgs()...)
	infoArgs = append(infoArgs, url)

	infoCmd := exec.CommandContext(ctx, "yt-dlp", infoArgs...)
	var infoOut, infoErr bytes.Buffer
	infoCmd.Stdout = &infoOut
	infoCmd.Stderr = &infoErr

	if err := infoCmd.Run(); err != nil {
		return nil, fmt.Errorf("metadata fetch: %w (stderr: %s)", err, strings.TrimSpace(infoErr.String()))
	}

	var meta ytDlpMeta
	if err := json.Unmarshal(infoOut.Bytes(), &meta); err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	outputTemplate := filepath.Join(dir, "%(title)s-"+timestamp+".%(ext)s")

	dlArgs := []string{
		"-x", "--audio-format", "mp3", "--audio-quality", "0",
		"--embed-thumbnail", "--add-metadata",
		"--output", outputTemplate, "--no-playlist",
		"--no-warnings", "--no-update", "--socket-timeout", "10",
		"--print", "filename", "--extractor-args", "youtube:skip=webpage",
	}
	dlArgs = append(dlArgs, s.proxyArgs()...)
	dlArgs = append(dlArgs, url)

	dlCmd := exec.CommandContext(ctx, "yt-dlp", dlArgs...)
	var dlOut, dlErr bytes.Buffer
	dlCmd.Stdout = &dlOut
	dlCmd.Stderr = &dlErr

	if err := dlCmd.Run(); err != nil {
		return nil, fmt.Errorf("download: %w (stderr: %s)", err, strings.TrimSpace(dlErr.String()))
	}

	localPath := strings.TrimSpace(dlOut.String())
	if localPath == "" {
		return nil, fmt.Errorf("no output file from yt-dlp")
	}

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
	err := w.rdb.XGroupCreateMkStream(ctx, streamKey, consumerGroup, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

func (w *ImportWorker) poll(ctx context.Context) {
	msgs, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"},
		Count:    1,
		Block:    1 * time.Second,
	}).Result()

	if err != nil || len(msgs) == 0 || len(msgs[0].Messages) == 0 {
		return
	}

	msg := msgs[0].Messages[0]
	msgID := msg.ID

	var job Job
	payload, ok := msg.Values["payload"].(string)
	if !ok {
		w.rdb.XAck(ctx, streamKey, consumerGroup, msgID)
		return
	}
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
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
	resolveQuery := acquisition.ResolveQuery{
		Title:  job.Title,
		Artist: job.Artist,
	}

	var downloadURL string
	if job.Source == "spotify" || job.Source == "deezer" || job.Source == "musicbrainz" || job.Source == "lastfm" {
		resolveStart := time.Now()
		candidate, err := w.resolver.Resolve(ctx, resolveQuery)
		metrics.ObserveImportJobStage("resolve", resolveStart)
		if err != nil {
			w.setProgress(job.ID, StatusFailed, 0, "resolve_failed", "", err.Error())
			metrics.IncImportJob("failed")
			return
		}
		if candidate == nil {
			w.setProgress(job.ID, StatusFailed, 0, "resolve_failed", "", "no downloadable source found")
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
		downloadURL = job.URL
	}

	w.setProgress(job.ID, StatusDownloading, 20, "downloading", "", "")

	if err := os.MkdirAll(w.downloadDir, 0755); err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "mkdir_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
	}

	dlStart := time.Now()
	result, err := w.downloader.Download(ctx, downloadURL, w.downloadDir)
	metrics.ObserveImportJobStage("download", dlStart)
	if err != nil {
		w.setProgress(job.ID, StatusFailed, 0, "download_failed", "", err.Error())
		metrics.IncImportJob("failed")
		return
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

	header := &multipart.FileHeader{
		Filename: result.Title + ".mp3",
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
