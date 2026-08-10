package pix

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/internal/domain/models"
)

type paymentServicePix struct {
	pixUseCase  models.PixPaymentWebHook
}

func NewpaymentServicePix(pix models.PixPaymentWebHook) *paymentServicePix{
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
			resp := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
			service.JSON(w,resp, http.StatusBadRequest)
			return
		}

		if err := params.Validate(); err != nil {
			resp := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
			service.JSON(w,resp, http.StatusBadRequest)
			return
		}

		err := s.pixUseCase.Webhook(r.Context(), &params)
		if err != nil {
			resp := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w,resp, http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
