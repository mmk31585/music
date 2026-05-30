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
