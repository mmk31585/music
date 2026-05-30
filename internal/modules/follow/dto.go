package follow

type FollowResponse struct {
	Message string `json:"message"`
}

type FollowStateResponse struct {
	Following bool `json:"following"`
}
