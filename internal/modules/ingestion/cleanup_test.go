package ingestion

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func TestCleanupService_FlagStaleDrafts(t *testing.T) {
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
		nil, nil, "mp3", "", false, "pending",
		`{}`, nil, nil,
		now.Add(-48*time.Hour), now.Add(-48*time.Hour),
	).AddRow(
		"draft-2", "user-2", "test2.mp3", "/path/file2.mp3", int64(2000),
		nil, nil, "mp3", "", false, "enriching",
		`{}`, nil, nil,
		now.Add(-36*time.Hour), now.Add(-36*time.Hour),
	)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE status IN .+ AND updated_at < \$1 AND stale = FALSE ORDER BY updated_at ASC`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE ingestion_drafts SET stale = TRUE, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("draft-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE ingestion_drafts SET stale = TRUE, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("draft-2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	svc := NewCleanupService(repo, zap.NewNop())
	flagged, err := svc.FlagStaleDrafts(context.Background(), 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if flagged != 2 {
		t.Errorf("expected 2 flagged, got %d", flagged)
	}
}

func TestCleanupService_FlagStaleDrafts_None(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE status IN .+ AND updated_at < \$1 AND stale = FALSE ORDER BY updated_at ASC`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "uploaded_by", "original_filename", "file_path", "file_size",
			"duration_seconds", "bitrate", "format", "file_hash", "stale", "status",
			"extracted_metadata", "enriched_metadata", "final_metadata",
			"created_at", "updated_at",
		}))

	svc := NewCleanupService(repo, zap.NewNop())
	flagged, err := svc.FlagStaleDrafts(context.Background(), 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if flagged != 0 {
		t.Errorf("expected 0 flagged, got %d", flagged)
	}
}

func TestCleanupService_FlagStaleDrafts_MarkError(t *testing.T) {
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
		nil, nil, "mp3", "", false, "pending",
		`{}`, nil, nil,
		now.Add(-48*time.Hour), now.Add(-48*time.Hour),
	)

	mock.ExpectQuery(`SELECT .+ FROM ingestion_drafts WHERE status IN .+ AND updated_at < \$1 AND stale = FALSE ORDER BY updated_at ASC`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE ingestion_drafts SET stale = TRUE, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("draft-1").
		WillReturnError(errors.New("db error"))

	svc := NewCleanupService(repo, zap.NewNop())
	flagged, err := svc.FlagStaleDrafts(context.Background(), 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if flagged != 0 {
		t.Errorf("expected 0 flagged (due to error), got %d", flagged)
	}
}
