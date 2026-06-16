package importcmd

type SearchResult struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	URL       string `json:"url"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Source    string `json:"source"`
}

type ImportRequest struct {
	URL string `json:"url" binding:"required"`
}

type ImportResponse struct {
	DraftID  string `json:"draftId"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration int    `json:"duration"`
	Message  string `json:"message"`
}
