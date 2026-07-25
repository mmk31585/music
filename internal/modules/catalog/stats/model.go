package stats

// StatsResponse is the lightweight response for the admin catalog stats endpoint.
type StatsResponse struct {
	TotalTracks  int64 `json:"total_tracks"`
	TotalAlbums  int64 `json:"total_albums"`
	TotalArtists int64 `json:"total_artists"`
	TotalGenres  int64 `json:"total_genres"`
}
