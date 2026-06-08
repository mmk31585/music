package subscription

import "time"

type PlanResponse struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	PriceCents  int64    `json:"priceCents"`
	Currency    string   `json:"currency"`
	Interval    string   `json:"interval"`
	IsPremium   bool     `json:"isPremium"`
	IsActive    bool     `json:"isActive"`
	Available   bool     `json:"available"`
	Features    []string `json:"features"`
}

type SubscriptionResponse struct {
	ID                 string        `json:"id"`
	UserID             string        `json:"userId"`
	PlanID             string        `json:"planId"`
	Plan               *PlanResponse `json:"plan,omitempty"`
	Status             string        `json:"status"`
	CurrentPeriodStart time.Time     `json:"currentPeriodStart"`
	CurrentPeriodEnd   *time.Time    `json:"currentPeriodEnd,omitempty"`
	CancelAtPeriodEnd  bool          `json:"cancelAtPeriodEnd"`
	CanceledAt         *time.Time    `json:"canceledAt,omitempty"`
	CreatedAt          time.Time     `json:"createdAt"`
}

type PaymentResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"userId"`
	SubscriptionID *string    `json:"subscriptionId,omitempty"`
	PlanID         *string    `json:"planId,omitempty"`
	AmountCents    int64      `json:"amountCents"`
	Currency       string     `json:"currency"`
	Status         string     `json:"status"`
	Provider       *string    `json:"provider,omitempty"`
	FailureReason  *string    `json:"failureReason,omitempty"`
	PaidAt         *time.Time `json:"paidAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type CheckoutRequest struct {
	PlanID string `json:"planId" binding:"required"`
}

type CheckoutResponse struct {
	Message      string                `json:"message"`
	Available    bool                  `json:"available"`
	Subscription *SubscriptionResponse `json:"subscription,omitempty"`
	Payment      *PaymentResponse      `json:"payment,omitempty"`
}

type CancelSubscriptionResponse struct {
	Message      string                `json:"message"`
	Subscription *SubscriptionResponse `json:"subscription"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
