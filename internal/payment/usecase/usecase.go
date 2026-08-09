package usecase

import (
	"context"
	"sync"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/aclgo/simple-api-gateway/internal/payment"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	proto "github.com/aclgo/simple-api-gateway/proto-service/balance"
)

type paymentUC struct {
	clientBalanceGPRC proto.WalletServiceClient
	providers         map[string]models.PaymentProcessor
	mu                sync.RWMutex
	logger            logger.Logger
}

func NewPaymentUC(clientBalanceGRPC proto.WalletServiceClient,logger logger.Logger) payment.PaymentInterface {
	return &paymentUC{
		clientBalanceGPRC: clientBalanceGRPC,
		providers:         make(map[string]models.PaymentProcessor),
		logger:            logger,
	}
}

func (w *paymentUC) RegisterProvider(method string, proccessor models.PaymentProcessor) {
	w.mu.Lock()
	w.providers[method] = proccessor
	w.mu.Unlock()
}

// func (u *paymentUC) Credit(ctx context.Context, in *payment.ParamCreditInput) (*payment.ParamCreditOutput, error) {


// 	ig := proto.ParamGetWalletByAccountRequest{
// 		AccountID: in.AccountId,
// 	}

// 	wlt, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &ig)
// 	if err != nil {
// 		return nil, fmt.Errorf("u.clientBalanceGPRC.GetWalletByAccount: %w", err)
// 	}

// 	ic := proto.ParamCreditWalletRequest{
// 		WalletID: wlt.WalletID,
// 		Amount:   in.Amount,
// 		ReferenceID: uuid.NewString(),
// 	}

// 	resp, err := u.clientBalanceGPRC.Credit(ctx, &ic)
// 	if err != nil {
// 		return nil, fmt.Errorf("u.clientBalanceGRPC.Credit: %w", err)
// 	}

// 	out := payment.ParamCreditOutput{
// 		WalletID:  resp.WalletID,
// 		AccountID: resp.AccountID,
// 		Balance:   resp.Balance,
// 		CreatedAT: resp.CreatedAT.AsTime(),
// 		UpdatedAT: resp.UpdatedAT.AsTime(),
// 	}

// 	return &out, nil
// }

func (u *paymentUC) GeneratePayment(ctx context.Context, in *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {
	u.mu.RLock()
	provider, ok := u.providers[in.Method]
	u.mu.RUnlock()

	if !ok {
		return nil, payment.ErrPaymentMethodNotSupported
	}

	return provider.Proccess(ctx,in)
}


