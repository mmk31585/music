package library

import "time"

type LikedTrack struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TrackID   int64     `json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikedAlbum struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	AlbumID   int64     `json:"album_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowedArtist struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ArtistID  int64     `json:"artist_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PlayHistory struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	TrackID  int64     `json:"track_id"`
	PlayedAt time.Time `json:"played_at"`
}

type LibraryTrackItem struct {
	TrackID         int64      `json:"track_id"`
	Title           string     `json:"title"`
	ArtistID        *int64     `json:"artist_id,omitempty"`
	ArtistName      *string    `json:"artist_name,omitempty"`
	AlbumID         *int64     `json:"album_id,omitempty"`
	AlbumTitle      *string    `json:"album_title,omitempty"`
	CoverURL        *string    `json:"cover_url,omitempty"`
	AudioURL        *string    `json:"audio_url,omitempty"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	AddedAt         time.Time  `json:"added_at"`
	LastPlayedAt    *time.Time `json:"last_played_at,omitempty"`
}

type LibraryAlbumItem struct {
	AlbumID     int64      `json:"album_id"`
	Title       string     `json:"title"`
	ArtistID    *int64     `json:"artist_id,omitempty"`
	ArtistName  *string    `json:"artist_name,omitempty"`
	CoverURL    *string    `json:"cover_url,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	AddedAt     time.Time  `json:"added_at"`
}

type LibraryArtistItem struct {
	ArtistID   int64     `json:"artist_id"`
	Name       string    `json:"name"`
	CoverURL   *string   `json:"cover_url,omitempty"`
	FollowedAt time.Time `json:"followed_at"`
}
