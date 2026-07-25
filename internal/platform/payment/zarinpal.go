package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	zarinpalAPI           = "https://api.zarinpal.com/pg/v4/payment/request.json"
	zarinpalVerifyAPI     = "https://api.zarinpal.com/pg/v4/payment/verify.json"
	zarinpalStartPay      = "https://www.zarinpal.com/pg/StartPay/"
	zarinpalSandboxAPI    = "https://sandbox.zarinpal.com/pg/v4/payment/request.json"
	zarinpalSandboxVerify = "https://sandbox.zarinpal.com/pg/v4/payment/verify.json"
	zarinpalSandboxPay    = "https://sandbox.zarinpal.com/pg/StartPay/"
)

type ZarinpalConfig struct {
	MerchantID string
	Sandbox    bool
	HTTPClient *http.Client
}

type zarinpalRequest struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	CallbackURL string `json:"callback_url"`
	Currency    string `json:"currency"`
	Metadata    any    `json:"metadata,omitempty"`
}

type zarinpalResponse struct {
	Data struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Authority string `json:"authority"`
		Fee       int    `json:"fee"`
	} `json:"data"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type zarinpalVerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int64  `json:"amount"`
	Authority  string `json:"authority"`
}

type zarinpalVerifyResponse struct {
	Data struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		CardPAN string `json:"card_pan"`
		RefID   int    `json:"ref_id"`
		Fee     int    `json:"fee"`
	} `json:"data"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type ZarinpalGateway struct {
	config ZarinpalConfig
}

func NewZarinpalGateway(cfg ZarinpalConfig) *ZarinpalGateway {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &ZarinpalGateway{config: cfg}
}

func (g *ZarinpalGateway) Name() string {
	return "zarinpal"
}

func (g *ZarinpalGateway) requestURL() string {
	if g.config.Sandbox {
		return zarinpalSandboxAPI
	}
	return zarinpalAPI
}

func (g *ZarinpalGateway) verifyURL() string {
	if g.config.Sandbox {
		return zarinpalSandboxVerify
	}
	return zarinpalVerifyAPI
}

func (g *ZarinpalGateway) payURL(authority string) string {
	if g.config.Sandbox {
		return zarinpalSandboxPay + authority
	}
	return zarinpalStartPay + authority
}

func (g *ZarinpalGateway) postJSON(ctx context.Context, url string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (g *ZarinpalGateway) RequestPayment(ctx context.Context, req PayRequest) (*PayResponse, error) {
	meta := map[string]string{}
	if req.Metadata != nil {
		meta = req.Metadata
	}

	payload := zarinpalRequest{
		MerchantID:  g.config.MerchantID,
		Amount:      req.AmountCents,
		Description: req.Description,
		CallbackURL: req.CallbackURL,
		Currency:    req.Currency,
		Metadata:    meta,
	}

	data, err := g.postJSON(ctx, g.requestURL(), payload)
	if err != nil {
		return nil, fmt.Errorf("zarinpal request: %w", err)
	}

	var result zarinpalResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("zarinpal decode: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("%w: %s (code=%d)", ErrPaymentFailed, result.Errors[0].Message, result.Errors[0].Code)
	}

	if result.Data.Code != 100 {
		return nil, fmt.Errorf("%w: unexpected code %d", ErrPaymentFailed, result.Data.Code)
	}

	return &PayResponse{
		Authority:   result.Data.Authority,
		RedirectURL: g.payURL(result.Data.Authority),
		Status:      "pending",
	}, nil
}

func (g *ZarinpalGateway) VerifyPayment(ctx context.Context, req VerifyReq) (*VerifyResp, error) {
	payload := zarinpalVerifyRequest{
		MerchantID: g.config.MerchantID,
		Amount:     req.AmountCents,
		Authority:  req.Authority,
	}

	data, err := g.postJSON(ctx, g.verifyURL(), payload)
	if err != nil {
		return nil, fmt.Errorf("zarinpal verify request: %w", err)
	}

	var result zarinpalVerifyResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("zarinpal verify decode: %w", err)
	}

	if len(result.Errors) > 0 {
		if result.Errors[0].Code == 101 {
			return &VerifyResp{
				Success: false,
				Message: "Payment already verified",
			}, nil
		}
		return nil, fmt.Errorf("%w: %s (code=%d)", ErrPaymentFailed, result.Errors[0].Message, result.Errors[0].Code)
	}

	if result.Data.Code != 100 {
		return nil, fmt.Errorf("%w: unexpected verify code %d", ErrPaymentFailed, result.Data.Code)
	}

	return &VerifyResp{
		Success:       true,
		ProviderPayID: fmt.Sprintf("%d", result.Data.RefID),
		Message:       result.Data.Message,
	}, nil
}
