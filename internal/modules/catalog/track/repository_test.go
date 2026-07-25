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
	artistID := uuid.New()
	albumID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "artist_id", "album_id", "title", "slug", "duration_seconds", "track_number",
		"explicit", "audio_url", "cover_url", "audio_media_id", "cover_media_id",
		"play_count", "is_public", "created_at", "updated_at",
	}).AddRow(
		trackID, artistID, &albumID, "Test Track", "test-track", 240, 1,
		false, nil, nil, nil, nil,
		int64(0), true, now, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM tracks WHERE id = $1`)).WithArgs(trackID).
		WillReturnRows(rows)

	// hydrate: listArtistsByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_artists ta`)).WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "slug", "role", "position"}))

	// hydrate: listGenresByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`FROM genres g`)).WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "created_at"}))

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

	track1ID := uuid.New()
	track2ID := uuid.New()
	artistID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "artist_id", "album_id", "title", "slug", "duration_seconds", "track_number",
		"explicit", "audio_url", "cover_url", "audio_media_id", "cover_media_id",
		"play_count", "is_public", "created_at", "updated_at",
	}).AddRow(
		track1ID, artistID, nil, "Track 1", "track-1", 200, 1,
		false, nil, nil, nil, nil,
		int64(100), true, now, nil,
	).AddRow(
		track2ID, artistID, nil, "Track 2", "track-2", 180, 1,
		true, nil, nil, nil, nil,
		int64(50), true, now, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM tracks t WHERE t.is_public = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).WithArgs(20, 0).
		WillReturnRows(rows)

	// hydrate for each track: listArtistsByTrack + listGenresByTrack
	for i := 0; i < 2; i++ {
		mock.ExpectQuery(regexp.QuoteMeta(`FROM track_artists ta`)).WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "slug", "role", "position"}))
		mock.ExpectQuery(regexp.QuoteMeta(`FROM genres g`)).WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "created_at"}))
	}

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
	artistID := uuid.New()
	now := time.Now()

	explicit := false
	trackNum := 1
	albumID := uuid.Nil

	// Begin transaction
	mock.ExpectBegin()

	// generateUniqueSlug
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT slug FROM tracks WHERE artist_id = $1 AND slug LIKE $2`)).
		WithArgs(artistID, "new-track%").
		WillReturnRows(sqlmock.NewRows([]string{"slug"}))

	// INSERT INTO tracks
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tracks`)).
		WithArgs(
			artistID, &albumID, "New Track", "new-track", 240, &trackNum,
			false, nil, nil, nil, nil, true,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "artist_id", "album_id", "title", "slug", "duration_seconds", "track_number",
			"explicit", "audio_url", "cover_url", "audio_media_id", "cover_media_id",
			"play_count", "is_public", "created_at", "updated_at",
		}).AddRow(
			trackID, artistID, &albumID, "New Track", "new-track", 240, 1,
			false, nil, nil, nil, nil,
			int64(0), true, now, nil,
		))

	// replaceArtistsTx: DELETE old artists
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM track_artists WHERE track_id = $1`)).
		WithArgs(trackID).WillReturnResult(sqlmock.NewResult(0, 0))

	// replaceArtistsTx: INSERT new artist
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO track_artists`)).
		WithArgs(trackID, artistID, "primary", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// replaceGenresTx: DELETE old genres (no genres, but DELETE always runs)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM track_genres WHERE track_id = $1`)).
		WithArgs(trackID).WillReturnResult(sqlmock.NewResult(0, 0))

	// Commit
	mock.ExpectCommit()

	// GetByID reload: SELECT from tracks
	mock.ExpectQuery(regexp.QuoteMeta(`FROM tracks WHERE id = $1`)).WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "artist_id", "album_id", "title", "slug", "duration_seconds", "track_number",
			"explicit", "audio_url", "cover_url", "audio_media_id", "cover_media_id",
			"play_count", "is_public", "created_at", "updated_at",
		}).AddRow(
			trackID, artistID, &albumID, "New Track", "new-track", 240, 1,
			false, nil, nil, nil, nil,
			int64(0), true, now, nil,
		))

	// hydrate: listArtistsByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_artists ta`)).WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "slug", "role", "position"}).
			AddRow(artistID, "Test Artist", "test-artist", "primary", 1))

	// hydrate: listGenresByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`FROM genres g`)).WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "created_at"}))

	track, err := repo.Create(context.Background(), CreateRequest{
		Title:           "New Track",
		DurationSeconds: 240,
		Explicit:        &explicit,
		AlbumID:         &albumID,
		TrackNumber:     &trackNum,
		Artists: []TrackArtistRequest{
			{ArtistID: artistID, Role: "primary", Position: 1},
		},
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

	// generateUniqueSlug: check for existing slugs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT slug FROM tracks WHERE artist_id = $1 AND slug LIKE $2 AND id != $3`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"slug"}))

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

	// GetByID reload after update
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "artist_id", "album_id", "title", "slug",
			"duration_seconds", "track_number", "explicit",
			"audio_url", "cover_url", "audio_media_id", "cover_media_id",
			"play_count", "is_public", "created_at", "updated_at",
		}).AddRow(trackID, nil, nil, title, "updated-title",
			240, 1, explicit, nil, nil, nil, nil,
			0, true, now, now))

	// hydrate: listArtistsByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"artist_id", "name", "slug", "role", "position"}))

	// hydrate: listGenresByTrack
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(trackID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "created_at"}))

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

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tracks WHERE id = $1`)).
		WithArgs(trackID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), trackID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
