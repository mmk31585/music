package player

type PlaybackTrackResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	ArtistName      string  `json:"artistName"`
	AlbumTitle      *string `json:"albumTitle,omitempty"`
	CoverURL        *string `json:"coverUrl,omitempty"`
	DurationSeconds *int    `json:"durationSeconds,omitempty"`
	StreamURL       string  `json:"streamUrl"`
}
