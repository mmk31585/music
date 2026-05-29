package library

type LikeTrackRequest struct {
	TrackID int64 `json:"track_id" binding:"required"`
}

type LikeAlbumRequest struct {
	AlbumID int64 `json:"album_id" binding:"required"`
}

type FollowArtistRequest struct {
	ArtistID int64 `json:"artist_id" binding:"required"`
}

type AddPlayHistoryRequest struct {
	TrackID int64 `json:"track_id" binding:"required"`
}
