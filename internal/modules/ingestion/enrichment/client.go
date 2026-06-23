package enrichment

import "context"

type Source string

const (
	SourceFile        Source = "file"
	SourceMusicBrainz Source = "musicbrainz"
	SourceLastFM      Source = "lastfm"
	SourceSpotify     Source = "spotify"
	SourceLRCLib      Source = "lrclib"
	SourceML          Source = "ml"
)

// EnrichScope defines which sources to run during enrichment.
// Empty or nil means run all sources.
type EnrichScope []Source

type Confidence string

const (
	ConfidenceExact    Confidence = "exact_match"
	ConfidenceFuzzy    Confidence = "fuzzy"
	ConfidenceFallback Confidence = "fallback"
)

type TrackQuery struct {
	Title    string
	Artist   string
	Album    string
	Duration int // seconds, used by LRCLIB for matching
}

type MusicBrainzResult struct {
	MBID        string   `json:"mbid"`
	Title       string   `json:"title"`
	ArtistName  string   `json:"artistName"`
	ArtistMBID  string   `json:"artistMbid"`
	AlbumName   string   `json:"albumName"`
	AlbumMBID   string   `json:"albumMbid"`
	ReleaseYear int      `json:"releaseYear"`
	Duration    int      `json:"duration"`
	Genres      []string `json:"genres"`
}

type LastFMResult struct {
	PlayCount      int      `json:"playCount"`
	ListenerCount  int      `json:"listenerCount"`
	Tags           []string `json:"tags"`
	ArtistBio      string   `json:"artistBio"`
	ArtistImageURL string   `json:"artistImageUrl"`
	AlbumCoverURL  string   `json:"albumCoverUrl"`
	SimilarArtists []string `json:"similarArtists"`
}

type SpotifyResult struct {
	SpotifyID      string `json:"spotifyId"`
	PreviewURL     string `json:"previewUrl"`
	AlbumCoverURL  string `json:"albumCoverUrl"`
	ArtistImageURL string `json:"artistImageUrl"`
	Popularity     int    `json:"popularity"`
}

type MLResult struct {
	Lyrics         string  `json:"lyrics,omitempty"`
	LyricsType     string  `json:"lyricsType,omitempty"`
	AlbumCoverURL  string  `json:"albumCoverUrl,omitempty"`
	ArtistImageURL string  `json:"artistImageUrl,omitempty"`
	Confidence     float64 `json:"confidence,omitempty"`
}

type EnrichedSuggestion struct {
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Source     Source      `json:"source"`
	Confidence Confidence  `json:"confidence"`
}

type EnrichmentResult struct {
	MusicBrainz *MusicBrainzResult   `json:"musicbrainz,omitempty"`
	LastFM      *LastFMResult        `json:"lastfm,omitempty"`
	Spotify     *SpotifyResult       `json:"spotify,omitempty"`
	LRCLib      *LRCLibResult        `json:"lrclib,omitempty"`
	ML          *MLResult            `json:"ml,omitempty"`
	Suggestions []EnrichedSuggestion `json:"suggestions"`
	Attempted   bool                 `json:"enrichment_attempted"`
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

type MLClient interface {
	FetchLyrics(ctx context.Context, query TrackQuery) (*MLResult, error)
	FetchCover(ctx context.Context, query TrackQuery) (*MLResult, error)
}

type DeezerClient interface {
	SearchArtistImage(ctx context.Context, name string) (string, error)
}

type CoverArtClient interface {
	SearchAlbumCover(ctx context.Context, albumTitle, artistName string) (string, error)
}
