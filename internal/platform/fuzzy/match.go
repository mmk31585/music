package fuzzy

import (
	"math"
	"strings"
	"unicode"
)

func Normalize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func Levenshtein(a, b string) int {
	na, nb := len(a), len(b)
	if na == 0 {
		return nb
	}
	if nb == 0 {
		return na
	}

	prev := make([]int, nb+1)
	curr := make([]int, nb+1)

	for j := 0; j <= nb; j++ {
		prev[j] = j
	}

	for i := 1; i <= na; i++ {
		curr[0] = i
		for j := 1; j <= nb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[nb]
}

func TitleSimilarity(a, b string) float64 {
	na, nb := Normalize(a), Normalize(b)
	if na == "" && nb == "" {
		return 1.0
	}
	if na == "" || nb == "" {
		return 0.0
	}
	dist := Levenshtein(na, nb)
	maxLen := max(len(na), len(nb))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxLen)
}

func ArtistSimilarity(a, b string) float64 {
	na, nb := Normalize(a), Normalize(b)
	if na == "" || nb == "" {
		return 0.0
	}
	tokensA := strings.Fields(na)
	tokensB := strings.Fields(nb)

	intersection := 0
	for _, ta := range tokensA {
		for _, tb := range tokensB {
			if ta == tb {
				intersection++
				break
			}
		}
	}
	union := len(tokensA) + len(tokensB) - intersection
	if union == 0 {
		return 0.0
	}
	return float64(intersection) / float64(union)
}

func DurationSimilarity(a, b int) float64 {
	if a <= 0 || b <= 0 {
		return 0.5
	}
	diff := math.Abs(float64(a - b))
	switch {
	case diff <= 3:
		return 1.0
	case diff <= 10:
		return 0.8
	case diff <= 30:
		return 0.5
	default:
		return 0.0
	}
}

func CombinedScore(titleA, artistA, albumA string, durA int, titleB, artistB, albumB string, durB int) float64 {
	ts := TitleSimilarity(titleA, titleB)
	as := ArtistSimilarity(artistA, artistB)
	albs := TitleSimilarity(albumA, albumB)
	ds := DurationSimilarity(durA, durB)

	return (ts*0.35 + as*0.35 + albs*0.15 + ds*0.15)
}
