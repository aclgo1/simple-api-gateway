package pix

import (
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/payment"
)

type paymentServicePix struct {
	paymentInterface  payment.PaymentInterface
}

func NewpaymentServicePix(paymentInterface payment.PaymentInterface) *paymentServicePix {
	return &paymentServicePix{
		paymentInterface:  paymentInterface,
	}
}

func (s *paymentServicePix) WebHookPix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
