package queue

type AddTrackRequest struct {
	TrackID string `json:"track_id" binding:"required"`

	// mode:
	// "next"  = play next, insert at position 1
	// "later" = play later, insert at end
	//
	// Default: "later"
	Mode string `json:"mode"`
}

type QueueResponse struct {
	Items []QueueItem `json:"items"`
	Count int         `json:"count"`
}

type AddTrackResponse struct {
	Item QueueItem `json:"item"`
}

type ReorderQueueRequest struct {
	Items []ReorderQueueItemRequest `json:"items" binding:"required,min=1"`
}

type ReorderQueueItemRequest struct {
	ID       string `json:"id" binding:"required"`
	Position int    `json:"position" binding:"required,min=1"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
