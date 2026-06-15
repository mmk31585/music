package config

type EnrichmentConfig struct {
	LastFmAPIKey        string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func loadEnrichmentConfig() EnrichmentConfig {
	return EnrichmentConfig{
		LastFmAPIKey:        getEnv("LASTFM_API_KEY", ""),
		SpotifyClientID:     getEnv("SPOTIFY_CLIENT_ID", ""),
		SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", ""),
	}
}
