package payment

import (
	"context"
	"errors"
)

var (
	ErrGatewayNotSupported = errors.New("payment gateway not supported")
	ErrPaymentFailed       = errors.New("payment failed")
	ErrPaymentCanceled     = errors.New("payment canceled")
	ErrInvalidCallback     = errors.New("invalid callback data")
)

type PayRequest struct {
	AmountCents int64
	Currency    string
	Description string
	CallbackURL string
	Metadata    map[string]string
}

type PayResponse struct {
	Authority     string
	RedirectURL   string
	ProviderPayID string
	Status        string
}

type VerifyReq struct {
	Authority   string
	Status      string
	AmountCents int64
	Metadata    map[string]string
}

type VerifyResp struct {
	Success       bool
	ProviderPayID string
	Message       string
}

func NewPayRequest(amount int64, currency, desc, callback string, meta map[string]string) PayRequest {
	return PayRequest{
		AmountCents: amount,
		Currency:    currency,
		Description: desc,
		CallbackURL: callback,
		Metadata:    meta,
	}
}

type Gateway interface {
	Name() string
	RequestPayment(ctx context.Context, req PayRequest) (*PayResponse, error)
	VerifyPayment(ctx context.Context, req VerifyReq) (*VerifyResp, error)
}
