package tips

import (
	"context"
	"time"

	"github.com/google/uuid"
	"music/internal/platform/payment"
)

type Service struct {
	repo    *Repository
	gateway payment.Gateway
}

func NewService(repo *Repository, gw payment.Gateway) *Service {
	return &Service{repo: repo, gateway: gw}
}

func (s *Service) CreateTip(ctx context.Context, senderID, artistID string, trackID *string, amountCents int64, currency, message, callbackURL string) (*Tip, string, error) {
	tip := Tip{
		ID:          uuid.New().String(),
		SenderID:    senderID,
		ArtistID:    artistID,
		TrackID:     trackID,
		AmountCents: amountCents,
		Currency:    currency,
		Message:     &message,
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, tip); err != nil {
		return nil, "", err
	}

	payReq := payment.PayRequest{
		AmountCents: amountCents,
		Currency:    currency,
		Description: "Tip for artist " + artistID,
		CallbackURL: callbackURL,
		Metadata: map[string]string{
			"order_id":  tip.ID,
			"purpose":   "tip",
			"sender_id": senderID,
			"artist_id": artistID,
		},
	}

	resp, err := s.gateway.RequestPayment(ctx, payReq)
	if err != nil {
		return nil, "", err
	}

	return &tip, resp.RedirectURL, nil
}

func (s *Service) VerifyTip(ctx context.Context, tipID, authority, status string) (*Tip, error) {
	tip, err := s.repo.GetByID(ctx, tipID)
	if err != nil {
		return nil, err
	}

	if status != "OK" {
		_ = s.repo.MarkPaid(ctx, tipID, s.gateway.Name(), "")
		tip.Status = "failed"
		return tip, nil
	}

	verifyResp, err := s.gateway.VerifyPayment(ctx, payment.VerifyReq{
		Authority:   authority,
		AmountCents: tip.AmountCents,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	if !verifyResp.Success {
		return tip, nil
	}

	if err := s.repo.MarkPaid(ctx, tipID, s.gateway.Name(), verifyResp.ProviderPayID); err != nil {
		return nil, err
	}

	tip.Status = "completed"
	tip.Provider = &verifyResp.ProviderPayID
	now := time.Now().UTC()
	tip.PaidAt = &now
	return tip, nil
}

func (s *Service) ListBySender(ctx context.Context, userID string, limit, offset int) ([]Tip, error) {
	return s.repo.ListBySender(ctx, userID, limit, offset)
}

func (s *Service) ListByArtist(ctx context.Context, artistID string, limit, offset int) ([]Tip, error) {
	return s.repo.ListByArtist(ctx, artistID, limit, offset)
}

func (s *Service) TotalForArtist(ctx context.Context, artistID string) (int64, error) {
	return s.repo.TotalForArtist(ctx, artistID)
}
