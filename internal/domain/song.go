package domain

import "time"

type Song struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}
