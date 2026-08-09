package models

import (
	"context"
	"time"
)

type Plan string
type StatusPayment string

var (
	PaymentMethodPix             = "pix"
	PaymentMethodCard            = "credit-card"
	PaymentMethodBoleto = "boleto"
	PaymentMethodInternalBalance = "internal-balance"
	Plan_Week Plan = "7_days"
	Plan_Month Plan = "1_month"
	Plan_Year Plan = "1_year"
	Plan_Undefined Plan = "undefined"
	PaymentPaid StatusPayment = "PAID"
	PaymentFailed StatusPayment = "FAILED"
	PaymentRefunded StatusPayment = "REFUNDED"
	PaymentCancelled StatusPayment = "CANCELLED"
	PaymentPending StatusPayment =  "PENDING"
	PaymentUnspecified StatusPayment = "UNSPECIFIED"

)

type ParamPaymentProcessInput struct {
	Method    string  `json:"method"`
	AccountId string  `json:"account_id"`
	ReferenceId string `json:"reference_id"`
	Amount    int64 `json:"amount"`
	CardToken      string `json:"card_token,omitempty"`
	CardExpiration string `json:"card_expiration,omitempty"`
}

type ParamPaymentProcessOutput struct {
	Method               string     `json:"method"`
	Status StatusPayment `json:"status"`
	GatewayTransactionID string     `json:"gateway_transaction_id"`

	CardToken string `json:"card_token"`
    CardExpiration string `json:"card_expiration"`
	
	PixQRCode            string     `json:"pix_qr_code,omitempty"`
	PixExpiration        time.Time `json:"pix_expiration"`
	
	BoletoURL            string     `json:"boleto_url,omitempty"`
	BoletoBarcode        string     `json:"boleto_barcode,omitempty"`
	BoletoExpiration     time.Time `json:"boleto_expiration"`
}

type PaymentProcessor interface {
	Proccess(ctx context.Context, in *ParamPaymentProcessInput) (*ParamPaymentProcessOutput, error)
	Webhook(ctx context.Context, in *ParamPixWebHookInput)(error)
}


type ParamPixWebHookInput struct {
	GatewayTransactionId string `json:"gateway_transaction_id"`
}
