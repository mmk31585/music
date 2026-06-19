package app

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"music/internal/modules/importcmd/acquisition"
	importworker "music/internal/modules/importcmd/worker"
	"music/internal/modules/ingestion"
	"music/internal/modules/ingestion/enrichment"
	platformstorage "music/internal/platform/storage"
	"music/internal/workers/analytics"
	"music/internal/workers/cleanup"
	"music/internal/workers/notification"
	"music/internal/workers/recommendation"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BackgroundWorker interface {
	Run(ctx context.Context) error
	Name() string
}

type WorkerContainer struct {
	App *App

	AnalyticsWorker      BackgroundWorker
	RecommendationWorker BackgroundWorker
	NotificationWorker   BackgroundWorker
	EmailWorker          BackgroundWorker
	CleanupWorker        BackgroundWorker
	ImportWorker         BackgroundWorker

	workers []BackgroundWorker
}

func NewWorkerContainer(a *App) *WorkerContainer {
	c := &WorkerContainer{
		App: a,
	}

	c.buildWorkers()

	return c
}

func (c *WorkerContainer) buildWorkers() {
	c.AnalyticsWorker = analytics.NewWorker(c.App.DB, c.App.Logger)
	c.RecommendationWorker = recommendation.NewWorker(c.App.DB, c.App.Logger)
	c.NotificationWorker = notification.NewWorker(c.App.DB, c.App.Logger)
	c.CleanupWorker = cleanup.NewWorker(c.App.DB, c.App.Logger)

	proxy := os.Getenv("IMPORT_PROXY")
	downloadDir := os.Getenv("IMPORT_DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = filepath.Join(os.TempDir(), "music-imports")
	}

	dlSvc := importworker.NewDownloadService(proxy)
	acqSvc := acquisition.NewService(c.App.Logger, proxy)

	sqlDB := stdlib.OpenDBFromPool(c.App.DB)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")

	ingestionRepo := ingestion.NewRepository(sqlxDB)

	enrichCfg := enrichment.DefaultConfig()
	enrichCfg.LastFM.APIKey = c.App.Config.Enrichment.LastFmAPIKey
	enrichCfg.Spotify.ClientID = c.App.Config.Enrichment.SpotifyClientID
	enrichCfg.Spotify.ClientSecret = c.App.Config.Enrichment.SpotifyClientSecret

	mbClient := enrichment.NewMusicBrainzClient(enrichCfg.MusicBrainz)
	lfmClient := enrichment.NewLastFMClient(enrichCfg.LastFM)
	spotClient := enrichment.NewSpotifyClient(enrichCfg.Spotify)
	lrcClient := enrichment.NewLRCLibClient()
	enricher := enrichment.NewEnricher(mbClient, lfmClient, spotClient, lrcClient, c.App.Logger)

	storageClient, err := platformstorage.New(context.Background(), platformstorage.Config{
		Driver: c.App.Config.Storage.Driver,
		Local: platformstorage.LocalConfig{
			BaseDir: c.App.Config.Storage.Local.BaseDir,
			BaseURL: c.App.Config.Storage.Local.BaseURL,
		},
		S3: platformstorage.S3Config{
			Bucket:          c.App.Config.Storage.S3.Bucket,
			Region:          c.App.Config.Storage.S3.Region,
			Endpoint:        c.App.Config.Storage.S3.Endpoint,
			AccessKeyID:     c.App.Config.Storage.S3.AccessKeyID,
			SecretAccessKey: c.App.Config.Storage.S3.SecretAccessKey,
			PublicBaseURL:   c.App.Config.Storage.S3.PublicBaseURL,
			UsePathStyle:    c.App.Config.Storage.S3.UsePathStyle,
			PresignURLs:     c.App.Config.Storage.S3.PresignURLs,
			PresignTTL:      c.App.Config.Storage.S3.PresignTTL,
		},
	})
	if err != nil {
		c.App.Logger.Fatal("failed to create storage for worker", zap.Error(err))
	}

	ingestionSvc := ingestion.NewService(storageClient, ingestionRepo, enricher)

	c.ImportWorker = importworker.NewImportWorker(
		c.App.Logger,
		c.App.Redis,
		dlSvc,
		ingestionSvc,
		acqSvc,
		downloadDir,
	)

	c.workers = []BackgroundWorker{
		c.AnalyticsWorker,
		c.RecommendationWorker,
		c.NotificationWorker,
		c.EmailWorker,
		c.CleanupWorker,
		c.ImportWorker,
	}
}

func (c *WorkerContainer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, worker := range c.workers {
		if worker == nil {
			continue
		}

		wg.Add(1)

		go func(w BackgroundWorker) {
			defer wg.Done()

			if err := w.Run(ctx); err != nil && err != context.Canceled {
				c.App.Logger.Error("worker failed", zap.String("name", w.Name()), zap.Error(err))
				errCh <- err
			}
		}(worker)
	}

	select {
	case <-ctx.Done():
		wg.Wait()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (c *WorkerContainer) Shutdown(ctx context.Context) error {
	return nil
}
