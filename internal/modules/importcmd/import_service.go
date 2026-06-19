package importcmd

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"music/internal/modules/importcmd/acquisition"
	"music/internal/modules/importcmd/search"
	"music/internal/modules/importcmd/worker"
	"music/internal/modules/ingestion"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImportService struct {
	logger       *zap.Logger
	downloader   *Service
	ingestionSvc *ingestion.Service
	resolver     *acquisition.Service
	searchSvc    *search.Service
	importWorker *worker.ImportWorker
	downloadDir  string
}

func NewImportService(
	logger *zap.Logger,
	downloader *Service,
	ingestionSvc *ingestion.Service,
	resolver *acquisition.Service,
	searchSvc *search.Service,
	importWorker *worker.ImportWorker,
	downloadDir string,
) *ImportService {
	return &ImportService{
		logger:       logger,
		downloader:   downloader,
		ingestionSvc: ingestionSvc,
		resolver:     resolver,
		searchSvc:    searchSvc,
		importWorker: importWorker,
		downloadDir:  downloadDir,
	}
}

func (s *ImportService) Search(ctx context.Context, query string) ([]search.Result, error) {
	results, err := s.searchSvc.Search(ctx, query)
	if err != nil {
		s.logger.Warn("API search failed, falling back to yt-dlp search", zap.Error(err))
		return s.fallbackSearch(ctx, query)
	}
	if len(results) == 0 {
		s.logger.Info("API search returned no results, falling back to yt-dlp search", zap.String("query", query))
		return s.fallbackSearch(ctx, query)
	}
	return results, nil
}

func (s *ImportService) fallbackSearch(ctx context.Context, query string) ([]search.Result, error) {
	oldResults, err := s.downloader.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	converted := make([]search.Result, len(oldResults))
	for i, r := range oldResults {
		converted[i] = search.Result{
			Title:     r.Title,
			Artist:    r.Artist,
			URL:       r.URL,
			Duration:  r.Duration,
			Thumbnail: r.Thumbnail,
			Source:    r.Source,
		}
	}
	return converted, nil
}

func (s *ImportService) Import(ctx context.Context, url, source, title, artist string, externalIDs map[string]string, userID string) (*ImportResponse, error) {
	if s.importWorker == nil {
		return s.importSync(ctx, url, userID)
	}

	// If no URL, try to resolve one via the acquisition layer.
	// This handles metadata-only sources like MusicBrainz.
	downloadURL := url
	if downloadURL == "" && title != "" && artist != "" {
		resolveQuery := acquisition.ResolveQuery{
			Title:  title,
			Artist: artist,
		}
		// Map external IDs
		if isrc, ok := externalIDs["isrc"]; ok {
			resolveQuery.ISRC = isrc
		}
		if mbid, ok := externalIDs["mbid"]; ok {
			resolveQuery.ExternalIDs = map[string]string{"mbid": mbid}
		}

		resolveCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		candidate, err := s.resolver.Resolve(resolveCtx, resolveQuery)
		cancel()
		if err != nil {
			s.logger.Warn("failed to resolve download URL",
				zap.String("title", title),
				zap.String("artist", artist),
				zap.Error(err),
			)
			// Fall through — the worker can also try resolution
		} else if candidate != nil {
			downloadURL = candidate.URL
			source = candidate.Source
			s.logger.Info("resolved download URL",
				zap.String("title", title),
				zap.String("artist", artist),
				zap.String("url", downloadURL),
				zap.String("source", source),
			)
		}
	}

	// If we still have no URL and no worker, we can't proceed.
	if downloadURL == "" && s.importWorker == nil {
		return nil, fmt.Errorf("no downloadable URL available for %s - %s", artist, title)
	}

	jobID := uuid.New().String()
	job := &worker.Job{
		ID:     jobID,
		URL:    downloadURL,
		Source: source,
		Title:  title,
		Artist: artist,
		UserID: userID,
	}

	if s.importWorker != nil {
		if err := s.importWorker.Enqueue(ctx, job); err != nil {
			return nil, err
		}
		return &ImportResponse{
			JobID:   jobID,
			Message: "import job queued. check progress endpoint for status.",
		}, nil
	}

	// No worker — run sync
	return s.importSync(ctx, downloadURL, userID)
}

func (s *ImportService) importSync(ctx context.Context, url, userID string) (*ImportResponse, error) {
	if err := os.MkdirAll(s.downloadDir, 0755); err != nil {
		return nil, err
	}

	entry, localPath, err := s.downloader.Download(ctx, url, s.downloadDir)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(localPath)
	if err != nil {
		os.Remove(localPath)
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		os.Remove(localPath)
		return nil, err
	}

	header := &multipart.FileHeader{
		Filename: entry.Title + ".mp3",
		Size:     stat.Size(),
	}

	result, err := s.ingestionSvc.Upload(ctx, f, header, userID)
	f.Close()
	os.Remove(localPath)

	if err != nil {
		return nil, err
	}

	return &ImportResponse{
		DraftID:  result.DraftID,
		Title:    entry.Title,
		Artist:   entry.Uploader,
		Duration: int(entry.Duration),
		Message:  "track imported successfully. review and publish from the ingestion page.",
	}, nil
}

func (s *ImportService) GetProgress(ctx context.Context, jobID string) (*worker.ProgressUpdate, error) {
	return s.importWorker.GetProgress(ctx, jobID)
}

// ── Batch Import ─────────────────────────────────────────────────

func (s *ImportService) BatchImport(ctx context.Context, items []BatchImportItem, userID string) (*BatchImportResponse, error) {
	batchID := uuid.New().String()
	results := make([]BatchJobResult, 0, len(items))

	for _, item := range items {
		// For each track, construct external IDs
		extIDs := item.ExternalIDs
		if extIDs == nil {
			extIDs = make(map[string]string)
		}

		// If user provided a direct URL, pass it through; otherwise worker resolves
		downloadURL := item.URL

		// Enqueue directly without pre-resolution (worker handles it)
		resp, err := s.enqueueJob(ctx, item.Title, item.Artist, downloadURL, item.Source, extIDs, userID)
		result := BatchJobResult{
			Title:  item.Title,
			Artist: item.Artist,
		}
		if err != nil {
			result.Error = err.Error()
			s.logger.Warn("batch item import failed",
				zap.String("batch_id", batchID),
				zap.String("title", item.Title),
				zap.String("artist", item.Artist),
				zap.Error(err),
			)
		} else if resp != nil {
			result.JobID = resp.JobID
		}
		results = append(results, result)
	}

	// Store batch progress in context (simplified: in-memory for now)
	s.storeBatchProgress(batchID, results)

	return &BatchImportResponse{
		BatchID: batchID,
		Jobs:    results,
		Message: fmt.Sprintf("batch import initiated: %d tracks processed", len(items)),
	}, nil
}

// enqueueJob directly enqueues a job for the worker without pre-resolution.
// If url is empty, the worker will resolve the download URL during processing.
func (s *ImportService) enqueueJob(ctx context.Context, title, artist, url, source string, externalIDs map[string]string, userID string) (*ImportResponse, error) {
	jobID := uuid.New().String()
	job := &worker.Job{
		ID:     jobID,
		URL:    url, // If set, worker skips resolution
		Source: source,
		Title:  title,
		Artist: artist,
		UserID: userID,
	}

	if s.importWorker != nil {
		if err := s.importWorker.Enqueue(ctx, job); err != nil {
			return nil, err
		}
		return &ImportResponse{
			JobID:   jobID,
			Message: "import job queued. check progress endpoint for status.",
		}, nil
	}

	return nil, fmt.Errorf("no import worker available; cannot import '%s - %s'", artist, title)
}

// In-memory batch progress store (simple map; Redis-backed for production)
type batchProgressEntry struct {
	results []BatchJobResult
}

var (
	batchProgressMu    sync.Mutex
	batchProgressStore = make(map[string]*batchProgressEntry)
)

func (s *ImportService) storeBatchProgress(batchID string, results []BatchJobResult) {
	batchProgressMu.Lock()
	defer batchProgressMu.Unlock()
	batchProgressStore[batchID] = &batchProgressEntry{results: results}
}

func (s *ImportService) GetBatchProgress(ctx context.Context, batchID string) (*BatchProgressResponse, error) {
	batchProgressMu.Lock()
	entry, ok := batchProgressStore[batchID]
	batchProgressMu.Unlock()

	if !ok {
		return nil, nil
	}

	total := len(entry.results)
	completed := 0
	failed := 0
	inProgress := 0

	for _, r := range entry.results {
		if r.Error != "" {
			failed++
		} else if r.JobID != "" {
			// Check if this job is complete
			pu, err := s.GetProgress(ctx, r.JobID)
			if err != nil || pu == nil {
				inProgress++
			} else if pu.Status == worker.StatusComplete {
				completed++
			} else if pu.Status == worker.StatusFailed {
				failed++
			} else {
				inProgress++
			}
		} else {
			inProgress++
		}
	}

	progressPct := 0
	if total > 0 {
		progressPct = (completed + failed) * 100 / total
	}

	return &BatchProgressResponse{
		BatchID:     batchID,
		Total:       total,
		Completed:   completed,
		Failed:      failed,
		InProgress:  inProgress,
		ProgressPct: progressPct,
	}, nil
}
