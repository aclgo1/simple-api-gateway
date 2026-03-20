package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
	"github.com/aclgo/simple-api-gateway/internal/wallet/pix"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	proto "github.com/aclgo/simple-api-gateway/proto-service/balance"
	"github.com/redis/go-redis/v9"
)

type walletUC struct {
	clientBalanceGPRC proto.WalletServiceClient
	providers         map[string]wallet.PaymentProcessor
	mu                sync.RWMutex
	logger            logger.Logger
	repository        pix.Repository
}

func NewwalletUC(clientBalanceGRPC proto.WalletServiceClient, r pix.Repository, logger logger.Logger) wallet.WalletInterface {
	return &walletUC{
		clientBalanceGPRC: clientBalanceGRPC,
		providers:         make(map[string]wallet.PaymentProcessor),
		logger:            logger,
		repository:        r,
	}
}

func (w *walletUC) RegisterProvider(method string, proccessor wallet.PaymentProcessor) {
	w.mu.Lock()
	w.providers[method] = proccessor
	w.mu.Unlock()
}

func (u *walletUC) Credit(ctx context.Context, in *wallet.ParamCreditInput) (*wallet.ParamCreditOutput, error) {

	in.Amount = math.Abs(in.Amount)

	pct := proto.ParamCreateTransactionRequest{
		ReferenceId: in.ReferenceId,
	}

	if _, err := u.clientBalanceGPRC.CreateTransaction(ctx, &pct); err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.CreateTransaction: %w", err)
	}

	ig := proto.ParamGetWalletByAccountRequest{
		AccountID: in.AccountId,
	}

	wlt, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &ig)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.GetWalletByAccount: %w", err)
	}

	ic := proto.ParamCreditWalletRequest{
		WalletID: wlt.WalletID,
		Amount:   in.Amount,
	}

	resp, err := u.clientBalanceGPRC.Credit(ctx, &ic)
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
	err := u.repository.Get(ctx, in.AccountId)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == nil {
		return nil, wallet.ErrExceddedLimitGenPix
	}

	u.mu.RLock()
	provider, ok := u.providers[in.Method]
	if !ok {
		return nil, wallet.ErrPaymentMethodNotSupported
	}

	u.mu.RUnlock()

	switch in.Amount {
	case 60.00:
	case 40.00:
	case 20.00:
	case 10.00:
	default:
		return nil, errors.New("amount not supported")
	}

	ppi := wallet.ParamPaymentProcessorInput{
		AccountId: in.AccountId,
		Amount:    in.Amount,
	}

	data, err := provider.Proccess(ctx, &ppi)
	if err != nil {
		return nil, err
	}

	if err := u.repository.Set(ctx, in.AccountId); err != nil {
		u.logger.Errorf("u.repository.Set: %v", err)
	}

	out := wallet.ParamGeneratePaymentOutput{
		Type: in.Method,
		Data: data,
	}

	return &out, nil
}
