package notification

import (
	"context"
	"encoding/json"
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

	userID := uuid.New()
	notifID := uuid.New()
	now := time.Now()

	payload := map[string]interface{}{"type": "like", "actor_id": "user-2"}
	payloadBytes, _ := json.Marshal(payload)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "title", "body", "payload",
		"is_read", "created_at",
	}).AddRow(
		notifID, userID, "like", "New like", "Someone liked your track",
		string(payloadBytes), false, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO notifications`)).
		WithArgs(userID.String(), "like", "New like", "Someone liked your track", string(payloadBytes)).
		WillReturnRows(rows)

	input := CreateNotificationInput{
		UserID:  userID.String(),
		Type:    "like",
		Title:   "New like",
		Body:    "Someone liked your track",
		Payload: payload,
	}

	n, err := repo.Create(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, notifID, n.ID)
	assert.Equal(t, "like", n.Type)
	assert.False(t, n.IsRead)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_ListByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "title", "body", "payload", "is_read", "created_at",
	}).AddRow(
		uuid.New(), userID, "follow", "New follower", nil, "{}", false, now,
	).AddRow(
		uuid.New(), userID, "milestone", "Level up!", nil, "{}", true, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(userID.String(), 20, 0).
		WillReturnRows(rows)

	notifs, err := repo.ListByUser(context.Background(), userID.String(), 20, 0)
	require.NoError(t, err)
	assert.Len(t, notifs, 2)
	assert.Equal(t, "follow", notifs[0].Type)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_MarkAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	notifID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE notifications`)).
		WithArgs(notifID.String(), userID.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkAsRead(context.Background(), userID.String(), notifID.String())
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CountUnreadByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`)).
		WithArgs(userID.String()).
		WillReturnRows(rows)

	count, err := repo.CountUnreadByUser(context.Background(), userID.String())
	require.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
