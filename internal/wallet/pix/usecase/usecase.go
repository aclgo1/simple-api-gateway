package usecase

import "context"

type paymentProcessorPix struct {
}

func NewpaymentProcessorPix() *paymentProcessorPix {
	return &paymentProcessorPix{}
}

func (p *paymentProcessorPix) Proccess(ctx context.Context, amount float64) (any, error) {
	return nil, nil
}
