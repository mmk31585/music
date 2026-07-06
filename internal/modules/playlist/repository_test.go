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
		"id", "user_id", "owner_id", "name", "description", "cover_url", "is_public", "created_at", "updated_at",
	}).AddRow(
		playlistID, userID, userID, "Test Playlist", nil, nil, true, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO playlists (user_id, owner_id, name, description, cover_url, is_public)
		VALUES ($1, $1, $2, $3, $4, $5)
		RETURNING id, user_id, owner_id, name, description, cover_url, is_public, created_at, updated_at
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
		"id", "user_id", "owner_id", "name", "description", "cover_url", "is_public",
		"is_collaborative", "created_at", "updated_at",
	}).AddRow(
		playlistID, userID, userID, "My Playlist", nil, nil, true,
		false, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, owner_id, name, description, cover_url, is_public, is_collaborative, created_at, updated_at
		FROM playlists
		WHERE id = $1
	`)).WithArgs(playlistID).
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
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM tracks WHERE id = $1)`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(MAX(position), 0) + 1 FROM playlist_tracks WHERE playlist_id = $1`)).
		WithArgs(playlistID).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES ($1, $2, $3)`)).
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
	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM playlists WHERE id = $1 AND user_id = $2`)).
		WithArgs(playlistID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeletePlaylist(context.Background(), playlistID, userID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func strPtr(s string) *string { return &s }
