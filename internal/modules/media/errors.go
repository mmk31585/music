package media

import "errors"

var (
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidMimeType  = errors.New("invalid mime type")
	ErrNoFileProvided   = errors.New("no file provided")
	ErrStorageFailed    = errors.New("failed to store file")
	ErrInvalidFieldName = errors.New("invalid upload field name")
	ErrMultipleFiles    = errors.New("multiple upload fields provided")
	ErrEmptyFile        = errors.New("uploaded file is empty")
)
