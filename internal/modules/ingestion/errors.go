package ingestion

import "errors"

var (
	ErrInvalidFileType   = errors.New("unsupported audio file type")
	ErrFileTooLarge      = errors.New("file exceeds maximum allowed size")
	ErrNoFileProvided    = errors.New("no file provided")
	ErrStorageFailed     = errors.New("failed to store file")
	ErrExtractionFailed  = errors.New("failed to extract metadata from file")
	ErrDraftNotFound     = errors.New("ingestion draft not found")
	ErrInvalidStatus     = errors.New("invalid draft status")
)
