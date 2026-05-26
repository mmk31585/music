package repository

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalFileRepo struct {
	uploadDir string
}

func NewLocalFileRepo(uploadDir string) *LocalFileRepo {
	return &LocalFileRepo{uploadDir: uploadDir}
}

func (r *LocalFileRepo) Save(fileHeader *multipart.FileHeader, destPath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fullPath := filepath.Join(r.uploadDir, destPath)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (r *LocalFileRepo) Delete(filePath string) error {
	fullPath := filepath.Join(r.uploadDir, filePath)
	return os.Remove(fullPath)
}
