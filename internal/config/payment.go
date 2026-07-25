package config

import (
	"os"
	"strconv"
	"time"
)

type PaymentConfig struct {
	Gateway          string
	ZarinpalMerchant string
	ZarinpalSandbox  bool
	IDPayAPIKey      string
	IDPaySandbox     bool
	HTTPTimeout      time.Duration
}

func loadPaymentConfig() PaymentConfig {
	gateway := os.Getenv("PAYMENT_GATEWAY")
	if gateway == "" {
		gateway = "zarinpal"
	}

	zarinpalSandbox := true
	if v := os.Getenv("ZARINPAL_SANDBOX"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			zarinpalSandbox = false
		}
	}
	idpaySandbox := true
	if v := os.Getenv("IDPAY_SANDBOX"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			idpaySandbox = false
		}
	}

	httpTimeout := 30
	if v := os.Getenv("PAYMENT_HTTP_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			httpTimeout = parsed
		}
	}

	return PaymentConfig{
		Gateway:          gateway,
		ZarinpalMerchant: os.Getenv("ZARINPAL_MERCHANT_ID"),
		ZarinpalSandbox:  zarinpalSandbox,
		IDPayAPIKey:      os.Getenv("IDPAY_API_KEY"),
		IDPaySandbox:     idpaySandbox,
		HTTPTimeout:      time.Duration(httpTimeout) * time.Second,
	}
}
