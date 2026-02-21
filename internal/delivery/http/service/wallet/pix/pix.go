package pix

import (
	"context"
	"net/http"

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

	}
}

func (s *walletServicePix) WebHookPix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
