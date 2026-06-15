package ingestion

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"music/internal/pkg/audioinfo"

	"github.com/dhowden/tag"
)

type extractionResult struct {
	Tags     *ExtractedTags
	CoverArt []byte
	CoverExt string
	Audio    *audioinfo.Info
}

func extractMetadata(fileName string, reader io.ReadSeeker) (*extractionResult, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	format := formatFromExt(ext)
	if format == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidFileType, ext)
	}

	result := &extractionResult{}
	result.Tags = &ExtractedTags{}

	m, err := tag.ReadFrom(reader)
	if err == nil && m != nil {
		fillTags(m, result.Tags)

		if pic := m.Picture(); pic != nil {
			result.CoverArt = pic.Data
			result.CoverExt = extractExtFromMIME(pic.MIMEType)
			result.Tags.HasCoverArt = true
		}

		if _, err := reader.Seek(0, io.SeekStart); err == nil {
			audioInfo, err := audioinfo.Extract(reader)
			if err == nil && audioInfo != nil {
				result.Audio = audioInfo
				result.Tags.Duration = audioInfo.Duration
				result.Tags.Bitrate = audioInfo.Bitrate
				result.Tags.Format = audioInfo.Format
			}
		}
	} else {
		if _, err := reader.Seek(0, io.SeekStart); err == nil {
			audioInfo, err := audioinfo.Extract(reader)
			if err == nil && audioInfo != nil {
				result.Audio = audioInfo
				result.Tags.Duration = audioInfo.Duration
				result.Tags.Bitrate = audioInfo.Bitrate
				result.Tags.Format = audioInfo.Format
			}
		}
	}

	if result.Tags.Format == "" {
		result.Tags.Format = format
	}

	return result, nil
}

func fillTags(m tag.Metadata, t *ExtractedTags) {
	if v := m.Title(); v != "" {
		t.Title = v
	}
	if v := m.Artist(); v != "" {
		t.Artist = v
	}
	if v := m.Album(); v != "" {
		t.Album = v
	}
	if v := m.AlbumArtist(); v != "" {
		t.AlbumArtist = v
	}
	if track, total := m.Track(); track > 0 {
		t.TrackNumber = track
		t.TrackTotal = total
	}
	if disc, total := m.Disc(); disc > 0 {
		t.DiscNumber = disc
		t.DiscTotal = total
	}
	if v := m.Year(); v > 0 {
		t.Year = v
	}
	if v := m.Genre(); v != "" {
		t.Genre = v
	}
	if v := m.Comment(); v != "" {
		t.Comment = v
	}
	if v := m.Composer(); v != "" {
		t.Composer = v
	}
	if v := m.Lyrics(); v != "" {
		t.Lyrics = v
	}
}

func formatFromExt(ext string) string {
	switch ext {
	case ".mp3":
		return "mp3"
	case ".flac":
		return "flac"
	case ".ogg":
		return "ogg"
	case ".m4a", ".aac":
		return "aac"
	default:
		return ""
	}
}

func extractExtFromMIME(mime string) string {
	switch mime {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".jpg"
	}
}

func extractedTagsToJSON(t *ExtractedTags) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal extracted tags: %w", err)
	}
	return string(data), nil
}

func extractFirstLines(s string, maxLines int) string {
	if s == "" {
		return ""
	}
	lines := strings.SplitN(s, "\n", maxLines+1)
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return strings.Join(lines, "\n")
}

func asJSONPointer(v interface{}) string {
	if v == nil {
		return "null"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func parseExtractedTags(rawJSON string) *ExtractedTags {
	var tags ExtractedTags
	if err := json.Unmarshal([]byte(rawJSON), &tags); err != nil {
		return &ExtractedTags{}
	}
	return &tags
}


