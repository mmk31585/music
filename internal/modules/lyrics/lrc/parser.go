// internal/modules/lyrics/lrc/parser.go

package lrc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	timestampRegex      = regexp.MustCompile(`^\[(\d{1,2}):(\d{2})(?:\.(\d{2,3}))?\]`) // Updated for 2 or 3 digit ms
	ErrInvalidTimestamp = errors.New("invalid timestamp format")
)

// Changed from LyricLine to LyricsLine to match Service
type LyricsLine struct {
	TimeSeconds float64
	Text        string
}

func ParseLRC(content string) ([]LyricsLine, error) {
	if content == "" {
		return nil, nil
	}

	lines := strings.Split(content, "\n")
	result := make([]LyricsLine, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "[") {
			continue
		}

		lyricLine, err := parseLRCLine(line)
		if err != nil {
			continue // Skip malformed lines rather than failing the whole file
		}
		result = append(result, lyricLine)
	}
	return result, nil
}

func parseLRCLine(line string) (LyricsLine, error) {
	matches := timestampRegex.FindStringSubmatch(line)
	if matches == nil {
		return LyricsLine{}, ErrInvalidTimestamp
	}

	min, _ := strconv.Atoi(matches[1])
	sec, _ := strconv.Atoi(matches[2])

	var ms float64
	if matches[3] != "" {
		m, _ := strconv.Atoi(matches[3])
		if len(matches[3]) == 2 {
			ms = float64(m) / 100
		} else {
			ms = float64(m) / 1000
		}
	}

	return LyricsLine{
		TimeSeconds: float64(min*60) + float64(sec) + ms,
		Text:        strings.TrimSpace(line[len(matches[0]):]),
	}, nil
}
