package ai

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type mockRepoClient struct {
	aiClient
	mock sqlmock.Sqlmock
	db   *sqlx.DB
}

type aiClient struct {
	embedding []float64
	embedErr  error
	mood      *TrackMood
	moodErr   error
	playlist  []string
	playErr   error
}

func (m *aiClient) GenerateEmbedding(ctx context.Context, input string) ([]float64, error) {
	return m.embedding, m.embedErr
}
func (m *aiClient) AnalyzeMood(ctx context.Context, title, artist, genre string) (*TrackMood, error) {
	return m.mood, m.moodErr
}
func (m *aiClient) GeneratePlaylist(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, error) {
	return m.playlist, m.playErr
}

func TestService_GenerateEmbedding(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("test-id", "Test", "Artist", "Album", "Pop", 200, ""))

	mock.ExpectExec(`INSERT INTO track_embeddings_text`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "v1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{embedding: []float64{0.1, 0.2, 0.3}}, zap.NewNop(), true)

	result, err := svc.GenerateEmbedding(context.Background(), "test-id")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Embedding) != 3 {
		t.Fatalf("expected 3 dims, got %d", len(result.Embedding))
	}
}

func TestService_GenerateEmbedding_Disabled(t *testing.T) {
	svc := NewService(nil, nil, zap.NewNop(), false)
	_, err := svc.GenerateEmbedding(context.Background(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error when disabled")
	}
}

func TestService_GeneratePlaylist_FallbackSelect(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, "").
			AddRow("2", "B", "Artist", "Album", "Rock", 180, ""))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, ""))

	mock.ExpectExec(`INSERT INTO ai_generation_log`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "ai-v1", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{playErr: errors.New("AI failed")}, zap.NewNop(), true)

	req := GeneratePlaylistRequest{Prompt: "test", Limit: 1}
	result, err := svc.GeneratePlaylist(context.Background(), req, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Tracks) == 0 {
		t.Fatal("expected fallback tracks")
	}
}

func TestService_GeneratePlaylist_AllFallbackLevels(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, "").
			AddRow("2", "B", "Artist", "Album", "Rock", 180, ""))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, ""))

	mock.ExpectExec(`INSERT INTO ai_generation_log`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "ai-v1", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{playErr: errors.New("AI failed")}, zap.NewNop(), true)

	req := GeneratePlaylistRequest{Prompt: "test", SeedTrackID: "seed-1", Mood: "calm", Limit: 1}
	result, err := svc.GeneratePlaylist(context.Background(), req, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Tracks) == 0 {
		t.Fatal("expected tracks from fallback")
	}
}

func TestService_GetSimilarByMood_ByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id, t.title FROM tracks t`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("1", "A"))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	results, err := svc.GetSimilarByMood(context.Background(), "track-1", "calm", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestService_GetSimilarByMood_ByTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT track_id, mood_tags, energy, valence, tempo`)).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "mood_tags", "energy", "valence", "tempo", "danceability", "acousticness", "instrumentalness", "liveness", "speechiness", "updated_at"}).
			AddRow(uuid.New(), `[]`, 0.5, 0.5, 120.0, 0.5, 0.3, 0.1, 0.3, 0.1, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id, t.title FROM tracks t`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("2", "B"))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	results, err := svc.GetSimilarByMood(context.Background(), uuid.New().String(), "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestService_GetSimilarByEmbedding(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT track_id, embedding, model_version, updated_at FROM track_embeddings_text`)).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "embedding", "model_version", "updated_at"}).
			AddRow(uuid.New(), "{0.1,0.2}", "v1", nil))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_embeddings_audio te`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, "").
			AddRow("2", "B", "Artist", "Album", "Rock", 180, ""))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	results, err := svc.GetSimilarByEmbedding(context.Background(), uuid.New().String(), 5, "audio")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRedactPII_Email(t *testing.T) {
	result := redactPII("contact me at test@example.com for info")
	if result != "contact me at [EMAIL REDACTED] for info" {
		t.Fatalf("unexpected: %s", result)
	}
}

func TestRedactPII_Phone(t *testing.T) {
	result := redactPII("my phone is 09123456789")
	if result != "my phone is [PHONE REDACTED]" {
		t.Fatalf("unexpected: %s", result)
	}
}

func TestRedactPII_Empty(t *testing.T) {
	result := redactPII("")
	if result != "" {
		t.Fatal("expected empty")
	}
}

func TestRedactPII_Clean(t *testing.T) {
	result := redactPII("play some chill music")
	if result != "play some chill music" {
		t.Fatalf("unexpected: %s", result)
	}
}
