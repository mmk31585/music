package ingestion

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSearchByFileHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{
		"id", "uploaded_by", "original_filename", "file_path", "file_size",
		"duration_seconds", "bitrate", "format", "file_hash", "stale", "status",
		"extracted_metadata", "enriched_metadata", "final_metadata",
		"created_at", "updated_at",
	}).AddRow(
		"draft-1", "user-1", "test.mp3", "/path/file.mp3", int64(1000),
		nil, nil, "mp3", "abc123", false, "published",
		`{}`, nil, nil,
		time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE file_hash = \$1 ORDER BY created_at DESC`).
		WithArgs("abc123").
		WillReturnRows(rows)

	result, err := repo.SearchByFileHash(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].ID != "draft-1" || result[0].Status != DraftStatusPublished {
		t.Errorf("unexpected first result: %+v", result[0])
	}
}

func TestSearchByFileHash_NoResults(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE file_hash = \$1 ORDER BY created_at DESC`).
		WithArgs("nonexistent").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "uploaded_by", "original_filename", "file_path", "file_size",
			"duration_seconds", "bitrate", "format", "file_hash", "stale", "status",
			"extracted_metadata", "enriched_metadata", "final_metadata",
			"created_at", "updated_at",
		}))

	result, err := repo.SearchByFileHash(context.Background(), "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 results, got %d", len(result))
	}
}

func TestListStaleDrafts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "uploaded_by", "original_filename", "file_path", "file_size",
		"duration_seconds", "bitrate", "format", "file_hash", "stale", "status",
		"extracted_metadata", "enriched_metadata", "final_metadata",
		"created_at", "updated_at",
	}).AddRow(
		"draft-1", "user-1", "test.mp3", "/path/file.mp3", int64(1000),
		nil, nil, "mp3", "", true, "pending",
		`{}`, nil, nil,
		now.Add(-48*time.Hour), now.Add(-48*time.Hour),
	)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE status IN .+ AND updated_at < \$1 AND stale = FALSE ORDER BY updated_at ASC`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	cutoff := now.Add(-24 * time.Hour)
	result, err := repo.ListStaleDrafts(context.Background(), cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestMarkDraftStale(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mock.ExpectExec(`UPDATE ingestion_drafts SET stale = TRUE, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("draft-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkDraftStale(context.Background(), "draft-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCountDraftsByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"status", "cnt"}).
		AddRow("pending", 5).
		AddRow("review", 3).
		AddRow("published", 10)

	mock.ExpectQuery(`SELECT status, COUNT\(\*\) AS cnt FROM ingestion_drafts GROUP BY status`).
		WillReturnRows(rows)

	result, err := repo.CountDraftsByStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 statuses, got %d", len(result))
	}
	if result[DraftStatusPending] != 5 {
		t.Errorf("expected pending=5, got %d", result[DraftStatusPending])
	}
	if result[DraftStatusReview] != 3 {
		t.Errorf("expected review=3, got %d", result[DraftStatusReview])
	}
}

func TestCountPublishedThisMonth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM ingestion_drafts WHERE status = 'published' AND updated_at >= date_trunc\('month', NOW\(\)\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))

	result, err := repo.CountPublishedThisMonth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result != 7 {
		t.Errorf("expected 7, got %d", result)
	}
}
