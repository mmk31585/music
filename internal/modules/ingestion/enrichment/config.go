package enrichment

type Config struct {
	MusicBrainz MusicBrainzConfig
	LastFM      LastFMConfig
	Spotify     SpotifyConfig
}

type MusicBrainzConfig struct {
	UserAgent string
	BaseURL   string
}

type LastFMConfig struct {
	APIKey  string
	BaseURL string
}

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	AuthBaseURL  string
}

func DefaultConfig() Config {
	return Config{
		MusicBrainz: MusicBrainzConfig{
			UserAgent: "MojaMusic/1.0 (music-ingestion)",
			BaseURL:   "https://musicbrainz.org/ws/2",
		},
		LastFM: LastFMConfig{
			BaseURL: "https://ws.audioscrobbler.com/2.0",
		},
		Spotify: SpotifyConfig{
			BaseURL:     "https://api.spotify.com/v1",
			AuthBaseURL: "https://accounts.spotify.com/api/token",
		},
	}
}
