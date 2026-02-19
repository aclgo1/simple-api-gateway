package usecase

import (
	"context"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
	"github.com/aclgo/simple-api-gateway/internal/wallet/card"
)

type paymentProcessorCard struct {
}

func NewpaymentProcessorCard() card.PaymentProcessor {
	return &paymentProcessorCard{}
}

func (p *paymentProcessorCard) Proccess(ctx context.Context, in *wallet.ParamPaymentProcessorInput) (any, error) {
	return nil, nil
}
