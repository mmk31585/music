package subscription

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

func TestRepository_ListPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "code", "name", "description", "price_cents",
		"currency", "interval", "is_premium", "is_active", "features", "created_at", "updated_at",
	}).AddRow(
		uuid.New().String(), "premium", "Premium", "Premium plan", int64(999),
		"USD", "month", false, true, []byte(`{"max_bitrate": 320}`), now, now,
	).AddRow(
		uuid.New().String(), "family", "Family", "Family plan", int64(1499),
		"USD", "month", false, true, []byte(`{"max_bitrate": 320, "max_users": 6}`), now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)

	plans, err := repo.ListPlans(context.Background())
	require.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.Equal(t, "premium", plans[0].Code)
	assert.Equal(t, int64(999), plans[0].PriceCents)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetCurrentSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	subID := uuid.New()
	planID := uuid.New()
	now := time.Now()
	future := now.Add(30 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "plan_id", "status", "current_period_start",
		"current_period_end", "cancel_at_period_end", "created_at", "updated_at",
	}).AddRow(
		subID, userID, planID, "active", now, future, false, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(userID.String()).
		WillReturnRows(rows)

	sub, err := repo.GetCurrentSubscription(context.Background(), userID.String())
	require.NoError(t, err)
	assert.Equal(t, subID.String(), sub.ID)
	assert.Equal(t, "active", sub.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	subID := uuid.New()
	userID := uuid.New()
	planID := uuid.New()
	now := time.Now()
	future := now.Add(30 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "plan_id", "status", "current_period_start",
		"current_period_end", "cancel_at_period_end", "created_at", "updated_at",
	}).AddRow(
		subID, userID, planID, "active", now, future, false, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO subscriptions`)).
		WithArgs(userID.String(), planID.String(), "active").
		WillReturnRows(rows)

	sub, err := repo.CreateSubscription(context.Background(), userID.String(), planID.String(), "active")
	require.NoError(t, err)
	assert.Equal(t, "active", sub.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CancelCurrentSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewRepository(sqlxDB)

	userID := uuid.New()
	subID := uuid.New()
	planID := uuid.New()
	now := time.Now()
	future := now.Add(30 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "plan_id", "status", "current_period_start",
		"current_period_end", "cancel_at_period_end", "canceled_at",
		"provider", "provider_subscription_id", "created_at", "updated_at",
	}).AddRow(
		subID.String(), userID.String(), planID.String(), "canceled",
		now, future, true, now, nil, nil, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE subscriptions`)).
		WithArgs(userID.String()).
		WillReturnRows(rows)

	sub, err := repo.CancelCurrentSubscription(context.Background(), userID.String())
	require.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, "canceled", sub.Status)
	assert.Equal(t, subID.String(), sub.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
