package ai

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func TestUpsertEmbedding(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectExec(`INSERT INTO track_embeddings_text`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "v1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	err = repo.UpsertEmbedding(context.Background(), uuid.New(), []float64{0.1, 0.2, 0.3}, "v1")
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetEmbeddingNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(`SELECT track_id, embedding, model_version, updated_at FROM track_embeddings_text`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	repo := NewRepository(sqlxDB)
	_, err = repo.GetEmbedding(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetSimilarByEmbeddingDefaultTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	rows := sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
		AddRow("1", "Track 1", "Artist 1", "Album 1", "Pop", 200, "")

	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_embeddings_audio te`)).
		WithArgs(sqlmock.AnyArg(), 10).
		WillReturnRows(rows)

	repo := NewRepository(sqlxDB)
	results, err := repo.GetSimilarByEmbedding(context.Background(), []float64{0.1, 0.2}, 10, "audio")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestGetSimilarByEmbeddingTextTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	rows := sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
		AddRow("2", "Track 2", "Artist 2", "Album 2", "Rock", 180, "")

	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_embeddings_text te`)).
		WithArgs(sqlmock.AnyArg(), 10).
		WillReturnRows(rows)

	repo := NewRepository(sqlxDB)
	results, err := repo.GetSimilarByEmbedding(context.Background(), []float64{0.1, 0.2}, 10, "text")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestUpsertMood(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectExec(`INSERT INTO track_moods`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 0.5, 0.5, 120.0, 0.5, 0.3, 0.1, 0.3, 0.1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	mood := TrackMood{
		TrackID:          uuid.New(),
		Energy:           0.5,
		Valence:          0.5,
		Tempo:            120,
		Danceability:     0.5,
		Acousticness:     0.3,
		Instrumentalness: 0.1,
		Liveness:         0.3,
		Speechiness:      0.1,
	}
	err = repo.UpsertMood(context.Background(), mood)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogGeneration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectExec(`INSERT INTO ai_generation_log`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 5, "test-model", 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	pid := uuid.New()
	err = repo.LogGeneration(context.Background(), GenerationLog{
		UserID:     uuid.New(),
		PlaylistID: &pid,
		Prompt:     "test prompt",
		TrackCount: 5,
		ModelUsed:  "test-model",
		LatencyMs:  100,
	})
	if err != nil {
		t.Fatal(err)
	}
}
