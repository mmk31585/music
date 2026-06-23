package media

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mediaID := uuid.New()
	now := time.Now()
	fileSize := int64(1024000)
	mimeType := "audio/mpeg"
	origFilename := "test.mp3"

	rows := sqlmock.NewRows([]string{
		"id", "media_type", "storage_provider", "bucket", "object_key", "public_url",
		"mime_type", "file_size", "checksum_sha256", "duration_seconds", "width", "height",
		"original_filename", "metadata", "created_by", "created_at", "updated_at",
	}).AddRow(
		mediaID, "audio", "local", nil, "/uploads/test.mp3", nil,
		&mimeType, &fileSize, nil, nil, nil, nil,
		&origFilename, "{}", nil, now, now,
	)

	mock.ExpectQuery(`INSERT INTO media`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := repo.Create(context.Background(), CreateMediaRequest{
		MediaType:       "audio",
		StorageProvider: "local",
		ObjectKey:       "/uploads/test.mp3",
		MimeType:        &mimeType,
		FileSize:        &fileSize,
	})
	require.NoError(t, err)
	assert.Equal(t, mediaID, result.ID)
	assert.Equal(t, &fileSize, result.FileSize)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mediaID := uuid.New()
	now := time.Now()
	mimeType := "audio/mpeg"
	fileSize := int64(2048000)
	origFilename := "track.mp3"
	rows := sqlmock.NewRows([]string{
		"id", "media_type", "storage_provider", "bucket", "object_key", "public_url",
		"mime_type", "file_size", "checksum_sha256", "duration_seconds", "width", "height",
		"original_filename", "metadata", "created_by", "created_at", "updated_at",
	}).AddRow(
		mediaID, "audio", "s3", nil, "/tracks/track.mp3", nil,
		&mimeType, &fileSize, nil, nil, nil, nil,
		&origFilename, `{"duration": 240}`, nil, now, now,
	)

	mock.ExpectQuery(`SELECT`).WithArgs(mediaID).WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), mediaID)
	require.NoError(t, err)
	assert.Equal(t, mediaID, result.ID)
	assert.Equal(t, "s3", result.StorageProvider)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mock.ExpectQuery(`SELECT`).WithArgs(sqlmock.AnyArg()).WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	mediaID := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM media WHERE id = $1`)).
		WithArgs(mediaID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), mediaID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
