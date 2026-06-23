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
		"id", "name", "bio", "avatar_url", "cover_url",
		"monthly_listeners", "status", "created_at", "updated_at",
	}).AddRow(
		artistID, "Test Artist", "A test artist", nil, nil,
		10000, "active", now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(artistID).
		WillReturnRows(rows)

	artist, err := repo.GetByID(context.Background(), artistID)
	require.NoError(t, err)
	assert.Equal(t, artistID, artist.ID)
	assert.Equal(t, "Test Artist", artist.Name)
	assert.Equal(t, 10000, artist.MonthlyListeners)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	artistID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(artistID).
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

	rows := sqlmock.NewRows([]string{
		"id", "name", "bio", "avatar_url", "cover_url",
		"monthly_listeners", "status", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), "Artist A", "Bio A", nil, nil, 5000, "active", time.Now(), time.Now(),
	).AddRow(
		uuid.New(), "Artist B", "Bio B", nil, nil, 3000, "active", time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(20, 0).
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

	rows := sqlmock.NewRows([]string{
		"id", "name", "bio", "avatar_url", "cover_url",
		"monthly_listeners", "status", "created_at", "updated_at",
	}).AddRow(
		artistID, "New Artist", "New bio", nil, nil, 0, "active", now, now,
	)

	bio := "New bio"
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO artists`)).
		WithArgs("New Artist", &bio, sqlmock.AnyArg()).
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
	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE artists SET`)).
		WithArgs(artistID, "Updated Name", "updated-name", "Updated bio", nil, nil, nil, nil, nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "bio", "image_url", "avatar_media_id", "banner_media_id", "country", "is_verified", "monthly_listeners", "created_at", "updated_at"}).
			AddRow(artistID, "Updated Name", "updated-name", "Updated bio", nil, nil, nil, nil, nil, nil, time.Now(), time.Now()))

	name := "Updated Name"
	bio := "Updated bio"
	_, err = repo.Update(context.Background(), artistID, UpdateRequest{
		Name: &name,
		Bio:  &bio,
	})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
