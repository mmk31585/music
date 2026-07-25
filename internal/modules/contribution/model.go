package contribution

import "time"

type Contribution struct {
	ID                   string     `json:"id" db:"id"`
	ContributorID        string     `json:"contributor_id" db:"contributor_id"`
	ContributionType     string     `json:"contribution_type" db:"contribution_type"`
	TargetType           string     `json:"target_type" db:"target_type"`
	TargetID             string     `json:"target_id" db:"target_id"`
	Locale               *string    `json:"locale" db:"locale"`
	Data                 string     `json:"data" db:"data"`
	Summary              *string    `json:"summary" db:"summary"`
	Status               string     `json:"status" db:"status"`
	AIVerdict            *string    `json:"ai_verdict" db:"ai_verdict"`
	AIConfidence         *float64   `json:"ai_confidence" db:"ai_confidence"`
	AIReason             *string    `json:"ai_reason" db:"ai_reason"`
	ModeratorID          *string    `json:"moderator_id" db:"moderator_id"`
	ModeratorNote        *string    `json:"moderator_note" db:"moderator_note"`
	Version              int        `json:"version" db:"version"`
	IsMinor              bool       `json:"is_minor" db:"is_minor"`
	XPAwarded            int        `json:"xp_awarded" db:"xp_awarded"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	DecidedAt            *time.Time `json:"decided_at" db:"decided_at"`
	AppliedAt            *time.Time `json:"applied_at" db:"applied_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	ContributorUsername  *string    `json:"contributor_username,omitempty" db:"contributor_username"`
	ContributorAvatarURL *string    `json:"contributor_avatar_url,omitempty" db:"contributor_avatar_url"`
}

type ContributionHistory struct {
	ID             string    `json:"id" db:"id"`
	ContributionID string    `json:"contribution_id" db:"contribution_id"`
	Data           string    `json:"data" db:"data"`
	PreviousData   *string   `json:"previous_data" db:"previous_data"`
	ChangedBy      string    `json:"changed_by" db:"changed_by"`
	ChangeType     string    `json:"change_type" db:"change_type"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ContentVersion struct {
	ID         string    `json:"id" db:"id"`
	TargetType string    `json:"target_type" db:"target_type"`
	TargetID   string    `json:"target_id" db:"target_id"`
	Version    int       `json:"version" db:"version"`
	Data       string    `json:"data" db:"data"`
	AppliedBy  string    `json:"applied_by" db:"applied_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ContributionType string

const (
	ContributionTypeLyrics      ContributionType = "lyrics"
	ContributionTypeTranslation ContributionType = "translation"
	ContributionTypeCredits     ContributionType = "credits"
	ContributionTypeMetadata    ContributionType = "metadata"
	ContributionTypeAlbumArt    ContributionType = "album_art"
	ContributionTypeBio         ContributionType = "bio"
)

type TargetType string

const (
	TargetTypeTrack  TargetType = "track"
	TargetTypeAlbum  TargetType = "album"
	TargetTypeArtist TargetType = "artist"
)

type ContributionStatus string

const (
	StatusPending     ContributionStatus = "pending"
	StatusApproved    ContributionStatus = "approved"
	StatusRejected    ContributionStatus = "rejected"
	StatusNeedsReview ContributionStatus = "needs_review"
)

type ChangeType string

const (
	ChangeTypeCreate  ChangeType = "create"
	ChangeTypeApprove ChangeType = "approve"
	ChangeTypeReject  ChangeType = "reject"
	ChangeTypeRevert  ChangeType = "revert"
)

type CreateContributionRequest struct {
	ContributionType string `json:"contribution_type" binding:"required,oneof=lyrics translation credits metadata album_art bio"`
	TargetType       string `json:"target_type" binding:"required,oneof=track album artist"`
	TargetID         string `json:"target_id" binding:"required,uuid"`
	Locale           string `json:"locale" binding:"omitempty,min=2,max=10"`
	Data             any    `json:"data" binding:"required"`
	Summary          string `json:"summary" binding:"omitempty,max=500"`
	IsMinor          bool   `json:"is_minor"`
}

type ReviewContributionRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Note   string `json:"note" binding:"omitempty,max=1000"`
}

type ContributionResponse struct {
	ID                   string     `json:"id"`
	ContributorID        string     `json:"contributor_id"`
	ContributorUsername  *string    `json:"contributor_username,omitempty"`
	ContributorAvatarURL *string    `json:"contributor_avatar_url,omitempty"`
	ContributionType     string     `json:"contribution_type"`
	TargetType           string     `json:"target_type"`
	TargetID             string     `json:"target_id"`
	Locale               *string    `json:"locale"`
	Data                 any        `json:"data"`
	Summary              *string    `json:"summary"`
	Status               string     `json:"status"`
	AIVerdict            *string    `json:"ai_verdict"`
	AIConfidence         *float64   `json:"ai_confidence"`
	Version              int        `json:"version"`
	IsMinor              bool       `json:"is_minor"`
	XPAwarded            int        `json:"xp_awarded"`
	CreatedAt            time.Time  `json:"created_at"`
	DecidedAt            *time.Time `json:"decided_at"`
	AppliedAt            *time.Time `json:"applied_at"`
	ModeratorNote        *string    `json:"moderator_note,omitempty"`
}

type ContributionHistoryResponse struct {
	ID         string    `json:"id"`
	Data       any       `json:"data"`
	Previous   any       `json:"previous,omitempty"`
	ChangedBy  string    `json:"changed_by"`
	ChangeType string    `json:"change_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type ContributorStats struct {
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	AvatarURL          string `json:"avatar_url"`
	TotalContributions int    `json:"total_contributions"`
	ApprovedCount      int    `json:"approved_count"`
	PendingCount       int    `json:"pending_count"`
	RejectedCount      int    `json:"rejected_count"`
	XPEarned           int    `json:"xp_earned"`
	Rank               int    `json:"rank"`
}

func ToContributionResponse(c Contribution, data any) ContributionResponse {
	return ContributionResponse{
		ID:                   c.ID,
		ContributorID:        c.ContributorID,
		ContributorUsername:  c.ContributorUsername,
		ContributorAvatarURL: c.ContributorAvatarURL,
		ContributionType:     c.ContributionType,
		TargetType:           c.TargetType,
		TargetID:             c.TargetID,
		Locale:               c.Locale,
		Data:                 data,
		Summary:              c.Summary,
		Status:               c.Status,
		AIVerdict:            c.AIVerdict,
		AIConfidence:         c.AIConfidence,
		Version:              c.Version,
		IsMinor:              c.IsMinor,
		XPAwarded:            c.XPAwarded,
		CreatedAt:            c.CreatedAt,
		DecidedAt:            c.DecidedAt,
		AppliedAt:            c.AppliedAt,
		ModeratorNote:        c.ModeratorNote,
	}
}

func ToHistoryResponse(h ContributionHistory, data any, prev any) ContributionHistoryResponse {
	return ContributionHistoryResponse{
		ID:         h.ID,
		Data:       data,
		Previous:   prev,
		ChangedBy:  h.ChangedBy,
		ChangeType: h.ChangeType,
		CreatedAt:  h.CreatedAt,
	}
}
