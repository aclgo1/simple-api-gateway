package pix

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/internal/wallet"
)

type walletServicePix struct {
	paymentProcessor wallet.PaymentProcessor
	walletInterface  wallet.WalletInterface
}

func NewwalletServicePix(paymentProcessor wallet.PaymentProcessor, walletInterface wallet.WalletInterface) *walletServicePix {
	return &walletServicePix{
		paymentProcessor: paymentProcessor,
		walletInterface:  walletInterface,
	}
}

func (s *walletServicePix) CreatePix(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input wallet.ParamGeneratePaymentInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			resp := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, resp, http.StatusBadRequest)
			return
		}

		if err := input.Validate(); err != nil {
			resp := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, resp, http.StatusBadRequest)
			return
		}

		out, err := s.walletInterface.GeneratePayment(ctx, &input)
		if err != nil {
			resp := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, resp, http.StatusInternalServerError)
			return
		}

		service.JSON(w, out, http.StatusCreated)
	}
}

func (s *walletServicePix) WebHookPix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
