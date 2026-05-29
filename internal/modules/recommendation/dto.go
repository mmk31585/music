package recommendation

type TrackItem struct {
	ID              string   `json:"id" db:"id"`
	Title           string   `json:"title" db:"title"`
	ArtistID        *string  `json:"artist_id,omitempty" db:"artist_id"`
	ArtistName      *string  `json:"artist_name,omitempty" db:"artist_name"`
	AlbumID         *string  `json:"album_id,omitempty" db:"album_id"`
	AlbumTitle      *string  `json:"album_title,omitempty" db:"album_title"`
	Genre           *string  `json:"genre,omitempty" db:"genre"`
	CoverURL        *string  `json:"cover_url,omitempty" db:"cover_url"`
	AudioURL        *string  `json:"audio_url,omitempty" db:"audio_url"`
	DurationSeconds *int     `json:"duration_seconds,omitempty" db:"duration_seconds"`
	Score           *float64 `json:"score,omitempty"`
}

type RecommendationResponse struct {
	Type  string      `json:"type"`
	Items []TrackItem `json:"items"`
	Limit int         `json:"limit"`
}
