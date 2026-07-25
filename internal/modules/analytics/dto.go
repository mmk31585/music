package analytics

type TrackEventRequest struct {
	EventType  string         `json:"event_type" binding:"required"`
	TrackID    string         `json:"track_id,omitempty"`
	ArtistID   string         `json:"artist_id,omitempty"`
	AlbumID    string         `json:"album_id,omitempty"`
	PlaylistID string         `json:"playlist_id,omitempty"`
	Query      string         `json:"query,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type TrackEventResponse struct {
	Message string `json:"message"`
}

type OverviewResponse struct {
	TotalTracks      int     `json:"total_tracks"`
	TotalUsers       int     `json:"total_users"`
	TotalAlbums      int     `json:"total_albums"`
	TotalPlays       int     `json:"total_plays"`
	ActiveUsers24h   int     `json:"active_users_last_24h"`
	StorageUsedMB    float64 `json:"storage_used_mb"`
}
