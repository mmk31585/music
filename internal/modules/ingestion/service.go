package ingestion

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"music/internal/modules/ingestion/enrichment"
	platformstorage "music/internal/platform/storage"
)

const (
	MaxUploadSize       = 200 << 20
	StorageAudioDir     = "ingestion-audio"
	StorageCoverDir     = "ingestion-covers"
	StorageEnrichDir    = "enrichment-images"
)

type Service struct {
	storage   platformstorage.Storage
	repo      *Repository
	enricher  *enrichment.Enricher
}

func NewService(storage platformstorage.Storage, repo *Repository, enricher *enrichment.Enricher) *Service {
	return &Service{
		storage:  storage,
		repo:     repo,
		enricher: enricher,
	}
}

func (s *Service) EnrichDraft(ctx context.Context, draftID string) error {
	draft, err := s.repo.GetDraftByID(ctx, draftID)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}

	if draft.Status != DraftStatusPending && draft.Status != DraftStatusEnrichmentFailed {
		return nil
	}

	if err := s.repo.UpdateDraftStatus(ctx, draftID, DraftStatusEnriching); err != nil {
		return err
	}

	tags := parseExtractedTags(draft.ExtractedMetadata)

	go s.enrichAsync(draftID, tags.Title, tags.Artist, tags.Album)

	return nil
}

func (s *Service) enrichAsync(draftID, title, artist, album string) {
	ctx := context.Background()

	result, err := s.enricher.Enrich(ctx, title, artist, album)
	if err != nil {
		zap.L().Error("enrichment failed", zap.String("draft_id", draftID), zap.Error(err))
		_ = s.repo.UpdateDraftStatus(ctx, draftID, DraftStatusEnrichmentFailed)
		return
	}

	if result.Spotify != nil {
		if result.Spotify.AlbumCoverURL != "" {
			if localURL, err := s.downloadAndStoreImage(ctx, result.Spotify.AlbumCoverURL, draftID, "album-cover"); err != nil {
				zap.L().Warn("failed to download album cover", zap.String("draft_id", draftID), zap.Error(err))
			} else {
				result.Spotify.AlbumCoverURL = localURL
				for i, sug := range result.Suggestions {
					if sug.Field == "album_cover_url" {
						result.Suggestions[i].Value = localURL
						break
					}
				}
			}
		}
		if result.Spotify.ArtistImageURL != "" {
			if localURL, err := s.downloadAndStoreImage(ctx, result.Spotify.ArtistImageURL, draftID, "artist-image"); err != nil {
				zap.L().Warn("failed to download artist image", zap.String("draft_id", draftID), zap.Error(err))
			} else {
				result.Spotify.ArtistImageURL = localURL
				for i, sug := range result.Suggestions {
					if sug.Field == "artist_image_url" {
						result.Suggestions[i].Value = localURL
						break
					}
				}
			}
		}
	}

	enrichedJSON, err := enrichment.SerializeResult(result)
	if err != nil {
		zap.L().Error("failed to serialize enrichment", zap.String("draft_id", draftID), zap.Error(err))
		_ = s.repo.UpdateDraftStatus(ctx, draftID, DraftStatusEnrichmentFailed)
		return
	}

	newStatus := DraftStatusReview

	if err := s.repo.UpdateDraftEnrichedMetadata(ctx, draftID, enrichedJSON, newStatus); err != nil {
		zap.L().Error("failed to save enrichment", zap.String("draft_id", draftID), zap.Error(err))
		_ = s.repo.UpdateDraftStatus(ctx, draftID, DraftStatusEnrichmentFailed)
		return
	}
}

func (s *Service) downloadAndStoreImage(ctx context.Context, imageURL, draftID, suffix string) (string, error) {
	parsed, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	ext := filepath.Ext(parsed.Path)
	if ext == "" {
		ext = ".jpg"
	}
	key := fmt.Sprintf("%s/%s-%s%s", StorageEnrichDir, draftID, suffix, ext)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	if err := s.storage.Upload(ctx, key, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		return "", fmt.Errorf("store: %w", err)
	}

	localURL, err := s.storage.GetURL(ctx, key)
	if err != nil {
		return "", fmt.Errorf("get url: %w", err)
	}

	return localURL, nil
}

func (s *Service) GetDraftSuggestions(ctx context.Context, id string) (*enrichment.EnrichmentResult, error) {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, ErrDraftNotFound
	}

	if draft.EnrichedMetadata == nil || *draft.EnrichedMetadata == "" {
		return &enrichment.EnrichmentResult{Suggestions: []enrichment.EnrichedSuggestion{}, Attempted: false}, nil
	}

	return enrichment.DeserializeResult(*draft.EnrichedMetadata)
}

func (s *Service) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploadedBy string) (*UploadResponse, error) {
	if file == nil || header == nil {
		return nil, ErrNoFileProvided
	}

	if header.Size > MaxUploadSize {
		return nil, ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	format := formatFromExt(ext)
	if format == "" {
		return nil, ErrInvalidFileType
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if len(content) == 0 {
		return nil, ErrNoFileProvided
	}

	hash := sha256.Sum256(content)
	fileHash := hex.EncodeToString(hash[:])

	existing, err := s.repo.SearchByFileHash(ctx, fileHash)
	if err != nil {
		zap.L().Warn("failed to check duplicate hash", zap.Error(err))
	} else if len(existing) > 0 {
		for _, ex := range existing {
			if ex.Status == DraftStatusPublished {
				return nil, ErrDuplicateFile
			}
		}
		zap.L().Warn("potential duplicate file (draft already exists)",
			zap.String("existing_draft", existing[0].ID),
		)
	}

	draftID := ulid.Make().String()

	audioKey := fmt.Sprintf("%s/%s-%s", StorageAudioDir, draftID, sanitizeFilename(header.Filename))

	if err := s.storage.Upload(ctx, audioKey, bytes.NewReader(content), int64(len(content)), detectMimeByExt(ext)); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	audioURL, err := s.storage.GetURL(ctx, audioKey)
	if err != nil {
		_ = s.storage.Delete(ctx, audioKey)
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	result, err := extractMetadata(header.Filename, bytes.NewReader(content))
	if err != nil {
		zap.L().Warn("metadata extraction failed", zap.String("filename", header.Filename), zap.Error(err))
		result = &extractionResult{
			Tags: &ExtractedTags{Format: format},
		}
	}

	if result.Audio != nil {
		if result.Audio.Format != "" {
			format = result.Audio.Format
		}
	}

	extractedJSON, err := extractedTagsToJSON(result.Tags)
	if err != nil {
		extractedJSON = "{}"
	}

	var duration *float64
	var bitrate *int
	if result.Audio != nil {
		duration = &result.Audio.Duration
		bitrate = &result.Audio.Bitrate
	}

	draft, err := s.repo.CreateDraft(ctx, CreateDraftParams{
		ID:                draftID,
		UploadedBy:        uploadedBy,
		OriginalFilename:  header.Filename,
		FilePath:          audioKey,
		FileSize:          header.Size,
		DurationSeconds:   duration,
		Bitrate:           bitrate,
		Format:            format,
		FileHash:          fileHash,
		ExtractedMetadata: extractedJSON,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, audioKey)
		return nil, err
	}

	assets := []AssetResponse{}
	var coverArtURL *string

	assets = append(assets, AssetResponse{
		ID:        ulid.Make().String(),
		AssetType: "audio",
		URL:       audioURL,
		Source:    "uploaded",
	})

	if len(result.CoverArt) > 0 {
		coverExt := result.CoverExt
		if coverExt == "" {
			coverExt = ".jpg"
		}
		coverKey := fmt.Sprintf("%s/%s-cover%s", StorageCoverDir, draftID, coverExt)
		coverMime := detectCoverMime(coverExt)

		if err := s.storage.Upload(ctx, coverKey, bytes.NewReader(result.CoverArt), int64(len(result.CoverArt)), coverMime); err != nil {
			zap.L().Warn("failed to upload embedded cover art", zap.Error(err))
		} else {
			coverURL, err := s.storage.GetURL(ctx, coverKey)
			if err == nil {
				coverArtURL = &coverURL
				assetID := ulid.Make().String()
				_, err := s.repo.CreateAsset(ctx, CreateAssetParams{
					ID:        assetID,
					DraftID:   draftID,
					AssetType: "cover",
					URL:       coverURL,
					Source:    "embedded",
				})
				if err == nil {
					assets = append(assets, AssetResponse{
						ID:        assetID,
						AssetType: "cover",
						URL:       coverURL,
						Source:    "embedded",
					})
				}
			}
		}
	}

	go s.enrichAsync(draftID, result.Tags.Title, result.Tags.Artist, result.Tags.Album)

	return &UploadResponse{
		DraftID:           draft.ID,
		OriginalFilename:  draft.OriginalFilename,
		FileSize:          draft.FileSize,
		Format:            draft.Format,
		DurationSeconds:   draft.DurationSeconds,
		Bitrate:           draft.Bitrate,
		Status:            DraftStatusEnriching,
		ExtractedMetadata: result.Tags,
		CoverArtURL:       coverArtURL,
		Assets:            assets,
		CreatedAt:         draft.CreatedAt,
	}, nil
}

func (s *Service) GetDraftByID(ctx context.Context, id string) (*DraftDetailResponse, error) {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, ErrDraftNotFound
	}

	assets, err := s.repo.GetAssetsByDraftID(ctx, id)
	if err != nil {
		return nil, err
	}

	assetResponses := make([]AssetResponse, 0, len(assets))
	for _, a := range assets {
		assetResponses = append(assetResponses, AssetResponse{
			ID:        a.ID,
			AssetType: a.AssetType,
			URL:       a.URL,
			Source:    a.Source,
			CreatedAt: a.CreatedAt,
		})
	}

	tags := parseExtractedTags(draft.ExtractedMetadata)

	resp := &DraftDetailResponse{
		ID:                draft.ID,
		UploadedBy:        draft.UploadedBy,
		OriginalFilename:  draft.OriginalFilename,
		FileSize:          draft.FileSize,
		Format:            draft.Format,
		DurationSeconds:   draft.DurationSeconds,
		Bitrate:           draft.Bitrate,
		Status:            draft.Status,
		ExtractedMetadata: tags,
		Assets:            assetResponses,
		CreatedAt:         draft.CreatedAt,
		UpdatedAt:         draft.UpdatedAt,
	}

	if draft.EnrichedMetadata != nil {
		resp.EnrichedMetadata = jsonToMap(*draft.EnrichedMetadata)
	}
	if draft.FinalMetadata != nil {
		resp.FinalMetadata = jsonToMap(*draft.FinalMetadata)
	}

	return resp, nil
}

func (s *Service) ListDrafts(ctx context.Context, status string, page, limit int) (*ListDraftsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if status != "" {
		valid := false
		for _, s := range []DraftStatus{DraftStatusPending, DraftStatusEnriching, DraftStatusReview, DraftStatusAccepted, DraftStatusRejected, DraftStatusPublished, DraftStatusEnrichmentFailed} {
			if DraftStatus(status) == s {
				valid = true
				break
			}
		}
		if !valid {
			return nil, ErrInvalidStatus
		}
	}

	rows, total, err := s.repo.ListDrafts(ctx, ListDraftsParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]DraftListItem, 0, len(rows))
	for _, row := range rows {
		tags := parseExtractedTags(row.ExtractedMetadata)
		item := DraftListItem{
			ID:               row.ID,
			OriginalFilename: row.OriginalFilename,
			FileSize:         row.FileSize,
			Format:           row.Format,
			DurationSeconds:  row.DurationSeconds,
			Status:           row.Status,
			FileHash:         row.FileHash,
			Stale:            row.Stale,
			Title:            tags.Title,
			Artist:           tags.Artist,
			Album:            tags.Album,
			CoverArtURL:      row.CoverArtURL,
			HasCoverArt:      tags.HasCoverArt,
			CreatedAt:        row.CreatedAt,
		}
		items = append(items, item)
	}

	return &ListDraftsResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	hash := sha256.Sum256([]byte(name + uuid.New().String()))
	return hex.EncodeToString(hash[:8]) + "-" + strings.ReplaceAll(strings.ToLower(filepath.Base(name)), " ", "-")
}

func detectMimeByExt(ext string) string {
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".flac":
		return "audio/flac"
	case ".ogg":
		return "audio/ogg"
	case ".m4a", ".aac":
		return "audio/mp4"
	default:
		return "application/octet-stream"
	}
}

func detectCoverMime(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}

func (s *Service) SaveFinalMetadata(ctx context.Context, id string, req SaveFinalMetadataRequest) error {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}
	if draft.Status != DraftStatusReview && draft.Status != DraftStatusEnrichmentFailed {
		return ErrInvalidStatus
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal final metadata: %w", err)
	}

	return s.repo.UpdateDraftFinalMetadata(ctx, id, string(data), DraftStatusAccepted)
}

func (s *Service) RejectDraft(ctx context.Context, id string, reason string) error {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}

	rejectMeta := map[string]string{"reason": reason}
	data, err := json.Marshal(rejectMeta)
	if err != nil {
		return fmt.Errorf("marshal reject metadata: %w", err)
	}

	return s.repo.UpdateDraftFinalMetadata(ctx, id, string(data), DraftStatusRejected)
}

func (s *Service) SearchArtists(ctx context.Context, q string) ([]ArtistSearchResult, error) {
	if q == "" {
		return []ArtistSearchResult{}, nil
	}
	const limit = 20
	return s.repo.SearchArtists(ctx, q, limit)
}

func (s *Service) SearchAlbums(ctx context.Context, q string) ([]AlbumSearchResult, error) {
	if q == "" {
		return []AlbumSearchResult{}, nil
	}
	const limit = 20
	return s.repo.SearchAlbums(ctx, q, limit)
}

func (s *Service) UploadDraftImage(ctx context.Context, draftID, entityType string, file multipart.File, ext string) (string, error) {
	draft, err := s.repo.GetDraftByID(ctx, draftID)
	if err != nil {
		return "", err
	}
	if draft == nil {
		return "", ErrDraftNotFound
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	imageKey := fmt.Sprintf("%s/%s-%s-image%s", StorageEnrichDir, draftID, entityType, ext)
	mime := detectImageMime(ext)

	if err := s.storage.Upload(ctx, imageKey, bytes.NewReader(content), int64(len(content)), mime); err != nil {
		return "", fmt.Errorf("store image: %w", err)
	}

	imageURL, err := s.storage.GetURL(ctx, imageKey)
	if err != nil {
		_ = s.storage.Delete(ctx, imageKey)
		return "", fmt.Errorf("get url: %w", err)
	}

	return imageURL, nil
}

func detectImageMime(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}

func (s *Service) GetDraftRaw(ctx context.Context, id string) (*IngestionDraft, error) {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, ErrDraftNotFound
	}
	return draft, nil
}

func (s *Service) GetIngestionStats(ctx context.Context) (*IngestionStats, error) {
	published, err := s.repo.CountPublishedThisMonth(ctx)
	if err != nil {
		return nil, err
	}
	byStatus, err := s.repo.CountDraftsByStatus(ctx)
	if err != nil {
		return nil, err
	}
	total := 0
	for _, c := range byStatus {
		total += c
	}
	return &IngestionStats{
		PublishedThisMonth: published,
		PendingReview:      byStatus[DraftStatusReview],
		TotalDrafts:        total,
		ByStatus:           map[string]int{"pending": byStatus[DraftStatusPending], "enriching": byStatus[DraftStatusEnriching], "review": byStatus[DraftStatusReview], "accepted": byStatus[DraftStatusAccepted], "rejected": byStatus[DraftStatusRejected], "published": byStatus[DraftStatusPublished], "enrichment_failed": byStatus[DraftStatusEnrichmentFailed]},
	}, nil
}

func jsonToMap(rawJSON string) map[string]interface{} {
	data := map[string]interface{}{}
	if rawJSON == "" || rawJSON == "{}" {
		return data
	}
	if err := json.Unmarshal([]byte(rawJSON), &data); err != nil {
		return map[string]interface{}{}
	}
	return data
}


