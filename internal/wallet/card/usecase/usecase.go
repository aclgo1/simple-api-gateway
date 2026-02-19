package usecase

import "context"

type paymentProcessorCard struct {
}

func NewpaymentProcessorCard() *paymentProcessorCard {
	return &paymentProcessorCard{}
}

func (p *paymentProcessorCard) Proccess(ctx context.Context, amount float64) (any, error) {
	return nil, nil
}
