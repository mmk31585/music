package enrichment

import "context"

type Source string

const (
	SourceFile       Source = "file"
	SourceMusicBrainz Source = "musicbrainz"
	SourceLastFM     Source = "lastfm"
	SourceSpotify    Source = "spotify"
)

type Confidence string

const (
	ConfidenceExact  Confidence = "exact_match"
	ConfidenceFuzzy  Confidence = "fuzzy"
	ConfidenceFallback Confidence = "fallback"
)

type TrackQuery struct {
	Title  string
	Artist string
	Album  string
}

type MusicBrainzResult struct {
	MBID        string
	Title       string
	ArtistName  string
	ArtistMBID  string
	AlbumName   string
	AlbumMBID   string
	ReleaseYear int
	Duration    int
	Genres      []string
}

type LastFMResult struct {
	PlayCount     int
	ListenerCount int
	Tags          []string
	ArtistBio     string
	SimilarArtists []string
}

type SpotifyResult struct {
	SpotifyID      string
	PreviewURL     string
	AlbumCoverURL  string
	ArtistImageURL string
	Popularity     int
}

type EnrichedSuggestion struct {
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Source     Source      `json:"source"`
	Confidence Confidence  `json:"confidence"`
}

type EnrichmentResult struct {
	MusicBrainz *MusicBrainzResult `json:"musicbrainz,omitempty"`
	LastFM      *LastFMResult      `json:"lastfm,omitempty"`
	Spotify     *SpotifyResult     `json:"spotify,omitempty"`
	Suggestions []EnrichedSuggestion `json:"suggestions"`
	Attempted   bool               `json:"enrichment_attempted"`
}

type MusicBrainzClient interface {
	SearchRecording(ctx context.Context, query TrackQuery) (*MusicBrainzResult, error)
}

type LastFMClient interface {
	SearchTrack(ctx context.Context, query TrackQuery) (*LastFMResult, error)
}

type SpotifyClient interface {
	SearchTrack(ctx context.Context, query TrackQuery) (*SpotifyResult, error)
}
