package playlist

type CollaboratorRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type CollaboratorResponse struct {
	UserID    string `json:"user_id"`
	AddedAt   string `json:"added_at"`
	IsCreator bool   `json:"is_creator"`
}

type CollaboratorsListResponse struct {
	Collaborators []CollaboratorResponse `json:"collaborators"`
}
