package media

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID uuid.UUID `db:"id" json:"id"`

	MediaType       string  `db:"media_type" json:"mediaType"`
	StorageProvider string  `db:"storage_provider" json:"storageProvider"`
	Bucket          *string `db:"bucket" json:"bucket,omitempty"`
	ObjectKey       string  `db:"object_key" json:"objectKey"`
	PublicURL       *string `db:"public_url" json:"publicUrl,omitempty"`

	MimeType       *string `db:"mime_type" json:"mimeType,omitempty"`
	FileSize       *int64  `db:"file_size" json:"fileSize,omitempty"`
	ChecksumSHA256 *string `db:"checksum_sha256" json:"checksumSha256,omitempty"`

	DurationSeconds *int `db:"duration_seconds" json:"durationSeconds,omitempty"`
	Width           *int `db:"width" json:"width,omitempty"`
	Height          *int `db:"height" json:"height,omitempty"`

	OriginalFilename *string `db:"original_filename" json:"originalFilename,omitempty"`
	Metadata         string  `db:"metadata" json:"metadata"`

	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type CreateMediaRequest struct {
	MediaType       string
	StorageProvider string
	Bucket          *string
	ObjectKey       string
	PublicURL       *string

	MimeType       *string
	FileSize       *int64
	ChecksumSHA256 *string

	DurationSeconds *int
	Width           *int
	Height          *int

	OriginalFilename *string
	Metadata         string
	CreatedBy        *uuid.UUID
}
