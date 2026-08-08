package usecase

import (
	"context"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
)

type paymentProcessorCard struct {
}

func NewpaymentProcessorCard() models.PaymentProcessor {
	return &paymentProcessorCard{}
}

func (p *paymentProcessorCard) Proccess(ctx context.Context, in *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {
	return nil, nil
}
