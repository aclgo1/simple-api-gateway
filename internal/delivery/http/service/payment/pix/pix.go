package pix

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
)

type paymentServicePix struct {
	pixUseCase  models.PaymentProcessor
}

func NewpaymentServicePix(pix models.PaymentProcessor) *paymentServicePix{
	if pix == nil {
		log.Fatal("pix usecase is nil")
	}
	return &paymentServicePix{
		pixUseCase:  pix,
	}
}

func (s *paymentServicePix) WebHookPix(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.ParamPixWebHookInput

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		}

		if err := params.Validate(); err != nil {
			
		}

		err := s.pixUseCase.Webhook(r.Context(), &params)
		if err != nil {
			
		}
	}
}
