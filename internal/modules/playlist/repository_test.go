package playlist

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

func TestRepository_CreatePlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	playlistID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "description", "cover_url", "is_public", "created_at", "updated_at",
	}).AddRow(
		playlistID, userID, "Test Playlist", nil, nil, true, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO playlists (user_id, name, description, cover_url, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, description, cover_url, is_public, created_at, updated_at
	`)).WithArgs(userID, "Test Playlist", sqlmock.AnyArg(), sqlmock.AnyArg(), true).
		WillReturnRows(rows)

	p, err := repo.CreatePlaylist(context.Background(), CreatePlaylistRequest{
		Name:        "Test Playlist",
		Description: strPtr("A test playlist"),
		IsPublic:    true,
	}, userID)

	require.NoError(t, err)
	assert.Equal(t, playlistID, p.ID)
	assert.Equal(t, "Test Playlist", p.Name)
	assert.True(t, p.IsPublic)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetPlaylistByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	playlistID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "description", "cover_url", "is_public",
		"track_count", "duration_seconds", "created_at", "updated_at",
	}).AddRow(
		playlistID, userID, "My Playlist", nil, nil, true,
		5, 1200, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(playlistID).
		WillReturnRows(rows)

	_, err = repo.GetPlaylistByID(context.Background(), playlistID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	playlistID := uuid.New()
	trackID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(MAX(position), 0) + 1 FROM playlist_tracks WHERE playlist_id = $1`)).
		WithArgs(playlistID).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO playlist_tracks`)).
		WithArgs(playlistID, trackID, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.AddTrack(context.Background(), playlistID, trackID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	playlistID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM playlist_tracks WHERE playlist_id = $1`)).
		WithArgs(playlistID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM playlists WHERE id = $1`)).
		WithArgs(playlistID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeletePlaylist(context.Background(), playlistID, uuid.New())
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func strPtr(s string) *string { return &s }
