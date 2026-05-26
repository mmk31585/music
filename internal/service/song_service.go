package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"music/internal/config"
	"music/internal/domain"
	"music/internal/repository"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SongService struct {
	songRepo repository.SongRepository
	fileRepo repository.FileRepository
	cfg      *config.Config
}

func NewSongService(songRepo repository.SongRepository, fileRepo repository.FileRepository, cfg *config.Config) *SongService {
	return &SongService{
		songRepo: songRepo,
		fileRepo: fileRepo,
		cfg:      cfg,
	}
}

func (s *SongService) Upload(ctx context.Context, title string, fileHeader *multipart.FileHeader) (*domain.Song, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, errors.New("title is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !s.cfg.AllowedAudioExt[ext] {
		return nil, fmt.Errorf("unsupported audio format: %s", ext)
	}

	if fileHeader.Size > s.cfg.MaxUploadSize {
		return nil, fmt.Errorf("file too large: max %d MB", s.cfg.MaxUploadSize>>20)
	}

	id := uuid.New().String()
	storedFileName := id + ext
	storedPath := storedFileName

	if err := s.fileRepo.Save(fileHeader, storedPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	contentType := detectContentType(fileHeader.Header.Get("Content-Type"), ext)

	song := &domain.Song{
		ID:          id,
		Title:       title,
		FileName:    fileHeader.Filename,
		FilePath:    storedPath,
		ContentType: contentType,
		UploadedAt:  time.Now().UTC(),
	}

	if err := s.songRepo.Save(song); err != nil {
		// rollback file
		_ = s.fileRepo.Delete(storedPath)
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return song, nil
}

func (s *SongService) List(ctx context.Context) ([]domain.Song, error) {
	return s.songRepo.FindAll()
}

func (s *SongService) Get(ctx context.Context, id string) (*domain.Song, error) {
	return s.songRepo.FindByID(id)
}

func (s *SongService) Delete(ctx context.Context, id string) error {
	song, err := s.songRepo.FindByID(id)
	if err != nil {
		return err
	}
	if err := s.songRepo.Delete(id); err != nil {
		return err
	}
	// delete the actual audio file
	_ = s.fileRepo.Delete(song.FilePath) // ignore missing
	return nil
}

func (s *SongService) GetAudioFilePath(id string) (string, string, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return "", "", errors.New("song id is required")
	}

	song, err := s.songRepo.FindByID(id)
	if err != nil {
		return "", "", err
	}

	fullPath := filepath.Join(s.cfg.UploadDir, song.FilePath)

	return fullPath, song.ContentType, nil
}

func detectContentType(uploadedContentType, ext string) string {
	if strings.HasPrefix(uploadedContentType, "audio/") {
		return uploadedContentType
	}
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	default:
		return "application/octet-stream"
	}
}
