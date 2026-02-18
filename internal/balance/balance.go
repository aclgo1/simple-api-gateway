package balance

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Balance interface {
	Credit(context.Context, *ParamCreditInput) (*ParamCreditOutput, error)
}

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
