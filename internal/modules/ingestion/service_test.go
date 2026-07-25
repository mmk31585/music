package ingestion

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestExtractMetadata_WithFullTags(t *testing.T) {
	data := buildMP3WithTagsAndCover()

	result, err := extractMetadata("test.mp3", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("extractMetadata failed: %v", err)
	}
	if result == nil || result.Tags == nil {
		t.Fatal("expected non-nil result and tags")
	}

	tags := result.Tags

	if tags.Title != testTitle {
		t.Errorf("Title = %q, want %q", tags.Title, testTitle)
	}
	if tags.Artist != testArtist {
		t.Errorf("Artist = %q, want %q", tags.Artist, testArtist)
	}
	if tags.Album != testAlbum {
		t.Errorf("Album = %q, want %q", tags.Album, testAlbum)
	}
	if tags.AlbumArtist != testAlbumArtist {
		t.Errorf("AlbumArtist = %q, want %q", tags.AlbumArtist, testAlbumArtist)
	}
	if tags.Genre != testGenre {
		t.Errorf("Genre = %q, want %q", tags.Genre, testGenre)
	}
	if tags.Composer != testComposer {
		t.Errorf("Composer = %q, want %q", tags.Composer, testComposer)
	}
	if tags.Comment != testComment {
		t.Errorf("Comment = %q, want %q", tags.Comment, testComment)
	}
	if tags.Year != testYear {
		t.Errorf("Year = %d, want %d", tags.Year, testYear)
	}
	if tags.TrackNumber != testTrackNum {
		t.Errorf("TrackNumber = %d, want %d", tags.TrackNumber, testTrackNum)
	}
	if tags.TrackTotal != testTrackTotal {
		t.Errorf("TrackTotal = %d, want %d", tags.TrackTotal, testTrackTotal)
	}
	if tags.DiscNumber != testDiscNum {
		t.Errorf("DiscNumber = %d, want %d", tags.DiscNumber, testDiscNum)
	}
	if tags.DiscTotal != testDiscTotal {
		t.Errorf("DiscTotal = %d, want %d", tags.DiscTotal, testDiscTotal)
	}

	if !strings.Contains(tags.Lyrics, "Test lyrics") {
		t.Errorf("Lyrics should contain 'Test lyrics', got %q", tags.Lyrics)
	}
}

func TestExtractMetadata_WithEmbeddedCover(t *testing.T) {
	data := buildMP3WithTagsAndCover()

	result, err := extractMetadata("test.mp3", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("extractMetadata failed: %v", err)
	}

	if !result.Tags.HasCoverArt {
		t.Error("expected HasCoverArt to be true")
	}
	if len(result.CoverArt) == 0 {
		t.Error("expected non-empty CoverArt bytes")
	}
	if result.CoverExt == "" {
		t.Error("expected non-empty CoverExt")
	}
}

func TestExtractMetadata_NoTags(t *testing.T) {
	data := buildNoTagMP3()

	result, err := extractMetadata("test.mp3", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("extractMetadata failed: %v", err)
	}
	if result == nil || result.Tags == nil {
		t.Fatal("expected non-nil result")
	}

	tags := result.Tags

	if tags.Title != "" {
		t.Errorf("Title should be empty for no-tag file, got %q", tags.Title)
	}
	if tags.Artist != "" {
		t.Errorf("Artist should be empty for no-tag file, got %q", tags.Artist)
	}
	if tags.Album != "" {
		t.Errorf("Album should be empty for no-tag file, got %q", tags.Album)
	}
	if tags.Year != 0 {
		t.Errorf("Year should be 0 for no-tag file, got %d", tags.Year)
	}
	if tags.HasCoverArt {
		t.Error("HasCoverArt should be false for no-tag file")
	}
	if len(result.CoverArt) != 0 {
		t.Error("CoverArt should be empty for no-tag file")
	}

	if tags.Format != "mp3" {
		t.Errorf("Format should be mp3, got %q", tags.Format)
	}
}

func TestExtractMetadata_InvalidFormat(t *testing.T) {
	data := []byte("this is not an audio file")

	_, err := extractMetadata("test.txt", bytes.NewReader(data))
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestExtractMetadata_EmptyData(t *testing.T) {
	data := []byte{}

	result, err := extractMetadata("test.mp3", bytes.NewReader(data))
	if err != nil && result == nil {
		return
	}
	if result != nil && result.Tags != nil {
		if result.Tags.Title != "" {
			t.Error("expected empty title for empty file")
		}
		return
	}
}

func TestFormatFromExt(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".mp3", "mp3"},
		{".flac", "flac"},
		{".ogg", "ogg"},
		{".m4a", "aac"},
		{".aac", "aac"},
		{".wav", ""},
		{".txt", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := formatFromExt(tt.ext)
		if result != tt.expected {
			t.Errorf("formatFromExt(%q) = %q, want %q", tt.ext, result, tt.expected)
		}
	}
}

func TestExtractFirstLines(t *testing.T) {
	tests := []struct {
		input    string
		maxLines int
		expected string
	}{
		{"", 3, ""},
		{"single line", 3, "single line"},
		{"line1\nline2\nline3\nline4", 3, "line1\nline2\nline3"},
		{"line1\nline2", 5, "line1\nline2"},
	}

	for _, tt := range tests {
		result := extractFirstLines(tt.input, tt.maxLines)
		if result != tt.expected {
			t.Errorf("extractFirstLines(%q, %d) = %q, want %q", tt.input, tt.maxLines, result, tt.expected)
		}
	}
}

func TestExtractedTagsToJSON(t *testing.T) {
	tags := &ExtractedTags{
		Title:  "Test",
		Artist: "Artist",
		Year:   2024,
	}

	json, err := extractedTagsToJSON(tags)
	if err != nil {
		t.Fatalf("extractedTagsToJSON failed: %v", err)
	}

	if !strings.Contains(json, "Test") {
		t.Errorf("JSON should contain 'Test', got %s", json)
	}
	if !strings.Contains(json, "Artist") {
		t.Errorf("JSON should contain 'Artist', got %s", json)
	}
}

func TestParseExtractedTags(t *testing.T) {
	json := `{"title":"Parsed Title","artist":"Parsed Artist","year":2022}`

	tags := parseExtractedTags(json)
	if tags == nil {
		t.Fatal("expected non-nil tags")
	}
	if tags.Title != "Parsed Title" {
		t.Errorf("Title = %q, want %q", tags.Title, "Parsed Title")
	}
	if tags.Artist != "Parsed Artist" {
		t.Errorf("Artist = %q, want %q", tags.Artist, "Parsed Artist")
	}
	if tags.Year != 2022 {
		t.Errorf("Year = %d, want %d", tags.Year, 2022)
	}
}

func testReader(t *testing.T) io.Reader {
	return bytes.NewReader([]byte("test"))
}
