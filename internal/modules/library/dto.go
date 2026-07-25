package library

type LikeTrackRequest struct {
	TrackID string `json:"track_id" binding:"required"`
}

type LikeAlbumRequest struct {
	AlbumID string `json:"album_id" binding:"required"`
}

type FollowArtistRequest struct {
	ArtistID string `json:"artist_id" binding:"required"`
}

type AddPlayHistoryRequest struct {
	TrackID   string `json:"track_id" binding:"required"`
	Duration  *int   `json:"duration,omitempty"`
	Completed *bool  `json:"completed,omitempty"`
}
