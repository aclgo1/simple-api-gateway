package card

import "github.com/aclgo/simple-api-gateway/internal/payment"

type paymentServiceCard struct {
	paymentInterface  payment.PaymentInterface
}

func NewPaymentServicePix(paymentInterface payment.PaymentInterface) *paymentServiceCard {
	return &paymentServiceCard{
		paymentInterface:  paymentInterface,
	}
}
