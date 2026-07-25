package video

import (
	"time"

	"github.com/google/uuid"
)

// VideoType defines the category of a video upload.
type VideoType string

const (
	VideoTypeOfficialMV VideoType = "official_mv" // uploaded by admin
	VideoTypeUserEdit   VideoType = "user_edit"   // uploaded by user
)

// VideoStatus tracks the processing state of a video.
type VideoStatus string

const (
	VideoStatusProcessing VideoStatus = "processing" // audio being replaced
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
)

// Video represents both official MVs and user edits.
// The type field determines who can upload and what processing happens.
type Video struct {
	ID          uuid.UUID   `db:"id" json:"id"`
	TrackID     uuid.UUID   `db:"track_id" json:"track_id"`
	UploaderID  uuid.UUID   `db:"uploader_id" json:"uploader_id"`
	Type        VideoType   `db:"type" json:"type"`
	Status      VideoStatus `db:"status" json:"status"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`

	// Storage paths (via existing storage module)
	RawVideoPath   *string `db:"raw_video_path" json:"raw_video_path,omitempty"`
	FinalVideoPath *string `db:"final_video_path" json:"final_video_path,omitempty"`
	ThumbnailPath  *string `db:"thumbnail_path" json:"thumbnail_path,omitempty"`

	// Video metadata
	DurationMs    int64  `db:"duration_ms" json:"duration_ms"`
	AspectRatio   string `db:"aspect_ratio" json:"aspect_ratio"`
	FileSizeBytes int64  `db:"file_size_bytes" json:"file_size_bytes"`

	// For user edits: which portion of the track's audio is used
	TrackStartMs int64 `db:"track_start_ms" json:"track_start_ms"`
	TrackEndMs   int64 `db:"track_end_ms" json:"track_end_ms"`

	// Social
	ViewCount int64 `db:"view_count" json:"view_count"`
	LikeCount int64 `db:"like_count" json:"like_count"` // denormalized
	IsPublic  bool  `db:"is_public" json:"is_public"`

	// Moderation
	IsApproved bool `db:"is_approved" json:"is_approved"` // false until reviewed

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// VideoLike represents a user's like on a video.
type VideoLike struct {
	ID        uuid.UUID `db:"id" json:"id"`
	VideoID   uuid.UUID `db:"video_id" json:"video_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TrackLikeVisibility stores whether a user's like on a track is public
// (shows on their profile) or private (only in their library).
// A row only exists if the user has LIKED the track — no row = not liked.
// This table works IN ADDITION TO the existing reactions/likes table;
// the reactions table still tracks the like for recommendation purposes,
// this table adds the public/private dimension.
//
// DECISION: New table created because the existing reactions model (reactions.go)
// uses a polymorphic pattern (target_type + target_id) without a visibility field.
// Rather than modifying the polymorphic schema (which could affect all reaction types),
// we add this lightweight lookup table specifically for track like visibility.
type TrackLikeVisibility struct {
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	TrackID    uuid.UUID `db:"track_id" json:"track_id"`
	Visibility string    `db:"visibility" json:"visibility"` // "public" | "private"
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	// PRIMARY KEY (user_id, track_id)
}

// UserMusicStatus is the "در حال گوش دادن" feature — a user's current
// or most-recently-played track, optionally visible to others.
type UserMusicStatus struct {
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`                             // PK
	CurrentTrackID *uuid.UUID `db:"current_track_id" json:"current_track_id,omitempty"` // null = not playing
	Visibility     string     `db:"visibility" json:"visibility"`                       // "public" | "followers" | "private"
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	// TTL: if UpdatedAt is >30 minutes ago, treat as "not currently playing"
	// even if current_track_id is set. Do not auto-delete — just apply this
	// staleness check in the service layer, not the DB.
}
