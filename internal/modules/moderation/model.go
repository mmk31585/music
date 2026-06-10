package moderation

import (
	"time"

	"github.com/google/uuid"
)

type ContentReport struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ReporterID     uuid.UUID  `db:"reporter_id" json:"reporter_id"`
	TargetID       string     `db:"target_id" json:"target_id"`
	TargetType     string     `db:"target_type" json:"target_type"`
	Reason         string     `db:"reason" json:"reason"`
	Description    *string    `db:"description" json:"description,omitempty"`
	Status         string     `db:"status" json:"status"`
	ModeratorID    *uuid.UUID `db:"moderator_id" json:"moderator_id,omitempty"`
	ResolvedAt     *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	ResolutionNote *string    `db:"resolution_note" json:"resolution_note,omitempty"`
}

type ContentFlag struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	TargetID   string     `db:"target_id" json:"target_id"`
	TargetType string     `db:"target_type" json:"target_type"`
	FlagType   string     `db:"flag_type" json:"flag_type"`
	FlaggedAt  time.Time  `db:"flagged_at" json:"flagged_at"`
	ExpiresAt  *time.Time `db:"expires_at" json:"expires_at,omitempty"`
}

type ModerationAction struct {
	ID             uuid.UUID              `db:"id" json:"id"`
	ReportID       *uuid.UUID             `db:"report_id" json:"report_id,omitempty"`
	ModeratorID    uuid.UUID              `db:"moderator_id" json:"moderator_id"`
	Action         string                 `db:"action" json:"action"`
	TargetID       string                 `db:"target_id" json:"target_id"`
	TargetType     string                 `db:"target_type" json:"target_type"`
	PreviousStatus *string                `db:"previous_status" json:"previous_status,omitempty"`
	NewStatus      *string                `db:"new_status" json:"new_status,omitempty"`
	Note           *string                `db:"note" json:"note,omitempty"`
	Metadata       map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt      time.Time              `db:"created_at" json:"created_at"`
}

type ModerationStats struct {
	TotalReports     int64            `db:"total_reports" json:"total_reports"`
	PendingReports   int64            `db:"pending_reports" json:"pending_reports"`
	ResolvedToday    int64            `db:"resolved_today" json:"resolved_today"`
	FlaggedContent   int64            `db:"flagged_content" json:"flagged_content"`
	UniqueReporters  int64            `db:"unique_reporters" json:"unique_reporters"`
	AvgResolutionHrs float64          `db:"hours" json:"avg_resolution_hours"`
	ByReason         map[string]int64 `json:"by_reason"`
	ByTargetType     map[string]int64 `json:"by_target_type"`
}

type ReportWithFlags struct {
	ContentReport
	Flags []ContentFlag `json:"flags,omitempty"`
}
