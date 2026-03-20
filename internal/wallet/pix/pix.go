package pix

import (
	"context"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
)

type PaymentProcessor interface {
	Proccess(context.Context, *wallet.ParamPaymentProcessorInput) (any, error)
}

type Repository interface {
	Get(ctx context.Context, key string) error
	Set(ctx context.Context, key string) error
}

func FormatPixKeyRepository(id string) string {
	return fmt.Sprintf("pix_generated:%s", id)
}

type ParamsPixOutput struct {
	Teste string `json:"teste"`
}
