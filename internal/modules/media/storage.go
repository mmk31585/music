package media

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type Storage interface {
	SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, category string, mimeType string) (*StoredFile, error)
	DeleteFile(ctx context.Context, relativePath string) error
}

type StoredFile struct {
	URL       string
	Path      string
	FileName  string
	Size      int64
	SHA256    string
	Duplicate bool
}

type LocalStorageConfig struct {
	BasePath   string
	PublicBase string
}

type LocalStorage struct {
	cfg LocalStorageConfig
}

func NewLocalStorage(cfg LocalStorageConfig) *LocalStorage {
	return &LocalStorage{cfg: cfg}
}

func (s *LocalStorage) SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, category string, mimeType string) (*StoredFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if s.cfg.BasePath == "" {
		return nil, fmt.Errorf("media base path is empty")
	}

	category = sanitizeCategory(category)
	if category == "" {
		category = "misc"
	}

	if err := os.MkdirAll(s.cfg.BasePath, 0o755); err != nil {
		return nil, err
	}

	categoryDir := filepath.Join(s.cfg.BasePath, category)
	if err := os.MkdirAll(categoryDir, 0o755); err != nil {
		return nil, err
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}

	ext := extensionForMime(mimeType)
	if ext == "" {
		ext = sanitizeExtension(strings.ToLower(filepath.Ext(header.Filename)))
	}
	if ext == "" {
		ext = extensionForMime(header.Header.Get("Content-Type"))
	}
	if ext == "" {
		ext = ".bin"
	}

	tempName := "." + uuid.NewString() + ".tmp"
	tempPath := filepath.Join(categoryDir, tempName)

	dst, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(dst, hasher), file)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(tempPath)
		return nil, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return nil, closeErr
	}
	if written == 0 {
		_ = os.Remove(tempPath)
		return nil, fmt.Errorf("uploaded file is empty")
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	filename := hash + ext
	fullPath := filepath.Join(categoryDir, filename)
	duplicate := false

	if _, err := os.Stat(fullPath); err == nil {
		duplicate = true
		_ = os.Remove(tempPath)
	} else if os.IsNotExist(err) {
		if err := os.Rename(tempPath, fullPath); err != nil {
			_ = os.Remove(tempPath)
			return nil, err
		}
	} else {
		_ = os.Remove(tempPath)
		return nil, err
	}

	relPath := filepath.ToSlash(filepath.Join(category, filename))
	publicURL := strings.TrimRight(s.cfg.PublicBase, "/") + "/" + relPath

	return &StoredFile{
		URL:       publicURL,
		Path:      relPath,
		FileName:  filename,
		Size:      written,
		SHA256:    hash,
		Duplicate: duplicate,
	}, nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, relativePath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.cfg.BasePath == "" {
		return fmt.Errorf("media base path is empty")
	}

	cleanPath := filepath.Clean(strings.TrimPrefix(relativePath, "/"))
	if cleanPath == "." || cleanPath == "" || strings.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("invalid media path")
	}

	baseAbs, err := filepath.Abs(s.cfg.BasePath)
	if err != nil {
		return err
	}

	fullPath, err := filepath.Abs(filepath.Join(baseAbs, cleanPath))
	if err != nil {
		return err
	}

	if fullPath != baseAbs && !strings.HasPrefix(fullPath, baseAbs+string(os.PathSeparator)) {
		return fmt.Errorf("media path escapes base directory")
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func sanitizeCategory(category string) string {
	category = strings.ToLower(strings.TrimSpace(category))

	re := regexp.MustCompile(`[^a-z0-9\-_]+`)
	category = re.ReplaceAllString(category, "-")

	category = strings.Trim(category, "-_")

	return category
}

func sanitizeExtension(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))

	if ext == "" {
		return ""
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	re := regexp.MustCompile(`^\.[a-z0-9]+$`)
	if !re.MatchString(ext) {
		return ""
	}

	return ext
}

func extensionForMime(mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	switch mimeType {
	case "audio/mpeg":
		return ".mp3"
	case "audio/ogg":
		return ".ogg"
	case "audio/flac":
		return ".flac"
	case "audio/wav", "audio/x-wav", "audio/wave":
		return ".wav"
	case "audio/mp4", "audio/aac":
		return ".m4a"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}

	if exts, _ := mime.ExtensionsByType(mimeType); len(exts) > 0 {
		return sanitizeExtension(exts[0])
	}

	return ""
}
