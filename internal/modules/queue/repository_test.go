package queue

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

func TestRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "track_id", "position", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), userID, uuid.New(), 1, now, now,
	).AddRow(
		uuid.New(), userID, uuid.New(), 2, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(userID).
		WillReturnRows(rows)

	items, err := repo.List(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddLater(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	trackID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`SELECT pg_advisory_xact_lock`)).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS`)).WithArgs(trackID).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE`)).WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO queue_items`)).
		WithArgs(userID, trackID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "track_id", "position", "created_at", "updated_at",
		}).AddRow(uuid.New(), userID, trackID, 1, time.Now(), time.Now()))
	mock.ExpectCommit()

	item, err := repo.AddLater(context.Background(), userID, trackID)
	require.NoError(t, err)
	assert.Equal(t, userID, item.UserID)
	assert.Equal(t, trackID, item.TrackID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Clear(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM queue_items WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err = repo.Clear(context.Background(), userID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
