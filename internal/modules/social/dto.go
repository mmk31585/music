package social

type FollowRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type FeedQuery struct {
	Limit  int    `form:"limit,default=20"`
	Offset int    `form:"offset,default=0"`
	Types  string `form:"types,omitempty"`
}

type FollowersResponse struct {
	Items      []UserFollow `json:"items"`
	TotalCount int          `json:"total_count"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
}

type ActivityResponse struct {
	Items      []ActivityFeedItem `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Count   int  `json:"count"`
	HasMore bool `json:"has_more"`
}
