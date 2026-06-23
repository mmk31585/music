package track

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

func TestRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	trackID := uuid.New()
	albumID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "artist", "album_id", "duration", "track_number",
		"disc_number", "status", "is_explicit", "language", "created_at", "updated_at",
	}).AddRow(
		trackID, "Test Track", "Test Artist", albumID, 240, 1, 1,
		"published", false, "fa", now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(trackID).
		WillReturnRows(rows)

	track, err := repo.GetByID(context.Background(), trackID)
	require.NoError(t, err)
	assert.Equal(t, trackID, track.ID)
	assert.Equal(t, "Test Track", track.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_ListPublic(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "artist", "album_id", "duration", "track_number",
		"disc_number", "status", "is_explicit", "language", "play_count",
		"created_at", "updated_at",
	}).AddRow(
		uuid.New(), "Track 1", "Artist 1", nil, 200, 1, 1,
		"published", false, "en", 100, now, now,
	).AddRow(
		uuid.New(), "Track 2", "Artist 2", nil, 180, 1, 1,
		"published", true, "fa", 50, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(20, 0).
		WillReturnRows(rows)

	tracks, err := repo.List(context.Background(), 20, 0, true, ListOptions{})
	require.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	trackID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"id", "title", "duration", "album_id", "track_number", "disc_number",
		"is_explicit", "language", "status", "created_at", "updated_at",
	}).AddRow(
		trackID, "New Track", 240, nil, 1, 1,
		false, "en", "draft", now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tracks`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	explicit := false
	trackNum := 1
	albumID := uuid.Nil
	track, err := repo.Create(context.Background(), CreateRequest{
		Title:           "New Track",
		DurationSeconds: 240,
		Explicit:        &explicit,
		AlbumID:         &albumID,
		TrackNumber:     &trackNum,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Track", track.Title)
	assert.Equal(t, trackID, track.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	trackID := uuid.New()
	now := time.Now()

	// GetByID: SELECT from tracks
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "artist_id", "album_id", "title", "slug",
			"duration_seconds", "track_number", "explicit",
			"audio_url", "cover_url", "audio_media_id", "cover_media_id",
			"play_count", "is_public", "created_at", "updated_at",
		}).AddRow(trackID, nil, nil, "Original Title", "original-title",
			240, 1, false, nil, nil, nil, nil,
			0, true, now, now))

	// hydrate: listArtistsByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "slug", "role", "position"}))

	// hydrate: listGenresByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "created_at"}))

	// Begin transaction
	mock.ExpectBegin()

	// UPDATE with title and explicit fields
	title := "Updated Title"
	explicit := true
	mock.ExpectQuery(`UPDATE tracks`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "artist_id", "album_id", "title", "slug",
			"duration_seconds", "track_number", "explicit",
			"audio_url", "cover_url", "audio_media_id", "cover_media_id",
			"play_count", "is_public", "created_at", "updated_at",
		}).AddRow(trackID, nil, nil, title, "updated-title",
			240, 1, explicit, nil, nil, nil, nil,
			0, true, now, now))

	// Commit
	mock.ExpectCommit()

	_, err = repo.Update(context.Background(), trackID, UpdateRequest{
		Title:    &title,
		Explicit: &explicit,
	})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	trackID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM track_artists WHERE track_id = $1`)).
		WithArgs(trackID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tracks WHERE id = $1`)).
		WithArgs(trackID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(context.Background(), trackID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
