package search

type SearchResponse struct {
	Query     string           `json:"query"`
	Tracks    []TrackResult    `json:"tracks"`
	Albums    []AlbumResult    `json:"albums"`
	Artists   []ArtistResult   `json:"artists"`
	Playlists []PlaylistResult `json:"playlists"`
}

type TrackResult struct {
	ID              string  `json:"id" db:"id"`
	Title           string  `json:"title" db:"title"`
	ArtistID        *string `json:"artist_id,omitempty" db:"artist_id"`
	ArtistName      *string `json:"artist_name,omitempty" db:"artist_name"`
	AlbumID         *string `json:"album_id,omitempty" db:"album_id"`
	AlbumTitle      *string `json:"album_title,omitempty" db:"album_title"`
	CoverURL        *string `json:"cover_url,omitempty" db:"cover_url"`
	AudioURL        *string `json:"audio_url,omitempty" db:"audio_url"`
	DurationSeconds *int    `json:"duration_seconds,omitempty" db:"duration_seconds"`
}

type AlbumResult struct {
	ID          string  `json:"id" db:"id"`
	Title       string  `json:"title" db:"title"`
	ArtistID    *string `json:"artist_id,omitempty" db:"artist_id"`
	ArtistName  *string `json:"artist_name,omitempty" db:"artist_name"`
	CoverURL    *string `json:"cover_url,omitempty" db:"cover_url"`
	ReleaseDate *string `json:"release_date,omitempty" db:"release_date"`
}

type ArtistResult struct {
	ID       string  `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	CoverURL *string `json:"cover_url,omitempty" db:"cover_url"`
}

type PlaylistResult struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
	CoverURL    *string `json:"cover_url,omitempty" db:"cover_url"`
	UserID      *string `json:"user_id,omitempty" db:"user_id"`
	IsPublic    *bool   `json:"is_public,omitempty" db:"is_public"`
}
