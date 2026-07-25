package config

import "os"

type EnrichmentConfig struct {
	LastFmAPIKey        string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func loadEnrichmentConfig() EnrichmentConfig {
	return EnrichmentConfig{
		LastFmAPIKey:        os.Getenv("LASTFM_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}
