package worker

type Status string

const (
	StatusQueued      Status = "queued"
	StatusResolving   Status = "resolving"
	StatusDownloading Status = "downloading"
	StatusExtracting  Status = "extracting"
	StatusUploading   Status = "uploading"
	StatusComplete    Status = "complete"
	StatusFailed      Status = "failed"
)

type Job struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Source   string `json:"source"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Query    string `json:"query"`
	UserID   string `json:"user_id"`
	Status   Status `json:"status"`
	Progress int    `json:"progress"`
	Stage    string `json:"stage"`
	DraftID  string `json:"draft_id,omitempty"`
	Error    string `json:"error,omitempty"`
}

type ProgressUpdate struct {
	JobID    string `json:"job_id"`
	Status   Status `json:"status"`
	Progress int    `json:"progress"`
	Stage    string `json:"stage"`
	Error    string `json:"error,omitempty"`
	DraftID  string `json:"draft_id,omitempty"`
}
