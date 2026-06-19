package importcmd

import (
	"context"
	"mime/multipart"
	"os"

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

func (s *ImportService) Import(ctx context.Context, url, source, userID string) (*ImportResponse, error) {
	if s.importWorker == nil {
		return s.importSync(ctx, url, userID)
	}

	jobID := uuid.New().String()
	job := &worker.Job{
		ID:     jobID,
		URL:    url,
		Source: source,
		UserID: userID,
	}

	if err := s.importWorker.Enqueue(ctx, job); err != nil {
		return nil, err
	}

	return &ImportResponse{
		JobID:   jobID,
		Message: "import job queued. check progress endpoint for status.",
	}, nil
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
