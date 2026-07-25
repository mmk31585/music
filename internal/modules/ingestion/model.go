package ingestion

import (
	"time"
)

type DraftStatus string

const (
	DraftStatusPending          DraftStatus = "pending"
	DraftStatusEnriching        DraftStatus = "enriching"
	DraftStatusReview           DraftStatus = "review"
	DraftStatusAccepted         DraftStatus = "accepted"
	DraftStatusRejected         DraftStatus = "rejected"
	DraftStatusPublished        DraftStatus = "published"
	DraftStatusEnrichmentFailed DraftStatus = "enrichment_failed"
)

type DraftReference struct {
	ArtistID *string `db:"artist_id" json:"artistId,omitempty"`
	AlbumID  *string `db:"album_id" json:"albumId,omitempty"`
	TrackID  *string `db:"track_id" json:"trackId,omitempty"`
}

type IngestionDraft struct {
	ID                string      `db:"id" json:"id"`
	UploadedBy        string      `db:"uploaded_by" json:"uploadedBy"`
	OriginalFilename  string      `db:"original_filename" json:"originalFilename"`
	FilePath          string      `db:"file_path" json:"filePath"`
	FileSize          int64       `db:"file_size" json:"fileSize"`
	DurationSeconds   *float64    `db:"duration_seconds" json:"durationSeconds,omitempty"`
	Bitrate           *int        `db:"bitrate" json:"bitrate,omitempty"`
	Format            string      `db:"format" json:"format"`
	FileHash          string      `db:"file_hash" json:"fileHash,omitempty"`
	Stale             bool        `db:"stale" json:"stale"`
	Status            DraftStatus `db:"status" json:"status"`
	ExtractedMetadata string      `db:"extracted_metadata" json:"extractedMetadata"`
	EnrichedMetadata  *string     `db:"enriched_metadata" json:"enrichedMetadata,omitempty"`
	FinalMetadata     *string     `db:"final_metadata" json:"finalMetadata,omitempty"`
	ArtistID          *string     `db:"artist_id" json:"artistId,omitempty"`
	AlbumID           *string     `db:"album_id" json:"albumId,omitempty"`
	TrackID           *string     `db:"track_id" json:"trackId,omitempty"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updatedAt"`
}

type IngestionDraftAsset struct {
	ID        string    `db:"id" json:"id"`
	DraftID   string    `db:"draft_id" json:"draftId"`
	AssetType string    `db:"asset_type" json:"assetType"`
	URL       string    `db:"url" json:"url"`
	Source    string    `db:"source" json:"source"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
