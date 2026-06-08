package subscription

import (
	"encoding/json"
	"time"
)

type Plan struct {
	ID          string          `db:"id"`
	Code        string          `db:"code"`
	Name        string          `db:"name"`
	Description *string         `db:"description"`
	PriceCents  int64           `db:"price_cents"`
	Currency    string          `db:"currency"`
	Interval    string          `db:"interval"`
	IsPremium   bool            `db:"is_premium"`
	IsActive    bool            `db:"is_active"`
	Features    json.RawMessage `db:"features"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

type Subscription struct {
	ID                 string     `db:"id"`
	UserID             string     `db:"user_id"`
	PlanID             string     `db:"plan_id"`
	Status             string     `db:"status"`
	CurrentPeriodStart time.Time  `db:"current_period_start"`
	CurrentPeriodEnd   *time.Time `db:"current_period_end"`
	CancelAtPeriodEnd  bool       `db:"cancel_at_period_end"`
	CanceledAt         *time.Time `db:"canceled_at"`
	Provider           *string    `db:"provider"`
	ProviderSubID      *string    `db:"provider_subscription_id"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

type Payment struct {
	ID             string     `db:"id"`
	UserID         string     `db:"user_id"`
	SubscriptionID *string    `db:"subscription_id"`
	PlanID         *string    `db:"plan_id"`
	AmountCents    int64      `db:"amount_cents"`
	Currency       string     `db:"currency"`
	Status         string     `db:"status"`
	Provider       *string    `db:"provider"`
	ProviderPayID  *string    `db:"provider_payment_id"`
	FailureReason  *string    `db:"failure_reason"`
	PaidAt         *time.Time `db:"paid_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
