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
	MaxUploadSize    = 200 << 20
	StorageAudioDir  = "ingestion-audio"
	StorageCoverDir  = "track-covers"
	StorageEnrichDir = "enrichment-images"
)

type Service struct {
	storage  platformstorage.Storage
	repo     *Repository
	enricher *enrichment.Enricher
}

func NewService(storage platformstorage.Storage, repo *Repository, enricher *enrichment.Enricher) *Service {
	return &Service{
		storage:  storage,
		repo:     repo,
		enricher: enricher,
	}
}

func (s *Service) EnrichDraft(ctx context.Context, draftID string, scope ...enrichment.EnrichScope) error {
	draft, err := s.repo.GetDraftByID(ctx, draftID)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}

	// Allow enrichment from any non-terminal state
	if draft.Status != DraftStatusPending && draft.Status != DraftStatusEnrichmentFailed && draft.Status != DraftStatusReview {
		return nil
	}

	// Only set status to enriching for full enrichment (no scope) or if coming from pending/failed
	if len(scope) == 0 || len(scope[0]) == 0 || draft.Status != DraftStatusReview {
		if err := s.repo.UpdateDraftStatus(ctx, draftID, DraftStatusEnriching); err != nil {
			return err
		}
	}

	tags := parseExtractedTags(draft.ExtractedMetadata)

	go s.enrichAsync(draftID, tags.Title, tags.Artist, tags.Album, int(tags.Duration), scope...)

	return nil
}

func (s *Service) enrichAsync(draftID, title, artist, album string, durationSeconds int, scope ...enrichment.EnrichScope) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("enrichment panicked", zap.String("draft_id", draftID), zap.Any("recover", r))
		}
	}()

	ctx := context.Background()

	result, err := s.enricher.Enrich(ctx, title, artist, album, durationSeconds, scope...)
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
				s.createAsset(ctx, draftID, "album_cover", localURL, "enrichment")
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
				s.createAsset(ctx, draftID, "artist_image", localURL, "enrichment")
			}
		}
	}

	if result.Spotify == nil || result.Spotify.ArtistImageURL == "" {
		if result.LastFM != nil && result.LastFM.ArtistImageURL != "" {
			if localURL, err := s.downloadAndStoreImage(ctx, result.LastFM.ArtistImageURL, draftID, "artist-image-lfm"); err != nil {
				zap.L().Warn("failed to download last.fm artist image", zap.String("draft_id", draftID), zap.Error(err))
			} else {
				result.LastFM.ArtistImageURL = localURL
				for i, sug := range result.Suggestions {
					if sug.Field == "artist_image_url" {
						result.Suggestions[i].Value = localURL
						break
					}
				}
				s.createAsset(ctx, draftID, "artist_image", localURL, "enrichment")
			}
		}
	}

	if (result.Spotify == nil || result.Spotify.AlbumCoverURL == "") && result.LastFM != nil && result.LastFM.AlbumCoverURL != "" {
		if localURL, err := s.downloadAndStoreImage(ctx, result.LastFM.AlbumCoverURL, draftID, "album-cover-lfm"); err != nil {
			zap.L().Warn("failed to download last.fm album cover", zap.String("draft_id", draftID), zap.Error(err))
		} else {
			result.LastFM.AlbumCoverURL = localURL
			for i, sug := range result.Suggestions {
				if sug.Field == "album_cover_url" {
					result.Suggestions[i].Value = localURL
					break
				}
			}
			s.createAsset(ctx, draftID, "album_cover", localURL, "enrichment")
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

func (s *Service) createAsset(ctx context.Context, draftID, assetType, url, source string) {
	_, err := s.repo.CreateAsset(ctx, CreateAssetParams{
		ID:        uuid.New().String(),
		DraftID:   draftID,
		AssetType: assetType,
		URL:       url,
		Source:    source,
	})
	if err != nil {
		zap.L().Warn("failed to create asset", zap.String("draft_id", draftID), zap.String("asset_type", assetType), zap.Error(err))
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
			zap.L().Warn("failed to store embedded cover art",
				zap.String("draft_id", draftID),
				zap.Int("cover_size", len(result.CoverArt)),
				zap.String("cover_ext", coverExt),
				zap.Error(err),
			)
		} else {
			coverURL, err := s.storage.GetURL(ctx, coverKey)
			if err != nil {
				zap.L().Warn("failed to get cover art URL",
					zap.String("draft_id", draftID),
					zap.String("cover_key", coverKey),
					zap.Error(err),
				)
			} else {
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

	go s.enrichAsync(draftID, result.Tags.Title, result.Tags.Artist, result.Tags.Album, int(result.Tags.Duration))

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

// PrepareUpload validates the requested upload and returns a presigned PUT URL
// for direct client-to-S3 upload. It does not create a draft yet.
func (s *Service) PrepareUpload(ctx context.Context, filename, contentType string, size int64, uploadedBy string) (*PrepareUploadResponse, error) {
	if filename == "" {
		return nil, ErrNoFileProvided
	}

	if size > MaxUploadSize {
		return nil, ErrFileTooLarge
	}

	if size == 0 {
		return nil, ErrNoFileProvided
	}

	ext := strings.ToLower(filepath.Ext(filename))
	format := formatFromExt(ext)
	if format == "" {
		return nil, ErrInvalidFileType
	}

	draftID := ulid.Make().String()
	key := fmt.Sprintf("%s/%s-%s", StorageAudioDir, draftID, sanitizeFilename(filename))

	presignTTL := 10 * time.Minute
	uploadURL, err := s.storage.PresignPut(ctx, key, contentType, presignTTL)
	if err != nil {
		zap.L().Error("presign put failed", zap.String("step", uploadStepPresign), zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	zap.L().Info("presigned upload url created",
		zap.String("step", uploadStepPresign),
		zap.String("draft_id", draftID),
		zap.String("key", key),
		zap.String("filename", filename),
		zap.Int64("size", size),
		zap.String("user_id", uploadedBy),
	)

	return &PrepareUploadResponse{
		DraftID:   draftID,
		Key:       key,
		UploadURL: uploadURL,
		ExpiresIn: int64(presignTTL.Seconds()),
	}, nil
}

// CompleteUpload is called after the client PUTs the file to S3 via the presigned URL.
// It verifies the object exists, downloads it, extracts metadata, creates the draft,
// and returns the upload response.
func (s *Service) CompleteUpload(ctx context.Context, draftID, key, uploadedBy string) (*UploadResponse, error) {
	if draftID == "" || key == "" {
		return nil, ErrNoFileProvided
	}

	// Verify the object exists in S3
	exists, err := s.storage.Exists(ctx, key)
	if err != nil {
		zap.L().Error("s3 head object failed", zap.String("step", "s3_head_object"), zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}
	if !exists {
		zap.L().Error("s3 object not found", zap.String("step", "s3_head_object"), zap.String("key", key))
		return nil, ErrNoFileProvided
	}

	// Download the object for hash + metadata extraction
	content, err := s.storage.Get(ctx, key)
	if err != nil {
		zap.L().Error("s3 get object failed", zap.String("step", "s3_get_object"), zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	if len(content) == 0 {
		return nil, ErrNoFileProvided
	}

	// Extract filename from key (format: ingestion-audio/<draftID>-<filename>)
	filename := key
	if parts := strings.SplitN(key, "-", 2); len(parts) == 2 {
		filename = parts[1]
	}

	ext := strings.ToLower(filepath.Ext(filename))
	format := formatFromExt(ext)
	if format == "" {
		return nil, ErrInvalidFileType
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

	result, err := extractMetadata(filename, bytes.NewReader(content))
	if err != nil {
		zap.L().Warn("metadata extraction failed", zap.String("filename", filename), zap.Error(err))
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

	audioURL, err := s.storage.GetURL(ctx, key)
	if err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	draft, err := s.repo.CreateDraft(ctx, CreateDraftParams{
		ID:                draftID,
		UploadedBy:        uploadedBy,
		OriginalFilename:  filename,
		FilePath:          key,
		FileSize:          int64(len(content)),
		DurationSeconds:   duration,
		Bitrate:           bitrate,
		Format:            format,
		FileHash:          fileHash,
		ExtractedMetadata: extractedJSON,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, key)
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
			zap.L().Warn("failed to store embedded cover art",
				zap.String("draft_id", draftID),
				zap.Int("cover_size", len(result.CoverArt)),
				zap.String("cover_ext", coverExt),
				zap.Error(err),
			)
		} else {
			coverURL, err := s.storage.GetURL(ctx, coverKey)
			if err != nil {
				zap.L().Warn("failed to get cover art URL",
					zap.String("draft_id", draftID),
					zap.String("cover_key", coverKey),
					zap.Error(err),
				)
			} else {
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

	go s.enrichAsync(draftID, result.Tags.Title, result.Tags.Artist, result.Tags.Album, int(result.Tags.Duration))

	zap.L().Info("upload complete (direct S3)",
		zap.String("step", uploadStepDBCreate),
		zap.String("draft_id", draftID),
		zap.String("user_id", uploadedBy),
	)

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

// UpdateExtractedMetadata overrides the draft's extracted (ffprobe) metadata
// with provided values. This is used by the import worker when the downloaded
// file has incorrect embedded tags (e.g., iTunes previews with generic IDs).
func (s *Service) UpdateExtractedMetadata(ctx context.Context, draftID, title, artist, album, coverURL string) error {
	draft, err := s.repo.GetDraftByID(ctx, draftID)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}

	tags := parseExtractedTags(draft.ExtractedMetadata)
	if title != "" {
		tags.Title = title
	}
	if artist != "" {
		tags.Artist = artist
	}
	if album != "" {
		tags.Album = album
	}
	if coverURL != "" {
		tags.CoverURL = coverURL
	}

	extractedJSON, err := extractedTagsToJSON(tags)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	return s.repo.UpdateDraftExtractedMetadata(ctx, draftID, extractedJSON)
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

	imageKey1 := fmt.Sprintf("%s/%s-%s-image%s", StorageEnrichDir, draftID, entityType, ext)
	imageKey2 := fmt.Sprintf("track-covers/%s-cover%s", draftID, ext)
	mime := detectImageMime(ext)

	// Upload to BOTH key formats for backward compatibility
	err1 := s.storage.Upload(ctx, imageKey1, bytes.NewReader(content), int64(len(content)), mime)
	err2 := s.storage.Upload(ctx, imageKey2, bytes.NewReader(content), int64(len(content)), mime)

	if err1 != nil && err2 != nil {
		return "", fmt.Errorf("store image: %w", err1)
	}

	// Return the new format URL (track-covers/)
	imageURL, err := s.storage.GetURL(ctx, imageKey2)
	if err != nil {
		// Fallback to the original format
		imageURL, err = s.storage.GetURL(ctx, imageKey1)
		if err != nil {
			_ = s.storage.Delete(ctx, imageKey1)
			_ = s.storage.Delete(ctx, imageKey2)
			return "", fmt.Errorf("get url: %w", err)
		}
	}

	// Save cover URL to draft record
	if err := s.repo.UpdateDraftCoverURL(ctx, draftID, imageURL); err != nil {
		zap.L().Warn("failed to save cover URL to draft", zap.String("draft_id", draftID), zap.Error(err))
	}

	// Clean up the original format key if new format succeeded
	_ = s.storage.Delete(ctx, imageKey1)

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

func (s *Service) DeleteDraft(ctx context.Context, id, requestingUserID string) error {
	draft, err := s.repo.GetDraftByID(ctx, id)
	if err != nil {
		return err
	}
	if draft == nil {
		return ErrDraftNotFound
	}

	// Only the uploader or an admin can delete a draft.
	// Admin check is enforced at the route level (RequireRole("admin")).
	if draft.UploadedBy != requestingUserID {
		return ErrUnauthorized
	}

	// Clean up storage files associated with the draft.
	if draft.FilePath != "" {
		if err := s.storage.Delete(ctx, draft.FilePath); err != nil {
			zap.L().Warn("failed to delete draft audio file from storage",
				zap.String("draft_id", id),
				zap.String("file_path", draft.FilePath),
				zap.Error(err),
			)
		}
	}

	// Clean up assets (cover art, images, etc.)
	assets, err := s.repo.GetAssetsByDraftID(ctx, id)
	if err != nil {
		zap.L().Warn("failed to fetch assets for cleanup", zap.String("draft_id", id), zap.Error(err))
	} else {
		for _, a := range assets {
			if a.URL != "" {
				// Extract storage key from URL if needed
				if key, ok := extractStorageKey(a.URL); ok {
					if err := s.storage.Delete(ctx, key); err != nil {
						zap.L().Warn("failed to delete asset from storage",
							zap.String("draft_id", id),
							zap.String("asset_url", a.URL),
							zap.Error(err),
						)
					}
				}
			}
		}
	}

	return s.repo.DeleteDraft(ctx, id)
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

// extractStorageKey attempts to extract a storage key from a stored URL.
// Returns the key and true if successful, or ("", false) if the URL doesn't
// match the local storage pattern (meaning it's an external URL, not stored).
//
// This is a best-effort helper for cleanup — it only handles local storage URLs
// that contain the known storage directory prefixes.
func extractStorageKey(url string) (string, bool) {
	for _, prefix := range []string{StorageAudioDir, StorageCoverDir, StorageEnrichDir} {
		if idx := strings.Index(url, prefix); idx != -1 {
			return url[idx:], true
		}
	}
	return "", false
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
