package importcmd

type SearchResult struct {
	Title       string            `json:"title"`
	Artist      string            `json:"artist"`
	URL         string            `json:"url"`
	Duration    int               `json:"duration"`
	Thumbnail   string            `json:"thumbnail"`
	Source      string            `json:"source"`
	Score       float64           `json:"score,omitempty"`
	ISRC        string            `json:"isrc,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
}

type ImportRequest struct {
	URL    string `json:"url" binding:"required"`
	Source string `json:"source"`
}

type ImportResponse struct {
	JobID    string `json:"jobId,omitempty"`
	DraftID  string `json:"draftId,omitempty"`
	Title    string `json:"title,omitempty"`
	Artist   string `json:"artist,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Message  string `json:"message"`
}

type ProgressResponse struct {
	JobID    string `json:"jobId"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Stage    string `json:"stage"`
	Error    string `json:"error,omitempty"`
	DraftID  string `json:"draftId,omitempty"`
}
