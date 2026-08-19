package usecase

import (
	"context"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/google/uuid"
)

type paymentProcessorCard struct {
}

func NewpaymentProcessorCard() models.PaymentProcessor {
	return &paymentProcessorCard{}
}

func (p *paymentProcessorCard) Proccess(ctx context.Context, in *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {
	return &models.ParamPaymentProcessOutput{Method: in.Method, Status: models.PaymentPaid, GatewayTransactionID: uuid.NewString()}, nil
}
