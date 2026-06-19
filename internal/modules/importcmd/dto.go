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
	// URL is the direct downloadable link (e.g. YouTube). Optional — if empty,
	// the system will resolve a downloadable URL via the acquisition resolver
	// using Title, Artist, and Source.
	URL    string `json:"url"`
	Source string `json:"source"`

	// Metadata fields — used when URL is empty to resolve a downloadable copy.
	Title       string            `json:"title,omitempty"`
	Artist      string            `json:"artist,omitempty"`
	Duration    int               `json:"duration,omitempty"`
	ISRC        string            `json:"isrc,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
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

// ── Artist Search DTOs ─────────────────────────────────────────────

type ArtistDiscographyResult struct {
	ArtistInfo ArtistInfoDTO   `json:"artist_info"`
	Albums     []AlbumGroupDTO `json:"albums"`
}

type ArtistInfoDTO struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type AlbumGroupDTO struct {
	Title  string           `json:"title"`
	Cover  string           `json:"cover"`
	Source string           `json:"source"`
	Tracks []TrackResultDTO `json:"tracks"`
}

type TrackResultDTO struct {
	Title       string            `json:"title"`
	Duration    int               `json:"duration"`
	Source      string            `json:"source"`
	Album       string            `json:"album,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
}

// ── Batch Import DTOs ─────────────────────────────────────────────

type BatchImportRequest struct {
	Tracks []BatchImportItem `json:"tracks" binding:"required,min=1"`
}

type BatchImportItem struct {
	Title       string            `json:"title" binding:"required"`
	Artist      string            `json:"artist" binding:"required"`
	URL         string            `json:"url,omitempty"`
	Album       string            `json:"album,omitempty"`
	Duration    int               `json:"duration,omitempty"`
	Source      string            `json:"source,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
}

type BatchImportResponse struct {
	BatchID string           `json:"batchId"`
	Jobs    []BatchJobResult `json:"jobs"`
	Message string           `json:"message"`
}

type BatchJobResult struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	JobID  string `json:"jobId,omitempty"`
	Error  string `json:"error,omitempty"`
}

type BatchProgressResponse struct {
	BatchID     string `json:"batchId"`
	Total       int    `json:"total"`
	Completed   int    `json:"completed"`
	Failed      int    `json:"failed"`
	InProgress  int    `json:"inProgress"`
	ProgressPct int    `json:"progressPct"`
}
