package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalConfig struct {
	BaseDir string
	BaseURL string
}

type LocalStorage struct {
	baseDir string
	baseURL string
}

func NewLocalStorage(cfg LocalConfig) (*LocalStorage, error) {
	if cfg.BaseDir == "" {
		return nil, errors.New("local storage base dir is required")
	}

	if err := os.MkdirAll(cfg.BaseDir, 0755); err != nil {
		return nil, err
	}

	return &LocalStorage{
		baseDir: cfg.BaseDir,
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
	}, nil
}

func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path, err := s.safePath(key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path, err := s.safePath(key)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	path, err := s.safePath(key)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	key = strings.TrimLeft(key, "/")

	if s.baseURL == "" {
		return key, nil
	}

	return s.baseURL + "/" + key, nil
}

func (s *LocalStorage) safePath(key string) (string, error) {
	cleanKey := filepath.Clean(key)

	if cleanKey == "." || cleanKey == "" {
		return "", errors.New("invalid storage key")
	}

	if strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return "", errors.New("invalid storage key")
	}

	return filepath.Join(s.baseDir, cleanKey), nil
}
