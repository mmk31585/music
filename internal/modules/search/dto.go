package search

type SearchResponse struct {
	Query     string           `json:"query"`
	Tracks    []TrackResult    `json:"tracks"`
	Albums    []AlbumResult    `json:"albums"`
	Artists   []ArtistResult   `json:"artists"`
	Playlists []PlaylistResult `json:"playlists"`
}

type TrackArtistResult struct {
	ArtistID string `json:"artistId" db:"artist_id"`
	Name     string `json:"name" db:"name"`
	Slug     string `json:"slug" db:"slug"`
	Role     string `json:"role" db:"role"`
	Position int    `json:"position" db:"position"`
}

type GenreResult struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

type TrackResult struct {
	ID              string              `json:"id" db:"id"`
	Title           string              `json:"title" db:"title"`
	Slug            string              `json:"slug" db:"slug"`
	ArtistID        *string             `json:"artistId,omitempty" db:"artist_id"`
	AlbumID         *string             `json:"albumId,omitempty" db:"album_id"`
	CoverURL        *string             `json:"coverUrl,omitempty" db:"cover_url"`
	AudioURL        *string             `json:"audioUrl,omitempty" db:"audio_url"`
	DurationSeconds int                 `json:"durationSeconds" db:"duration_seconds"`
	TrackNumber     *int                `json:"trackNumber,omitempty" db:"track_number"`
	Explicit        bool                `json:"explicit" db:"explicit"`
	PlayCount       int64               `json:"playCount" db:"play_count"`
	IsPublic        bool                `json:"isPublic" db:"is_public"`
	CreatedAt       string              `json:"createdAt" db:"created_at"`
	Artists         []TrackArtistResult `json:"artists,omitempty"`
	Genres          []GenreResult       `json:"genres,omitempty"`
}

type AlbumResult struct {
	ID          string              `json:"id" db:"id"`
	Title       string              `json:"title" db:"title"`
	Slug        string              `json:"slug" db:"slug"`
	ArtistID    *string             `json:"artistId,omitempty" db:"artist_id"`
	CoverURL    *string             `json:"coverUrl,omitempty" db:"cover_url"`
	ReleaseDate *string             `json:"releaseDate,omitempty" db:"release_date"`
	Artists     []TrackArtistResult `json:"artists,omitempty"`
}

type ArtistResult struct {
	ID               string  `json:"id" db:"id"`
	Name             string  `json:"name" db:"name"`
	Slug             string  `json:"slug" db:"slug"`
	Bio              *string `json:"bio,omitempty" db:"bio"`
	CoverURL         *string `json:"imageUrl,omitempty" db:"image_url"`
	IsVerified       bool    `json:"isVerified" db:"is_verified"`
	MonthlyListeners int     `json:"monthlyListeners" db:"monthly_listeners"`
}

type PlaylistResult struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
	CoverURL    *string `json:"coverUrl,omitempty" db:"cover_url"`
	UserID      *string `json:"userId,omitempty" db:"user_id"`
	IsPublic    *bool   `json:"isPublic,omitempty" db:"is_public"`
}
