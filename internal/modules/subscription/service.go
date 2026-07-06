package subscription

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"

	"music/internal/platform/events"
)

var (
	ErrPremiumNotAvailable = errors.New("premium is not available yet")
	ErrInactivePlan        = errors.New("plan is inactive")
)

type Service struct {
	repo      *Repository
	publisher events.Publisher
}

func NewService(repo *Repository, publisher events.Publisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) ListPlans(ctx context.Context) ([]PlanResponse, error) {
	plans, err := s.repo.ListPlans(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]PlanResponse, 0, len(plans))
	for _, plan := range plans {
		result = append(result, mapPlan(plan))
	}

	return result, nil
}

func (s *Service) CurrentSubscription(ctx context.Context, userID string) (*SubscriptionResponse, error) {
	sub, err := s.repo.GetCurrentSubscription(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		freePlan, freeErr := s.repo.GetPlanByCode(ctx, "free")
		if freeErr != nil {
			return nil, freeErr
		}

		newSub, createErr := s.repo.CreateSubscription(ctx, userID, freePlan.ID, "active")
		if createErr != nil {
			return nil, createErr
		}

		mapped := mapSubscription(*newSub)
		plan := mapPlan(*freePlan)
		mapped.Plan = &plan

		return &mapped, nil
	}

	if err != nil {
		return nil, err
	}

	plan, err := s.repo.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	mapped := mapSubscription(*sub)
	mappedPlan := mapPlan(*plan)
	mapped.Plan = &mappedPlan

	return &mapped, nil
}

func (s *Service) ListSubscriptions(ctx context.Context, userID string) ([]SubscriptionResponse, error) {
	subs, err := s.repo.ListSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		item := mapSubscription(sub)

		plan, planErr := s.repo.GetPlanByID(ctx, sub.PlanID)
		if planErr == nil {
			mappedPlan := mapPlan(*plan)
			item.Plan = &mappedPlan
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *Service) Checkout(ctx context.Context, userID string, req CheckoutRequest) (*CheckoutResponse, error) {
	plan, err := s.repo.GetPlanByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	if !plan.IsActive {
		return nil, ErrInactivePlan
	}

	// Premium is future-only now.
	// This protects you from accidentally selling premium before payment gateway exists.
	if plan.IsPremium {
		return nil, ErrPremiumNotAvailable
	}

	// Idempotency: check if the user already has an active subscription to this plan
	existing, _ := s.repo.GetCurrentSubscription(ctx, userID)
	if existing != nil && existing.PlanID == plan.ID && existing.Status == "active" {
		mappedSub := mapSubscription(*existing)
		mappedPlan := mapPlan(*plan)
		mappedSub.Plan = &mappedPlan
		return &CheckoutResponse{
			Message:      "Already subscribed to this plan.",
			Available:    true,
			Subscription: &mappedSub,
		}, nil
	}

	sub, err := s.repo.CreateSubscription(ctx, userID, plan.ID, "active")
	if err != nil {
		return nil, err
	}

	payment, err := s.repo.CreatePayment(
		ctx,
		userID,
		sub.ID,
		plan.ID,
		plan.PriceCents,
		plan.Currency,
		"paid",
	)
	if err != nil {
		return nil, err
	}

	s.publishSubscriptionPurchased(ctx, userID, sub.ID, plan.Name, int(plan.PriceCents), plan.Currency)

	mappedSub := mapSubscription(*sub)
	mappedPlan := mapPlan(*plan)
	mappedSub.Plan = &mappedPlan

	mappedPayment := mapPayment(*payment)

	return &CheckoutResponse{
		Message:      "Subscription created.",
		Available:    true,
		Subscription: &mappedSub,
		Payment:      &mappedPayment,
	}, nil
}

func (s *Service) CancelCurrentSubscription(ctx context.Context, userID string) (*CancelSubscriptionResponse, error) {
	sub, err := s.repo.CancelCurrentSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	mapped := mapSubscription(*sub)

	return &CancelSubscriptionResponse{
		Message:      "Subscription canceled.",
		Subscription: &mapped,
	}, nil
}

func (s *Service) ListPayments(ctx context.Context, userID string) ([]PaymentResponse, error) {
	payments, err := s.repo.ListPayments(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		result = append(result, mapPayment(payment))
	}

	return result, nil
}

func (s *Service) publishSubscriptionPurchased(
	ctx context.Context,
	userID string,
	subscriptionID string,
	planName string,
	amountCents int,
	currency string,
) {
	if s.publisher == nil {
		return
	}

	uid, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return
	}

	sid, err := uuid.Parse(strings.TrimSpace(subscriptionID))
	if err != nil {
		return
	}

	_ = s.publisher.Publish(ctx, events.SubscriptionPurchasedEvent{
		BaseEvent:      events.NewBaseEvent(),
		UserID:         uid,
		SubscriptionID: sid,
		Plan:           planName,
		Amount:         int64(amountCents),
		Currency:       currency,
	})
}

func mapPlan(plan Plan) PlanResponse {
	var features []string

	if len(plan.Features) > 0 {
		_ = json.Unmarshal(plan.Features, &features)
	}

	return PlanResponse{
		ID:          plan.ID,
		Code:        plan.Code,
		Name:        plan.Name,
		Description: plan.Description,
		PriceCents:  plan.PriceCents,
		Currency:    plan.Currency,
		Interval:    plan.Interval,
		IsPremium:   plan.IsPremium,
		IsActive:    plan.IsActive,
		Available:   plan.IsActive && !plan.IsPremium,
		Features:    features,
	}
}

func mapSubscription(sub Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                 sub.ID,
		UserID:             sub.UserID,
		PlanID:             sub.PlanID,
		Status:             sub.Status,
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
		CancelAtPeriodEnd:  sub.CancelAtPeriodEnd,
		CanceledAt:         sub.CanceledAt,
		CreatedAt:          sub.CreatedAt,
	}
}

func mapPayment(payment Payment) PaymentResponse {
	return PaymentResponse{
		ID:             payment.ID,
		UserID:         payment.UserID,
		SubscriptionID: payment.SubscriptionID,
		PlanID:         payment.PlanID,
		AmountCents:    payment.AmountCents,
		Currency:       payment.Currency,
		Status:         payment.Status,
		Provider:       payment.Provider,
		FailureReason:  payment.FailureReason,
		PaidAt:         payment.PaidAt,
		CreatedAt:      payment.CreatedAt,
	}
}
