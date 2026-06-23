package recommendation

type HomeFeedSection struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Subtitle  string      `json:"subtitle,omitempty"`
	Type      string      `json:"type"`
	Items     []TrackItem `json:"items"`
	SeedTrack *TrackItem  `json:"seed_track,omitempty"`
}

type HomeFeedResponse struct {
	Sections []HomeFeedSection `json:"sections"`
}
