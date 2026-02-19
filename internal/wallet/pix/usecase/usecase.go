package usecase

import (
	"context"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
	"github.com/aclgo/simple-api-gateway/internal/wallet/pix"
)

type paymentProcessorPix struct {
}

func NewpaymentProcessorPix() pix.PaymentProcessor {
	return &paymentProcessorPix{}
}

func (p *paymentProcessorPix) Proccess(ctx context.Context, in *wallet.ParamPaymentProcessorInput) (any, error) {
	return nil, nil
}
