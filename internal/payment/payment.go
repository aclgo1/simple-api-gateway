package payment

import (
	"context"
	"errors"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
)

type PaymentInterface interface {
	RegisterProvider(string, models.PaymentProcessor)
	GeneratePayment(context.Context, *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error)
}


var (
	ErrPaymentMethodNotSupported = errors.New("payment method not supported")
	ErrExceddedLimitGenPix       = errors.New("excedded limit generate pix")
)

type ParamCreditInput struct {
	Type        string  `json:"type"`
	AccountId   string  `json:"account_id"`
	Amount      int64 `json:"amount"`
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
	Balance   int64   `json:"balance"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

