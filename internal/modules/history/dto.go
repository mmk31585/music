package history

type HistoryResponse struct {
	Items      []ListeningHistoryItem `json:"items"`
	Pagination PaginationResponse     `json:"pagination"`
}

type PaginationResponse struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Count   int  `json:"count"`
	HasMore bool `json:"has_more"`
}

// This is not exposed as a route yet.
// It is useful for your playback/player module later.
type RecordListeningRequest struct {
	TrackID         string `json:"track_id" binding:"required"`
	Duration        int    `json:"duration" binding:"min=0"`
	Completed       bool   `json:"completed"`
	SessionID       string `json:"session_id,omitempty"`
	TrackDurationMs int64  `json:"track_duration_ms,omitempty"`
}
