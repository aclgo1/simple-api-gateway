package card

import "github.com/aclgo/simple-api-gateway/internal/wallet"

type walletServiceCard struct {
	paymentProcessor wallet.PaymentProcessor
	walletInterface  wallet.WalletInterface
}

func NewwalletServicePix(paymentProcessor wallet.PaymentProcessor, walletInterface wallet.WalletInterface) *walletServiceCard {
	return &walletServiceCard{
		paymentProcessor: paymentProcessor,
		walletInterface:  walletInterface,
	}
}
