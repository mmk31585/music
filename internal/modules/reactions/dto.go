package reactions

type ReactRequest struct {
	TargetID   string `json:"target_id" validate:"required"`
	TargetType string `json:"target_type" validate:"required,oneof=track album playlist artist comment"`
	Type       string `json:"type" validate:"required,oneof=like love dislike"`
}

type ReactionResponse struct {
	ID         string `json:"id"`
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	CreatedAt  string `json:"created_at"`
}

type CountsResponse struct {
	Like    int `json:"like"`
	Love    int `json:"love"`
	Dislike int `json:"dislike"`
	Total   int `json:"total"`
}
