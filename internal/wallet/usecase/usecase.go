package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	proto "github.com/aclgo/simple-api-gateway/proto-service/balance"
)

type walletUC struct {
	clientBalanceGPRC proto.WalletServiceClient
	providers         map[string]wallet.PaymentProcessor
	mu                sync.RWMutex
	logger            logger.Logger
}

func NewwalletUC(clientBalanceGRPC proto.WalletServiceClient, logger logger.Logger) wallet.WalletInterface {
	return &walletUC{
		clientBalanceGPRC: clientBalanceGRPC,
		providers:         make(map[string]wallet.PaymentProcessor),
		logger:            logger,
	}
}

func (w *walletUC) RegisterProvider(method string, proccessor wallet.PaymentProcessor) {
	w.mu.RLock()
	w.providers[method] = proccessor
	w.mu.RUnlock()
}

func (u *walletUC) Credit(ctx context.Context, in *wallet.ParamCreditInput) (*wallet.ParamCreditOutput, error) {

	ip := proto.ParamCreditWalletRequest{
		WalletID: in.WalletId,
		Amount:   in.Amount,
	}

	resp, err := u.clientBalanceGPRC.Credit(ctx, &ip)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGRPC.Credit: %w", err)
	}

	out := wallet.ParamCreditOutput{
		WalletID:  resp.WalletID,
		AccountID: resp.AccountID,
		Balance:   resp.Balance,
		CreatedAT: resp.CreatedAT.AsTime(),
		UpdatedAT: resp.UpdatedAT.AsTime(),
	}

	return &out, nil
}

func (u *walletUC) GeneratePayment(ctx context.Context, in *wallet.ParamGeneratePaymentInput) (*wallet.ParamGeneratePaymentOutput, error) {
	u.mu.RLock()
	provider, ok := u.providers[in.Method]
	if !ok {
		return nil, wallet.ErrPaymentMethodNotSupported
	}

	u.mu.RUnlock()

	data, err := provider.Proccess(ctx, in.Amount)
	if err != nil {
		return nil, err
	}

	out := wallet.ParamGeneratePaymentOutput{
		Type: in.Method,
		Data: data,
	}

	return &out, nil
}
