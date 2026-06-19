package search

import (
	"cmp"
	"slices"
	"strings"

	"music/internal/platform/fuzzy"
)

func Deduplicate(results []Result) []Result {
	seen := make(map[string]bool, len(results))
	deduped := make([]Result, 0, len(results))

	for _, r := range results {
		// Exact dedup by ISRC
		if r.ISRC != "" {
			key := "isrc:" + r.ISRC
			if seen[key] {
				continue
			}
			seen[key] = true
		}

		// Exact dedup by URL
		if r.URL != "" {
			key := "url:" + r.URL
			if seen[key] {
				continue
			}
			seen[key] = true
		}

		// Exact dedup by external ID
		for _, id := range []string{r.ExternalIDs.SpotifyID, r.ExternalIDs.DeezerID, r.ExternalIDs.MBID} {
			if id != "" {
				key := "ext:" + id
				if seen[key] {
					continue
				}
				seen[key] = true
			}
		}

		// Fuzzy dedup against already-accepted results
		isDuplicate := false
		for _, existing := range deduped {
			score := fuzzy.CombinedScore(
				r.Title, r.Artist, r.Album, r.Duration,
				existing.Title, existing.Artist, existing.Album, existing.Duration,
			)
			if score >= 0.85 {
				isDuplicate = true
				break
			}
		}
		if isDuplicate {
			continue
		}

		deduped = append(deduped, r)
	}

	return deduped
}

var sourcePriority = map[string]int{
	"spotify":     10,
	"deezer":      9,
	"musicbrainz": 8,
	"lastfm":      7,
	"local":       6,
}

func sourceRank(s string) int {
	if v, ok := sourcePriority[strings.ToLower(s)]; ok {
		return v
	}
	return 0
}

func Rank(results []Result) {
	slices.SortFunc(results, func(a, b Result) int {
		if a.Score != b.Score {
			if a.Score > b.Score {
				return -1
			}
			return 1
		}
		return cmp.Compare(sourceRank(b.Source), sourceRank(a.Source))
	})
}
