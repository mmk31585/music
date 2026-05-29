package media

type UploadResponse struct {
	MediaID         string `json:"mediaId"`
	URL             string `json:"url"`
	Path            string `json:"path"`
	FileName        string `json:"fileName"`
	OriginalName    string `json:"originalName"`
	Size            int64  `json:"size"`
	MimeType        string `json:"mimeType"`
	SHA256          string `json:"sha256"`
	Duplicate       bool   `json:"duplicate"`
	DurationSeconds *int   `json:"durationSeconds,omitempty"`
}
