package reputation

import "time"

type TrustTier struct {
	ID              int64     `json:"id" db:"id"`
	Slug            string    `json:"slug" db:"slug"`
	Label           string    `json:"label" db:"label"`
	Description     string    `json:"description,omitempty" db:"description"`
	MinScore        float64   `json:"minScore" db:"min_score"`
	MinAccepted     int       `json:"minAccepted" db:"min_accepted"`
	UploadSlots     int       `json:"uploadSlots" db:"upload_slots"`
	AutoPublish     bool      `json:"autoPublish" db:"auto_publish"`
	CanReview       bool      `json:"canReview" db:"can_review"`
	HierarchyLevel  int       `json:"hierarchyLevel" db:"hierarchy_level"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

type UserReputation struct {
	ID                      int64      `json:"id" db:"id"`
	UserID                  string     `json:"userId" db:"user_id"`
	TotalContributions      int        `json:"totalContributions" db:"total_contributions"`
	AcceptedContributions   int        `json:"acceptedContributions" db:"accepted_contributions"`
	RejectedContributions   int        `json:"rejectedContributions" db:"rejected_contributions"`
	PendingContributions    int        `json:"pendingContributions" db:"pending_contributions"`
	TrustScore              float64    `json:"trustScore" db:"trust_score"`
	Tier                    string     `json:"tier" db:"tier"`
	UploadSlots             int        `json:"uploadSlots" db:"upload_slots"`
	AutoPublish             bool       `json:"autoPublish" db:"auto_publish"`
	CanReview               bool       `json:"canReview" db:"can_review"`
	LastContributionAt      *time.Time `json:"lastContributionAt,omitempty" db:"last_contribution_at"`
	ContributionStreakDays  int        `json:"contributionStreakDays" db:"contribution_streak_days"`
	CreatedAt               time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt               time.Time  `json:"updatedAt" db:"updated_at"`
}

type ContributionScore struct {
	ID                 int64      `json:"id" db:"id"`
	UserID             string     `json:"userId" db:"user_id"`
	ContributionType   string     `json:"contributionType" db:"contribution_type"`
	ContributionID     *int64     `json:"contributionId,omitempty" db:"contribution_id"`
	ScoreDelta         float64    `json:"scoreDelta" db:"score_delta"`
	Reason             string     `json:"reason,omitempty" db:"reason"`
	ReviewedBy         *string    `json:"reviewedBy,omitempty" db:"reviewed_by"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
}

// ReputationSummary is the public-facing reputation profile.
type ReputationSummary struct {
	UserID                string  `json:"userId"`
	TrustScore            float64 `json:"trustScore"`
	Tier                  string  `json:"tier"`
	TierLabel             string  `json:"tierLabel"`
	AcceptedContributions int     `json:"acceptedContributions"`
	TotalContributions    int     `json:"totalContributions"`
	UploadSlots           int     `json:"uploadSlots"`
	AutoPublish           bool    `json:"autoPublish"`
	CanReview             bool    `json:"canReview"`
}
