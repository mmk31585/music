package library

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_LikeTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	trackID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO liked_tracks`)).
		WithArgs(userID.String(), trackID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.LikeTrack(context.Background(), userID.String(), trackID.String())
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UnlikeTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	trackID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM liked_tracks WHERE user_id = $1 AND track_id = $2`)).
		WithArgs(userID.String(), trackID.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UnlikeTrack(context.Background(), userID.String(), trackID.String())
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
