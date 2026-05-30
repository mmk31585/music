package search

type SearchResponse struct {
	Query     string           `json:"query"`
	Tracks    []TrackResult    `json:"tracks"`
	Albums    []AlbumResult    `json:"albums"`
	Artists   []ArtistResult   `json:"artist"`
	Playlists []PlaylistResult `json:"playlists"`
}

type TrackResult struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	ArtistID        *string `json:"artist_id,omitempty"`
	ArtistName      *string `json:"artist_name,omitempty"`
	AlbumID         *string `json:"album_id,omitempty"`
	AlbumTitle      *string `json:"album_title,omitempty"`
	CoverURL        *string `json:"cover_url,omitempty"`
	AudioURL        *string `json:"audio_url,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}

type AlbumResult struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	ArtistID    *string `json:"artist_id,omitempty"`
	ArtistName  *string `json:"artist_name,omitempty"`
	CoverURL    *string `json:"cover_url,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
}

type ArtistResult struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	CoverURL *string `json:"cover_url,omitempty"`
}

type PlaylistResult struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CoverURL    *string `json:"cover_url,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}
