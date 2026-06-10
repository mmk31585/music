package transcoder

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"music/internal/platform/audio"
)

type Quality struct {
	Name    string
	Bitrate string
}

var Qualities = []Quality{
	{Name: "128k", Bitrate: "128k"},
	{Name: "192k", Bitrate: "192k"},
	{Name: "320k", Bitrate: "320k"},
}

type Worker struct {
	db         *sqlx.DB
	logger     *zap.Logger
	ffmpegPath string
	outputDir  string
	analyzer   *audio.Analyzer
}

func NewWorker(db *sqlx.DB, logger *zap.Logger, ffmpegPath, outputDir string) *Worker {
	return &Worker{
		db:         db,
		logger:     logger,
		ffmpegPath: ffmpegPath,
		outputDir:  outputDir,
		analyzer:   audio.NewAnalyzer("ffprobe", ffmpegPath),
	}
}

func (w *Worker) Name() string { return "transcoder" }

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.processJobs(ctx); err != nil {
				w.logger.Error("transcoder worker failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) processJobs(ctx context.Context) error {
	rows, err := w.db.QueryContext(ctx, `
		SELECT id, track_id, source_path, qualities
		FROM transcode_jobs
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT 5
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, trackID, sourcePath string
		var qualities []string
		if err := rows.Scan(&id, &trackID, &sourcePath, &qualities); err != nil {
			continue
		}

		if err := w.transcode(ctx, id, trackID, sourcePath, qualities); err != nil {
			w.logger.Error("transcode failed", zap.String("job_id", id), zap.Error(err))
			w.db.ExecContext(ctx,
				`UPDATE transcode_jobs SET status = 'failed', error_message = $1 WHERE id = $2`,
				err.Error(), id,
			)
		}
	}
	return rows.Err()
}

func (w *Worker) transcode(ctx context.Context, jobID, trackID, sourcePath string, qualities []string) error {
	w.db.ExecContext(ctx,
		`UPDATE transcode_jobs SET status = 'processing', started_at = NOW() WHERE id = $1`,
		jobID,
	)

	normalizedPath := filepath.Join(w.outputDir, trackID, "normalized.flac")
	if err := w.analyzer.NormalizeLoudness(ctx, sourcePath, normalizedPath); err != nil {
		return fmt.Errorf("normalize loudness: %w", err)
	}

	trackDir := filepath.Join(w.outputDir, trackID, "hls")
	hlsPath := filepath.Join(trackDir, "master.m3u8")

	var mapArgs []string
	for i, q := range qualities {
		mapArgs = append(mapArgs,
			"-map", "0:a",
			fmt.Sprintf("-b:a:%d", i), q,
			"-var_stream_map", fmt.Sprintf("a:0,name:%s", q),
		)
	}

	hlsArgs := []string{
		"-i", normalizedPath,
		"-filter_complex", buildComplexFilter(qualities),
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(trackDir, "%v", "seg_%03d.ts"),
		"-master_pl_name", "master.m3u8",
		"-strftime", "1",
		"-y",
		filepath.Join(trackDir, "%v", "playlist.m3u8"),
	}

	cmd := exec.CommandContext(ctx, w.ffmpegPath, hlsArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg hls: %w\noutput: %s", err, string(output))
	}

	analysis, err := w.analyzer.Analyze(ctx, normalizedPath)
	if err != nil {
		w.logger.Warn("audio analysis failed, continuing without it",
			zap.String("job_id", jobID), zap.Error(err))
	}

	parsedTrackID, _ := uuid.Parse(trackID)

	if analysis != nil {
		waveformData := audio.WaveformData{
			TrackID:         trackID,
			Format:          "json",
			SampleRate:      analysis.WaveformSampleRate,
			DurationSeconds: analysis.DurationSeconds,
			Peaks:           analysis.WaveformPeaks,
			GeneratedAt:     time.Now().UTC(),
		}

		waveformJSON, err := json.Marshal(waveformData)
		if err == nil {
			w.db.ExecContext(ctx, `
				INSERT INTO waveform_data (track_id, data, sample_rate, duration_seconds, generated_at)
				VALUES ($1, $2::jsonb, $3, $4, NOW())
				ON CONFLICT (track_id) DO UPDATE SET data = $2::jsonb, generated_at = NOW()
			`, parsedTrackID, string(waveformJSON), analysis.WaveformSampleRate, analysis.DurationSeconds)
		}

		w.db.ExecContext(ctx, `
			UPDATE tracks SET
				duration_seconds = CASE WHEN duration_seconds IS NULL THEN $2 ELSE duration_seconds END,
				tempo = $3,
				audio_key = $4,
				integrated_loudness = $5,
				energy = $6,
				danceability = $7,
				acousticness = $8
			WHERE id = $1
		`, parsedTrackID, analysis.DurationSeconds, analysis.Tempo, analysis.Key,
			analysis.IntegratedLoudness, analysis.Energy, analysis.Danceability, analysis.Acousticness)

		for _, quality := range qualities {
			qualityPath := filepath.Join(trackDir, quality, "playlist.m3u8")
			w.db.ExecContext(ctx, `
				INSERT INTO audio_qualities (track_id, quality, file_path, bitrate, container, created_at)
				VALUES ($1, $2, $3, $4, 'mpegts', NOW())
				ON CONFLICT (track_id, quality) DO UPDATE SET file_path = $3
			`, parsedTrackID, quality, qualityPath, quality)
		}
	}

	w.db.ExecContext(ctx, `
		UPDATE tracks SET hls_path = $1, has_hls = true, hls_qualities = $2 WHERE id = $3
	`, hlsPath, qualities, parsedTrackID)

	w.db.ExecContext(ctx, `
		UPDATE transcode_jobs SET status = 'completed', completed_at = NOW() WHERE id = $1
	`, jobID)

	return nil
}

func buildComplexFilter(qualities []string) string {
	filter := ""
	for i := range qualities {
		if i > 0 {
			filter += "; "
		}
		filter += fmt.Sprintf("[0:a]aformat=sample_rates=44100|48000:channel_layouts=stereo[a%d]", i)
	}
	return filter
}
