package acquisition

import "context"

const (
	QualityFLAC = "flac"
	QualityV0   = "v0"
	Quality320  = "320kbps"
	Quality128  = "128kbps"
)

type Candidate struct {
	URL         string  `json:"url"`
	Source      string  `json:"source"`
	Quality     string  `json:"quality"`
	Confidence  float64 `json:"confidence"`
	AudioFormat string  `json:"audio_format"`
}

type ResolveQuery struct {
	Title       string
	Artist      string
	Album       string
	Duration    int
	ISRC        string
	ExternalIDs map[string]string
}

type Resolver interface {
	Name() string
	Priority() int
	Resolve(ctx context.Context, q ResolveQuery) (*Candidate, error)
}
