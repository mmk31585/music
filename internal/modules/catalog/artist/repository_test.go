package artist

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

	artistID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id",
		"country", "is_verified", "monthly_listeners", "created_at", "updated_at",
	}).AddRow(
		artistID, "Test Artist", "test-artist", nil, nil, nil, nil,
		nil, true, int64(10000), now, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM artists WHERE id = $1`)).WithArgs(artistID).
		WillReturnRows(rows)

	artist, err := repo.GetByID(context.Background(), artistID)
	require.NoError(t, err)
	assert.Equal(t, artistID, artist.ID)
	assert.Equal(t, "Test Artist", artist.Name)
	assert.Equal(t, int64(10000), artist.MonthlyListeners)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	artistID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM artists WHERE id = $1`)).WithArgs(artistID).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetByID(context.Background(), artistID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id",
		"country", "is_verified", "monthly_listeners", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), "Artist A", "artist-a", nil, nil, nil, nil,
		nil, true, int64(5000), now, nil,
	).AddRow(
		uuid.New(), "Artist B", "artist-b", nil, nil, nil, nil,
		nil, true, int64(3000), now, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM artists ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).WithArgs(20, 0).
		WillReturnRows(rows)

	artists, err := repo.List(context.Background(), 20, 0)
	require.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, "Artist A", artists[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	artistID := uuid.New()
	now := time.Now()

	bio := "New bio"

	rows := sqlmock.NewRows([]string{
		"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id",
		"country", "is_verified", "monthly_listeners", "created_at", "updated_at",
	}).AddRow(
		artistID, "New Artist", "new-artist", &bio, nil, nil, nil,
		nil, false, int64(0), now, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO artists`)).
		WithArgs("New Artist", "new-artist", &bio, nil, nil, nil, nil, nil, nil).
		WillReturnRows(rows)

	artist, err := repo.Create(context.Background(), CreateRequest{Name: "New Artist", Bio: &bio})
	require.NoError(t, err)
	assert.Equal(t, "New Artist", artist.Name)
	assert.Equal(t, artistID, artist.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	artistID := uuid.New()
	now := time.Now()

	name := "Updated Name"
	bio := "Updated bio"

	// GetByID: SELECT before UPDATE
	mock.ExpectQuery(regexp.QuoteMeta(`FROM artists WHERE id = $1`)).WithArgs(artistID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id",
			"country", "is_verified", "monthly_listeners", "created_at", "updated_at",
		}).AddRow(
			artistID, "Original", "original", nil, nil, nil, nil,
			nil, false, int64(0), now, nil,
		))

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE artists SET`)).
		WithArgs(artistID, &name, "updated-name", &bio, nil, nil, nil, nil, nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id", "country", "is_verified", "monthly_listeners", "created_at", "updated_at"}).
			AddRow(artistID, "Updated Name", "updated-name", &bio, nil, nil, nil, nil, false, int64(0), now, nil))

	_, err = repo.Update(context.Background(), artistID, UpdateRequest{
		Name: &name,
		Bio:  &bio,
	})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
