package finalization

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type mockStorage struct {
	copyFn   func(ctx context.Context, srcKey, dstKey string) error
	getURLFn func(ctx context.Context, key string) (string, error)
}

func (m *mockStorage) GetURL(ctx context.Context, key string) (string, error) {
	if m.getURLFn != nil {
		return m.getURLFn(ctx, key)
	}
	return "/storage/" + key, nil
}
func (m *mockStorage) Exists(ctx context.Context, key string) (bool, error) { return true, nil }
func (m *mockStorage) Copy(ctx context.Context, srcKey, dstKey string) error {
	if m.copyFn != nil {
		return m.copyFn(ctx, srcKey, dstKey)
	}
	return nil
}
func (m *mockStorage) Upload(ctx context.Context, key string, _ io.Reader, _ int64, _ string) error {
	return nil
}
func (m *mockStorage) Delete(ctx context.Context, _ string) error { return nil }
func (m *mockStorage) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }

func TestSplitGenres(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"Rock", 1},
		{"Rock, Pop", 2},
		{" Rock , Pop , Jazz ", 3},
		{",Rock,,Pop,", 2},
	}
	for _, tt := range tests {
		got := splitGenres(tt.input)
		if len(got) != tt.want {
			t.Errorf("splitGenres(%q) = %d items, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestExtractStorageKey(t *testing.T) {
	tests := []struct {
		url     string
		baseURL string
		want    string
	}{
		{"/storage/ingestion-audio/file.mp3", "/storage", "ingestion-audio/file.mp3"},
		{"https://cdn.example.com/ingestion-audio/file.mp3", "https://cdn.example.com", "ingestion-audio/file.mp3"},
		{"/ingestion-audio/file.mp3", "", "ingestion-audio/file.mp3"},
	}
	for _, tt := range tests {
		got := extractStorageKey(tt.url, tt.baseURL)
		if got != tt.want {
			t.Errorf("extractStorageKey(%q, %q) = %q, want %q", tt.url, tt.baseURL, got, tt.want)
		}
	}
}

func TestStrPtr(t *testing.T) {
	if v := strPtr(""); v != nil {
		t.Error("expected nil for empty string")
	}
	if v := strPtr("hello"); v == nil || *v != "hello" {
		t.Error("expected pointer to hello")
	}
}

func TestFinalize_InvalidStatus(t *testing.T) {
	s := &Service{logger: zap.NewNop()}
	_, err := s.Finalize(context.Background(), "draft-1", "pending", "audio/file.mp3", "mp3", `{}`)
	if err == nil {
		t.Fatal("expected error for non-accepted status")
	}
}

func TestFinalize_NoMetadata(t *testing.T) {
	s := &Service{logger: zap.NewNop()}
	_, err := s.Finalize(context.Background(), "draft-1", "accepted", "audio/file.mp3", "mp3", "")
	if err == nil {
		t.Fatal("expected error for missing metadata")
	}
}

func TestFinalize_ArtistLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	store := &mockStorage{}
	s := NewService(sqlxDB, store, zap.NewNop())

	meta := `{
		"artist": {"action": "link", "existingId": "550e8400-e29b-41d4-a716-446655440000", "name": "Test Artist"},
		"album": {"action": "skip"},
		"track": {"title": "Test Track", "durationSeconds": 180, "trackNumber": 1, "explicit": false}
	}`

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT asset_type, url FROM ingestion_draft_assets WHERE`).
		WithArgs("draft-1").
		WillReturnRows(sqlmock.NewRows([]string{"asset_type", "url"}))

	mock.ExpectQuery(`INSERT INTO tracks`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", nil, "Test Track", "test-track", 180, 1, false, "/storage/catalog-audio/550e8400-e29b-41d4-a716-446655440000/draft-1-file.mp3", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("660e8400-e29b-41d4-a716-446655440001"))

	mock.ExpectExec(`INSERT INTO track_artists`).
		WithArgs("660e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE ingestion_drafts`).
		WithArgs("published", "550e8400-e29b-41d4-a716-446655440000", sqlmock.AnyArg(), "660e8400-e29b-41d4-a716-446655440001", "draft-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	result, err := s.Finalize(context.Background(), "draft-1", "accepted", "ingestion-audio/draft-1-file.mp3", "mp3", meta)
	if err != nil {
		t.Fatal(err)
	}
	if result.ArtistID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected artist ID 550e8400..., got %s", result.ArtistID)
	}
}

func TestFinalize_CreateArtistAndAlbum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	store := &mockStorage{}
	s := NewService(sqlxDB, store, zap.NewNop())

	meta := `{
		"artist": {"action": "create", "name": "New Artist", "bio": "Great artist", "country": "US"},
		"album": {"action": "create", "title": "New Album", "releaseYear": 2024, "coverUrl": "https://img/cover.jpg"},
		"track": {"title": "New Track", "durationSeconds": 240, "trackNumber": 1, "explicit": true, "genre": "Rock, Pop"}
	}`

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO artists`).
		WithArgs("New Artist", "new-artist", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("aa0e8400-e29b-41d4-a716-446655440000"))

	mock.ExpectQuery(`INSERT INTO albums`).
		WithArgs("aa0e8400-e29b-41d4-a716-446655440000", "New Album", "new-album", sqlmock.AnyArg(), sqlmock.AnyArg(), "album").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("bb0e8400-e29b-41d4-a716-446655440001"))

	mock.ExpectExec(`INSERT INTO album_artists`).
		WithArgs("bb0e8400-e29b-41d4-a716-446655440001", "aa0e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`SELECT asset_type, url FROM ingestion_draft_assets WHERE`).
		WithArgs("draft-2").
		WillReturnRows(sqlmock.NewRows([]string{"asset_type", "url"}))

	mock.ExpectQuery(`INSERT INTO genres`).
		WithArgs("Rock", "rock", "Pop", "pop").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("cc0e8400-e29b-41d4-a716-446655440002").AddRow("dd0e8400-e29b-41d4-a716-446655440003"))

	mock.ExpectQuery(`INSERT INTO tracks`).
		WithArgs("aa0e8400-e29b-41d4-a716-446655440000", sqlmock.AnyArg(), "New Track", "new-track", 240, 1, true, "/storage/catalog-audio/aa0e8400-e29b-41d4-a716-446655440000/draft-2-file.mp3", "https://img/cover.jpg").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("ee0e8400-e29b-41d4-a716-446655440004"))

	mock.ExpectExec(`INSERT INTO track_artists`).
		WithArgs("ee0e8400-e29b-41d4-a716-446655440004", "aa0e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO track_genres`).
		WithArgs("ee0e8400-e29b-41d4-a716-446655440004", "cc0e8400-e29b-41d4-a716-446655440002", "dd0e8400-e29b-41d4-a716-446655440003").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE ingestion_drafts`).
		WithArgs("published", "aa0e8400-e29b-41d4-a716-446655440000", "bb0e8400-e29b-41d4-a716-446655440001", "ee0e8400-e29b-41d4-a716-446655440004", "draft-2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	result, err := s.Finalize(context.Background(), "draft-2", "accepted", "ingestion-audio/draft-2-file.mp3", "mp3", meta)
	if err != nil {
		t.Fatal(err)
	}
	if result.ArtistID != "aa0e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected artist ID, got %s", result.ArtistID)
	}
	if result.AlbumID != "bb0e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("expected album ID, got %s", result.AlbumID)
	}
	if result.TrackID != "ee0e8400-e29b-41d4-a716-446655440004" {
		t.Errorf("expected track ID, got %s", result.TrackID)
	}
}

func TestFinalize_StorageCopyFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	store := &mockStorage{
		copyFn: func(ctx context.Context, srcKey, dstKey string) error {
			return errors.New("storage unavailable")
		},
	}
	s := NewService(sqlxDB, store, zap.NewNop())

	meta := `{
		"artist": {"action": "link", "existingId": "550e8400-e29b-41d4-a716-446655440000", "name": "Test"},
		"album": {"action": "skip"},
		"track": {"title": "Test", "durationSeconds": 100, "trackNumber": 1}
	}`

	mock.ExpectBegin()
	mock.ExpectRollback()

	_, err = s.Finalize(context.Background(), "draft-3", "accepted", "ingestion-audio/file.mp3", "mp3", meta)
	if err == nil || !errorContains(err, "storage unavailable") {
		t.Fatalf("expected storage error, got: %v", err)
	}
}

func TestFinalize_DBInsertFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	store := &mockStorage{}
	s := NewService(sqlxDB, store, zap.NewNop())

	meta := `{
		"artist": {"action": "create", "name": "Fail Artist"},
		"album": {"action": "skip"},
		"track": {"title": "Fail Track", "durationSeconds": 100}
	}`

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO artists`).
		WillReturnError(errors.New("duplicate slug"))

	mock.ExpectRollback()

	_, err = s.Finalize(context.Background(), "draft-4", "accepted", "ingestion-audio/file.mp3", "mp3", meta)
	if err == nil {
		t.Fatal("expected error on insert failure")
	}
}

func errorContains(err error, s string) bool {
	if err == nil {
		return false
	}
	return len(err.Error()) >= len(s) && containsStr(err.Error(), s)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
