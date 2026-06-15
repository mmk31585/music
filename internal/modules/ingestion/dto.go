package ingestion

import (
	"time"
)

type UploadResponse struct {
	DraftID           string            `json:"draftId"`
	OriginalFilename  string            `json:"originalFilename"`
	FileSize          int64             `json:"fileSize"`
	Format            string            `json:"format"`
	DurationSeconds   *float64          `json:"durationSeconds,omitempty"`
	Bitrate           *int              `json:"bitrate,omitempty"`
	Status            DraftStatus       `json:"status"`
	ExtractedMetadata *ExtractedTags    `json:"extractedMetadata"`
	CoverArtURL       *string           `json:"coverArtUrl,omitempty"`
	Assets            []AssetResponse   `json:"assets,omitempty"`
	CreatedAt         time.Time         `json:"createdAt"`
}

type AssetResponse struct {
	ID        string    `json:"id"`
	AssetType string    `json:"assetType"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExtractedTags struct {
	Title       string `json:"title,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Album       string `json:"album,omitempty"`
	AlbumArtist string `json:"albumArtist,omitempty"`
	TrackNumber int    `json:"trackNumber,omitempty"`
	TrackTotal  int    `json:"trackTotal,omitempty"`
	DiscNumber  int    `json:"discNumber,omitempty"`
	DiscTotal   int    `json:"discTotal,omitempty"`
	Year        int    `json:"year,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Composer    string `json:"composer,omitempty"`
	Lyrics      string `json:"lyrics,omitempty"`
	Duration    float64  `json:"duration,omitempty"`
	Bitrate     int      `json:"bitrate,omitempty"`
	Format      string   `json:"format,omitempty"`
	HasCoverArt bool     `json:"hasCoverArt"`
}

type DraftListItem struct {
	ID               string       `json:"id"`
	OriginalFilename string       `json:"originalFilename"`
	FileSize         int64        `json:"fileSize"`
	Format           string       `json:"format"`
	DurationSeconds  *float64     `json:"durationSeconds,omitempty"`
	Status           DraftStatus  `json:"status"`
	Title            string       `json:"title,omitempty"`
	Artist           string       `json:"artist,omitempty"`
	Album            string       `json:"album,omitempty"`
	CoverArtURL      *string      `json:"coverArtUrl,omitempty"`
	HasCoverArt      bool         `json:"hasCoverArt"`
	CreatedAt        time.Time    `json:"createdAt"`
}

type DraftDetailResponse struct {
	ID                string          `json:"id"`
	UploadedBy        string          `json:"uploadedBy"`
	OriginalFilename  string          `json:"originalFilename"`
	FileSize          int64           `json:"fileSize"`
	Format            string          `json:"format"`
	DurationSeconds   *float64        `json:"durationSeconds,omitempty"`
	Bitrate           *int            `json:"bitrate,omitempty"`
	Status            DraftStatus     `json:"status"`
	ExtractedMetadata *ExtractedTags  `json:"extractedMetadata"`
	EnrichedMetadata  interface{}     `json:"enrichedMetadata,omitempty"`
	FinalMetadata     interface{}     `json:"finalMetadata,omitempty"`
	Assets            []AssetResponse `json:"assets,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

type ListDraftsResponse struct {
	Items []DraftListItem `json:"items"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type SaveFinalMetadataRequest struct {
	Artist FinalArtistMetadata `json:"artist" validate:"required"`
	Album  FinalAlbumMetadata  `json:"album" validate:"required"`
	Track  FinalTrackMetadata  `json:"track" validate:"required"`
}

type FinalArtistMetadata struct {
	Action          string  `json:"action"`
	ExistingID      *string `json:"existingId,omitempty"`
	Name            string  `json:"name"`
	Bio             string  `json:"bio,omitempty"`
	ImageURL        string  `json:"imageUrl,omitempty"`
	Country         string  `json:"country,omitempty"`
	MusicBrainzMBID string  `json:"musicbrainzMbid,omitempty"`
}

type FinalAlbumMetadata struct {
	Action               string  `json:"action"`
	ExistingID           *string `json:"existingId,omitempty"`
	Title                string  `json:"title"`
	ReleaseYear          int     `json:"releaseYear,omitempty"`
	Genre                string  `json:"genre,omitempty"`
	CoverURL             string  `json:"coverUrl,omitempty"`
	MusicBrainzReleaseID string  `json:"musicbrainzReleaseId,omitempty"`
}

type FinalTrackMetadata struct {
	Title            string `json:"title"`
	TrackNumber      int    `json:"trackNumber,omitempty"`
	DurationSeconds  int    `json:"durationSeconds"`
	Genre            string `json:"genre,omitempty"`
	Lyrics           string `json:"lyrics,omitempty"`
	Explicit         bool   `json:"explicit"`
	SpotifyPreviewURL string `json:"spotifyPreviewUrl,omitempty"`
	CoverURL         string `json:"coverUrl,omitempty"`
}

type RejectDraftRequest struct {
	Reason string `json:"reason,omitempty"`
}

type ArtistSearchResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Bio      string `json:"bio,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	Country  string `json:"country,omitempty"`
}

type AlbumSearchResult struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	ArtistName string  `json:"artistName"`
	CoverURL   string  `json:"coverUrl,omitempty"`
	ReleaseYear *int   `json:"releaseYear,omitempty"`
}
