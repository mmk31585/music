package subscription

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPlans(ctx context.Context) ([]Plan, error) {
	var plans []Plan

	err := r.db.SelectContext(ctx, &plans, `
		SELECT
			id,
			code,
			name,
			description,
			price_cents,
			currency,
			interval,
			is_premium,
			is_active,
			features,
			created_at,
			updated_at
		FROM plans
		ORDER BY price_cents ASC, created_at ASC
	`)

	return plans, err
}

func (r *Repository) GetPlanByID(ctx context.Context, planID string) (*Plan, error) {
	var plan Plan

	err := r.db.GetContext(ctx, &plan, `
		SELECT
			id,
			code,
			name,
			description,
			price_cents,
			currency,
			interval,
			is_premium,
			is_active,
			features,
			created_at,
			updated_at
		FROM plans
		WHERE id = $1
	`, planID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (r *Repository) GetPlanByCode(ctx context.Context, code string) (*Plan, error) {
	var plan Plan

	err := r.db.GetContext(ctx, &plan, `
		SELECT
			id,
			code,
			name,
			description,
			price_cents,
			currency,
			interval,
			is_premium,
			is_active,
			features,
			created_at,
			updated_at
		FROM plans
		WHERE code = $1
	`, code)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (r *Repository) GetCurrentSubscription(ctx context.Context, userID string) (*Subscription, error) {
	var sub Subscription

	err := r.db.GetContext(ctx, &sub, `
		SELECT
			id,
			user_id,
			plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			canceled_at,
			provider,
			provider_subscription_id,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
		  AND status IN ('active', 'trialing', 'pending')
		ORDER BY created_at DESC
		LIMIT 1
	`, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *Repository) ListSubscriptions(ctx context.Context, userID string) ([]Subscription, error) {
	var subs []Subscription

	err := r.db.SelectContext(ctx, &subs, `
		SELECT
			id,
			user_id,
			plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			canceled_at,
			provider,
			provider_subscription_id,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	return subs, err
}

func (r *Repository) CreateSubscription(ctx context.Context, userID string, planID string, status string) (*Subscription, error) {
	var sub Subscription

	err := r.db.GetContext(ctx, &sub, `
		INSERT INTO subscriptions (
			user_id,
			plan_id,
			status,
			current_period_start,
			current_period_end
		)
		VALUES (
			$1,
			$2,
			$3,
			NOW(),
			NOW() + INTERVAL '1 month'
		)
		RETURNING
			id,
			user_id,
			plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			canceled_at,
			provider,
			provider_subscription_id,
			created_at,
			updated_at
	`, userID, planID, status)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *Repository) CancelCurrentSubscription(ctx context.Context, userID string) (*Subscription, error) {
	var sub Subscription

	err := r.db.GetContext(ctx, &sub, `
		UPDATE subscriptions
		SET
			status = 'canceled',
			cancel_at_period_end = TRUE,
			canceled_at = NOW(),
			updated_at = NOW()
		WHERE id = (
			SELECT id
			FROM subscriptions
			WHERE user_id = $1
			  AND status IN ('active', 'trialing', 'pending')
			ORDER BY created_at DESC
			LIMIT 1
		)
		RETURNING
			id,
			user_id,
			plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			canceled_at,
			provider,
			provider_subscription_id,
			created_at,
			updated_at
	`, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *Repository) CreatePayment(ctx context.Context, userID string, subID string, planID string, amountCents int64, currency string, status string) (*Payment, error) {
	var payment Payment

	err := r.db.GetContext(ctx, &payment, `
		INSERT INTO payments (
			user_id,
			subscription_id,
			plan_id,
			amount_cents,
			currency,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			user_id,
			subscription_id,
			plan_id,
			amount_cents,
			currency,
			status,
			provider,
			provider_payment_id,
			failure_reason,
			paid_at,
			created_at,
			updated_at
	`, userID, subID, planID, amountCents, currency, status)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *Repository) ListPayments(ctx context.Context, userID string) ([]Payment, error) {
	var payments []Payment

	err := r.db.SelectContext(ctx, &payments, `
		SELECT
			id,
			user_id,
			subscription_id,
			plan_id,
			amount_cents,
			currency,
			status,
			provider,
			provider_payment_id,
			failure_reason,
			paid_at,
			created_at,
			updated_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	return payments, err
}
