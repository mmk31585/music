package video

// CreateVideoRequest is the payload for uploading a user edit video.
type CreateVideoRequest struct {
	TrackID     string `json:"track_id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	RawVideoURL string `json:"raw_video_url"` // pre-signed upload URL or direct path

	// For user edits: which portion of the track's audio is used
	TrackStartMs int64 `json:"track_start_ms"`
	TrackEndMs   int64 `json:"track_end_ms"`

	DurationMs    int64  `json:"duration_ms"`
	AspectRatio   string `json:"aspect_ratio"`
	FileSizeBytes int64  `json:"file_size_bytes"`
}

// UpdateVideoRequest is the payload for approving or updating a video.
type UpdateVideoRequest struct {
	Status       *string `json:"status"`      // for processing -> ready transition
	IsApproved   *bool   `json:"is_approved"` // admin approval
	IsPublic     *bool   `json:"is_public"`
	FinalURL     *string `json:"final_video_url"`
	ThumbnailURL *string `json:"thumbnail_url"`
}

// AdminUpdateVideoRequest is the payload for admin video editing.
type AdminUpdateVideoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"` // "official_mv" | "user_edit"
	IsPublic    *bool   `json:"is_public,omitempty"`
	IsApproved  *bool   `json:"is_approved,omitempty"`
	Status      *string `json:"status,omitempty"` // "processing" | "ready" | "failed"
}

// VideoVisibilityRequest is the payload for setting track like visibility.
type VideoVisibilityRequest struct {
	Visibility string `json:"visibility" validate:"required,oneof=public private"`
}

// UpdateMusicStatusRequest is the payload for setting the user's current music status.
type UpdateMusicStatusRequest struct {
	TrackID    *string `json:"track_id"` // null = stop broadcasting
	Visibility string  `json:"visibility" validate:"oneof=public followers private"`
}

// VideoResponse is the public-facing video output.
type VideoResponse struct {
	ID             string  `json:"id"`
	TrackID        string  `json:"track_id"`
	UploaderID     string  `json:"uploader_id"`
	Type           string  `json:"type"`
	Status         string  `json:"status"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	FinalVideoPath *string `json:"final_video_path,omitempty"`
	ThumbnailPath  *string `json:"thumbnail_path,omitempty"`
	TrackCoverURL  *string `json:"track_cover_url,omitempty"`
	DurationMs     int64   `json:"duration_ms"`
	AspectRatio    string  `json:"aspect_ratio"`
	ViewCount      int64   `json:"view_count"`
	LikeCount      int64   `json:"like_count"`
	IsPublic       bool    `json:"is_public"`
	IsApproved     bool    `json:"is_approved"`
	TrackStartMs   int64   `json:"track_start_ms"`
	TrackEndMs     int64   `json:"track_end_ms"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ExploreResponse wraps a list of videos for the explore endpoint.
type ExploreResponse struct {
	Items  []VideoResponse `json:"items"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Count  int             `json:"count"`
}

// LikedTrackResponse wraps a liked track ID with its visibility.
type LikedTrackResponse struct {
	TrackID    string `json:"track_id"`
	Visibility string `json:"visibility"`
	CreatedAt  string `json:"created_at"`
}

// MusicStatusResponse is the public-facing music status output.
type MusicStatusResponse struct {
	Playing        bool    `json:"playing"`
	CurrentTrackID *string `json:"current_track_id,omitempty"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
}
