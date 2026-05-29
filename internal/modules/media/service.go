package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Config struct {
	MaxFileSizeBytes int64

	MaxImageSizeBytes int64
	MaxAudioSizeBytes int64

	AllowedImageMime []string
	AllowedAudioMime []string
}

type Service struct {
	storage Storage
	cfg     Config
}

func NewService(storage Storage, cfg Config) *Service {
	return &Service{
		storage: storage,
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

	mimeType, err := detectMimeType(file)
	if err != nil {
		return nil, ErrInvalidMimeType
	}

	if !s.isAllowedMime(category, mimeType) {
		return nil, ErrInvalidMimeType
	}

	stored, err := s.storage.SaveFile(ctx, file, header, string(category), mimeType)
	if err != nil {
		if strings.Contains(err.Error(), "uploaded file is empty") {
			return nil, ErrEmptyFile
		}
		return nil, fmt.Errorf("%w: %v", ErrStorageFailed, err)
	}

	resp := &UploadResponse{
		MediaID:      stored.SHA256,
		URL:          stored.URL,
		Path:         stored.Path,
		FileName:     stored.FileName,
		OriginalName: header.Filename,
		Size:         stored.Size,
		MimeType:     mimeType,
		SHA256:       stored.SHA256,
		Duplicate:    stored.Duplicate,
	}

	return resp, nil
}

func (s *Service) DeleteUploadedFile(ctx context.Context, upload *UploadResponse) error {
	if upload == nil || upload.Path == "" || upload.Duplicate {
		return nil
	}

	return s.storage.DeleteFile(ctx, upload.Path)
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
