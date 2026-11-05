package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetURL(ctx context.Context, key string) (string, error)
	Copy(ctx context.Context, srcKey, dstKey string) error
	// PresignPut returns a time-limited presigned PUT URL for direct client-to-S3 upload.
	PresignPut(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error)
}
