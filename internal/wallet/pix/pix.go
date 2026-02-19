package pix

import (
	"context"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
)

type PaymentProcessor interface {
	Proccess(context.Context, *wallet.ParamPaymentProcessorInput) (any, error)
}
