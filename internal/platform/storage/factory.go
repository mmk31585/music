package storage

import (
	"context"
	"errors"
)

type Config struct {
	Driver string
	Local  LocalConfig
	S3     S3Config
}

func New(ctx context.Context, cfg Config) (Storage, error) {
	switch cfg.Driver {
	case "", "local":
		return NewLocalStorage(cfg.Local)

	case "s3":
		return NewS3Storage(ctx, cfg.S3)

	default:
		return nil, errors.New("unsupported storage driver: " + cfg.Driver)
	}
}
