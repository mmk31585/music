package upload

import "time"

type UploadSource string

const (
	SourceAdmin   UploadSource = "admin"
	SourceListener UploadSource = "listener"
	SourceClub    UploadSource = "club"
)

type CoUploaderRole string

const (
	RoleUploader    CoUploaderRole = "uploader"
	RoleContributor CoUploaderRole = "contributor"
	RoleFeatured    CoUploaderRole = "featured"
)

type Draft struct {
	ID                 string     `json:"id" db:"id"`
	UploadedBy         string     `json:"uploadedBy" db:"uploaded_by"`
	OriginalFilename   string     `json:"originalFilename" db:"original_filename"`
	FilePath           string     `json:"filePath" db:"file_path"`
	FileSize           int64      `json:"fileSize" db:"file_size"`
	DurationSeconds    *float64   `json:"durationSeconds,omitempty" db:"duration_seconds"`
	Bitrate            *int       `json:"bitrate,omitempty" db:"bitrate"`
	Format             string     `json:"format" db:"format"`
	Status             string     `json:"status" db:"status"`
	ExtractedMetadata  string     `json:"extractedMetadata" db:"extracted_metadata"`
	EnrichedMetadata   *string    `json:"enrichedMetadata,omitempty" db:"enriched_metadata"`
	FinalMetadata      *string    `json:"finalMetadata,omitempty" db:"final_metadata"`
	UploadSource       string     `json:"uploadSource" db:"upload_source"`
	ClubID             *string    `json:"clubId,omitempty" db:"club_id"`
	NeedsReview        bool       `json:"needsReview" db:"needs_review"`
	ReviewNotes        *string    `json:"reviewNotes,omitempty" db:"review_notes"`
	ReviewedBy         *string    `json:"reviewedBy,omitempty" db:"reviewed_by"`
	ReviewedAt         *time.Time `json:"reviewedAt,omitempty" db:"reviewed_at"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
}

type CoUploader struct {
	ID                int64     `json:"id" db:"id"`
	TrackID           int64     `json:"trackId" db:"track_id"`
	UserID            string    `json:"userId" db:"user_id"`
	Role              string    `json:"role" db:"role"`
	XpSharePercent    float64   `json:"xpSharePercent" db:"xp_share_percent"`
	ContributionType  *string   `json:"contributionType,omitempty" db:"contribution_type"`
	AddedAt           time.Time `json:"addedAt" db:"added_at"`
}

type UploadSlot struct {
	ID           int64      `json:"id" db:"id"`
	UserID       string     `json:"userId" db:"user_id"`
	UsedSlots    int        `json:"usedSlots" db:"used_slots"`
	MaxSlots     int        `json:"maxSlots" db:"max_slots"`
	LastUploadAt *time.Time `json:"lastUploadAt,omitempty" db:"last_upload_at"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

type UploadSlotStatus struct {
	Available bool `json:"available"`
	Used      int  `json:"used"`
	Max       int  `json:"max"`
}

// CreateUploadRequest is the DTO for initiating a listener upload.
type CreateUploadRequest struct {
	Filename string `json:"filename" binding:"required"`
	ClubID   *string `json:"clubId,omitempty"`
}

// ReviewDraftRequest is the DTO for reviewing a pending draft.
type ReviewDraftRequest struct {
	Action  string  `json:"action" binding:"required,oneof=accept reject"`
	Notes   *string `json:"notes,omitempty"`
}

// CoUploaderRequest is the DTO for adding a co-uploader.
type CoUploaderRequest struct {
	UserID           string  `json:"userId" binding:"required"`
	Role             string  `json:"role" binding:"required,oneof=uploader contributor featured"`
	XpSharePercent   float64 `json:"xpSharePercent" binding:"required,min=0,max=100"`
	ContributionType *string `json:"contributionType,omitempty"`
}
