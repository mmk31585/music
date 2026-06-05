package media

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"

	platformstorage "music/internal/platform/storage"
)

type Config struct {
	MaxFileSizeBytes int64

	MaxImageSizeBytes int64
	MaxAudioSizeBytes int64

	AllowedImageMime []string
	AllowedAudioMime []string
}

type Service struct {
	storage platformstorage.Storage
	repo    *Repository
	cfg     Config
}

func NewService(storage platformstorage.Storage, repo *Repository, cfg Config) *Service {
	return &Service{
		storage: storage,
		repo:    repo,
		cfg:     cfg,
	}
}

type UploadCategory string

const (
	UploadCategoryArtistImage UploadCategory = "artist-images"
	UploadCategoryAlbumCover  UploadCategory = "album-covers"
	UploadCategoryTrackCover  UploadCategory = "track-covers"
	UploadCategoryTrackAudio  UploadCategory = "track-audio"
)

func (s *Service) Upload(
	ctx context.Context,
	category UploadCategory,
	file multipart.File,
	header *multipart.FileHeader,
) (*UploadResponse, error) {
	if file == nil || header == nil {
		return nil, ErrNoFileProvided
	}

	if !s.isKnownCategory(category) {
		return nil, ErrInvalidFieldName
	}

	maxSize := s.maxSizeForCategory(category)
	if maxSize > 0 && header.Size > maxSize {
		return nil, ErrFileTooLarge
	}

	log.Println("CATEGORY:", category)
	log.Println("FILE:", header.Filename)
	log.Println("SIZE:", header.Size)

	mimeType, err := detectMimeType(file)
	if err != nil {
		log.Println("MIME ERROR:", err)
		return nil, ErrInvalidMimeType
	}

	log.Println("MIME:", mimeType)

	if !s.isAllowedMime(category, mimeType) {
		return nil, ErrInvalidMimeType
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
		}
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	if len(content) == 0 {
		return nil, ErrEmptyFile
	}

	hashBytes := sha256.Sum256(content)
	hash := hex.EncodeToString(hashBytes[:])

	existing, err := s.repo.FindByChecksum(ctx, hash)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return uploadResponseFromMedia(existing, true), nil
	}

	ext := extensionForMime(mimeType)
	if ext == "" {
		ext = sanitizeExtension(strings.ToLower(filepath.Ext(header.Filename)))
	}
	if ext == "" {
		ext = ".bin"
	}

	original := sanitizeFilename(header.Filename)
	key := buildObjectKey(category, original)

	err = s.storage.Upload(
		ctx,
		key,
		bytes.NewReader(content),
		int64(len(content)),
		mimeType,
	)
	if err != nil {
		log.Printf("[UPLOAD STORAGE FAIL] key=%s size=%d mime=%s err=%+v",
			key, len(content), mimeType, err,
		)
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	url, err := s.storage.GetURL(ctx, key)
	if err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	size := int64(len(content))
	metadata := buildUploadMetadata(category, header.Filename)

	mediaItem, err := s.repo.Create(ctx, CreateMediaRequest{
		MediaType:        mediaTypeForCategory(category),
		StorageProvider:  "local",
		Bucket:           nil,
		ObjectKey:        key,
		PublicURL:        &url,
		MimeType:         &mimeType,
		FileSize:         &size,
		ChecksumSHA256:   &hash,
		DurationSeconds:  nil,
		Width:            nil,
		Height:           nil,
		OriginalFilename: &header.Filename,
		Metadata:         metadata,
		CreatedBy:        nil,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, err
	}

	resp := uploadResponseFromMedia(mediaItem, false)
	resp.FileName = original

	return resp, nil
}
func sanitizeFilename(name string) string {
	name = strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9.\-_]+`)
	name = re.ReplaceAllString(name, "-")

	return name
}
func uploadResponseFromMedia(item *Media, duplicate bool) *UploadResponse {
	if item == nil {
		return nil
	}

	url := ""
	if item.PublicURL != nil {
		url = *item.PublicURL
	}

	originalName := ""
	if item.OriginalFilename != nil {
		originalName = *item.OriginalFilename
	}

	size := int64(0)
	if item.FileSize != nil {
		size = *item.FileSize
	}

	mimeType := ""
	if item.MimeType != nil {
		mimeType = *item.MimeType
	}

	sha256Value := ""
	if item.ChecksumSHA256 != nil {
		sha256Value = *item.ChecksumSHA256
	}

	return &UploadResponse{
		MediaID:         item.ID.String(),
		URL:             url,
		Path:            item.ObjectKey,
		FileName:        filepath.Base(item.ObjectKey),
		OriginalName:    originalName,
		Size:            size,
		MimeType:        mimeType,
		SHA256:          sha256Value,
		Duplicate:       duplicate,
		DurationSeconds: item.DurationSeconds,
	}
}

func (s *Service) DeleteUploadedFile(ctx context.Context, upload *UploadResponse) error {
	if upload == nil || upload.Path == "" || upload.Duplicate {
		return nil
	}

	return s.storage.Delete(ctx, upload.Path)
}

func buildObjectKey(category UploadCategory, fileName string) string {
	return fmt.Sprintf("%s/%s-%s", sanitizeCategory(string(category)), uuid.NewString(), fileName)
}

func mediaTypeForCategory(category UploadCategory) string {
	switch category {
	case UploadCategoryArtistImage, UploadCategoryAlbumCover, UploadCategoryTrackCover:
		return "image"
	case UploadCategoryTrackAudio:
		return "audio"
	default:
		return "other"
	}
}

func buildUploadMetadata(category UploadCategory, originalFilename string) string {
	return fmt.Sprintf(
		`{"uploadCategory":%q,"originalFilename":%q}`,
		string(category),
		originalFilename,
	)
}

func (s *Service) isKnownCategory(category UploadCategory) bool {
	switch category {
	case UploadCategoryArtistImage,
		UploadCategoryAlbumCover,
		UploadCategoryTrackCover,
		UploadCategoryTrackAudio:
		return true
	default:
		return false
	}
}

func (s *Service) maxSizeForCategory(category UploadCategory) int64 {
	switch category {
	case UploadCategoryArtistImage, UploadCategoryAlbumCover, UploadCategoryTrackCover:
		if s.cfg.MaxImageSizeBytes > 0 {
			return s.cfg.MaxImageSizeBytes
		}
	case UploadCategoryTrackAudio:
		if s.cfg.MaxAudioSizeBytes > 0 {
			return s.cfg.MaxAudioSizeBytes
		}
	}

	return s.cfg.MaxFileSizeBytes
}

func (s *Service) isAllowedMime(category UploadCategory, mimeType string) bool {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	if mimeType == "" {
		return false
	}

	switch category {
	case UploadCategoryArtistImage, UploadCategoryAlbumCover, UploadCategoryTrackCover:
		return containsMime(s.cfg.AllowedImageMime, mimeType)
	case UploadCategoryTrackAudio:
		return containsMime(s.cfg.AllowedAudioMime, mimeType)
	default:
		return false
	}
}

func containsMime(list []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))

	for _, v := range list {
		v = strings.ToLower(strings.TrimSpace(v))
		if v == value {
			return true
		}
	}

	return false
}

func detectMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return "", err
		}
	}

	return http.DetectContentType(buffer[:n]), nil
}
