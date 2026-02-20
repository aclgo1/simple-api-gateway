package wallet

import (
	"context"
	"errors"
	"time"
)

type WalletInterface interface {
	RegisterProvider(string, PaymentProcessor)
	Credit(context.Context, *ParamCreditInput) (*ParamCreditOutput, error)
	GeneratePayment(context.Context, *ParamGeneratePaymentInput) (*ParamGeneratePaymentOutput, error)
}

type PaymentProcessor interface {
	Proccess(context.Context, *ParamPaymentProcessorInput) (any, error)
}

var (
	PaymentMethodPix             = "pix"
	PaymentMethodCard            = "card"
	ErrPaymentMethodNotSupported = errors.New("payment method not supported")
)

type ParamCreditInput struct {
	Type        string  `json:"type"`
	AccountId   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	ReferenceId string  `json:"reference_id"`
}

func (i *ParamCreditInput) Validate() error {
	if i.AccountId == "" {
		return errors.New("account id empty")
	}

	// _, err := primitive.ObjectIDFromHex(i.WalletId)
	// if err != nil {
	// 	return errors.New("invalid wallet id")
	// }

	if i.Amount <= 0 {
		return errors.New("invalid amount")
	}

	return nil
}

type ParamCreditOutput struct {
	WalletID  string    `json:"wallet_id"`
	AccountID string    `json:"account_id"`
	Balance   float64   `json:"balance"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type ParamPaymentProcessorInput struct {
	Method    string  `json:"method"`
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type ParamGeneratePaymentInput struct {
	Method    string  `json:"method"`
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type ParamGeneratePaymentOutput struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
