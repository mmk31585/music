package config

import "time"

type PaymentConfig struct {
	Gateway          string
	ZarinpalMerchant string
	ZarinpalSandbox  bool
	IDPayAPIKey      string
	IDPaySandbox     bool
	HTTPTimeout      time.Duration
}

func loadPaymentConfig() PaymentConfig {
	return PaymentConfig{
		Gateway:          getEnv("PAYMENT_GATEWAY", "zarinpal"),
		ZarinpalMerchant: getEnv("ZARINPAL_MERCHANT_ID", ""),
		ZarinpalSandbox:  getEnvAsBool("ZARINPAL_SANDBOX", true),
		IDPayAPIKey:      getEnv("IDPAY_API_KEY", ""),
		IDPaySandbox:     getEnvAsBool("IDPAY_SANDBOX", true),
		HTTPTimeout:      getEnvAsDurationSeconds("PAYMENT_HTTP_TIMEOUT", 30),
	}
}
