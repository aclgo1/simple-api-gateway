package wallet

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletInterface interface {
	RegisterProvider(string, PaymentProcessor)
	Credit(context.Context, *ParamCreditInput) (*ParamCreditOutput, error)
	GeneratePayment(context.Context, *ParamGeneratePaymentInput) (*ParamGeneratePaymentOutput, error)
}

type PaymentProcessor interface {
	Proccess(context.Context, float64) (any, error)
}

var (
	PaymentMethodPix             = "pix"
	PaymentMethodCard            = "card"
	ErrPaymentMethodNotSupported = errors.New("payment method not supported")
)

type ParamCreditInput struct {
	WalletId string
	Amount   float64
}

func (i *ParamCreditInput) Validate() error {
	if i.WalletId == "" {
		return errors.New("wallet id empty")
	}

	_, err := primitive.ObjectIDFromHex(i.WalletId)
	if err != nil {
		return errors.New("invalid wallet id")
	}

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

type ParamGeneratePaymentInput struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type ParamGeneratePaymentOutput struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
