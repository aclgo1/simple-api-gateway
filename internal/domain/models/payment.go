package models

import (
	"context"
	"time"
)

type ParamPaymentProcessInput struct {
	Method    string  `json:"method"`
	AccountId string  `json:"account_id"`
	Amount    int64 `json:"amount"`
	
	CardToken      string `json:"card_token,omitempty"`
	CardExpiration string `json:"card_expiration,omitempty"`
}

type ParamPaymentProcessOutput struct {
	Method               string     `json:"method"`
	GatewayTransactionID string     `json:"gateway_transaction_id"`

	CardToken string `json:"card_token"`
    CardExpiration string `json:"card_expiration"`
	
	PixQRCode            string     `json:"pix_qr_code,omitempty"`
	PixExpiration        time.Time `json:"pix_expiration,omitempty"`
	
	BoletoURL            string     `json:"boleto_url,omitempty"`
	BoletoBarcode        string     `json:"boleto_barcode,omitempty"`
	BoletoExpiration     time.Time `json:"boleto_expiration,omitempty"`
}

type PaymentProcessor interface {
	Proccess(ctx context.Context, in *ParamPaymentProcessInput) (*ParamPaymentProcessOutput, error)
}

var (
	PaymentMethodPix             = "pix"
	PaymentMethodCard            = "credit-card"
	PaymentMethodBoleto = "boleto"
)