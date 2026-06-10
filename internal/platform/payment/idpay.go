package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	idpayAPI           = "https://api.idpay.ir/v1.1/payment"
	idpayVerifyAPI     = "https://api.idpay.ir/v1.1/payment/verify"
	idpaySandboxAPI    = "https://sandbox.idpay.ir/v1.1/payment"
	idpaySandboxVerify = "https://sandbox.idpay.ir/v1.1/payment/verify"
)

type IDPayConfig struct {
	APIKey     string
	Sandbox    bool
	HTTPClient *http.Client
}

type idpayRequest struct {
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	Name        string `json:"name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Description string `json:"desc,omitempty"`
	Callback    string `json:"callback"`
}

type idpayResponse struct {
	ID        string `json:"id"`
	Link      string `json:"link"`
	Status    int    `json:"status"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_message,omitempty"`
}

type idpayVerifyReq struct {
	ID      string `json:"id"`
	OrderID string `json:"order_id"`
}

type idpayVerifyResp struct {
	Status   int    `json:"status"`
	TrackID  int    `json:"track_id"`
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"`
	CardNo   string `json:"card_no"`
	HashedNo string `json:"hashed_card_no"`
	Date     int64  `json:"date"`
}

type IDPayGateway struct {
	config IDPayConfig
}

func NewIDPayGateway(cfg IDPayConfig) *IDPayGateway {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &IDPayGateway{config: cfg}
}

func (g *IDPayGateway) Name() string {
	return "idpay"
}

func (g *IDPayGateway) requestURL() string {
	if g.config.Sandbox {
		return idpaySandboxAPI
	}
	return idpayAPI
}

func (g *IDPayGateway) verifyURL() string {
	if g.config.Sandbox {
		return idpaySandboxVerify
	}
	return idpayVerifyAPI
}

func (g *IDPayGateway) doPost(ctx context.Context, url string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", g.config.APIKey)
	req.Header.Set("X-Sandbox", fmt.Sprintf("%t", g.config.Sandbox))
	return g.config.HTTPClient.Do(req)
}

func (g *IDPayGateway) RequestPayment(ctx context.Context, req PayRequest) (*PayResponse, error) {
	idReq := idpayRequest{
		OrderID:     req.Metadata["order_id"],
		Amount:      req.AmountCents,
		Description: req.Description,
		Callback:    req.CallbackURL,
	}

	resp, err := g.doPost(ctx, g.requestURL(), idReq)
	if err != nil {
		return nil, fmt.Errorf("idpay request: %w", err)
	}
	defer resp.Body.Close()

	var result idpayResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("idpay decode: %w", err)
	}

	if result.ErrorCode != "" {
		return nil, fmt.Errorf("%w: %s (code=%s)", ErrPaymentFailed, result.ErrorMsg, result.ErrorCode)
	}

	return &PayResponse{
		Authority:   result.ID,
		RedirectURL: result.Link,
		Status:      "pending",
	}, nil
}

func (g *IDPayGateway) VerifyPayment(ctx context.Context, req VerifyReq) (*VerifyResp, error) {
	orderID := ""
	if req.Metadata != nil {
		orderID = req.Metadata["order_id"]
	}
	verifyReq := idpayVerifyReq{
		ID:      req.Authority,
		OrderID: orderID,
	}

	resp, err := g.doPost(ctx, g.verifyURL(), verifyReq)
	if err != nil {
		return nil, fmt.Errorf("idpay verify: %w", err)
	}
	defer resp.Body.Close()

	var result idpayVerifyResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("idpay verify decode: %w", err)
	}

	// IDPay status codes: 1=paid, 2=failed, 3=canceled, ...
	if result.Status != 1 {
		return &VerifyResp{
			Success: false,
			Message: fmt.Sprintf("payment status: %d", result.Status),
		}, nil
	}

	return &VerifyResp{
		Success:       true,
		ProviderPayID: fmt.Sprintf("%d", result.TrackID),
		Message:       "verified",
	}, nil
}
