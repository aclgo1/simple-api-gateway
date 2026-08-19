package usecase

import (
	"context"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	proto "github.com/aclgo/simple-api-gateway/proto-service/balance"
	"github.com/google/uuid"
)

type paymentProcessorWallet struct {
	walletGRPC proto.WalletServiceClient
}

func NewPaymentProcessorWallet(walletClient proto.WalletServiceClient) models.PaymentProcessor {
	return &paymentProcessorWallet{
		walletGRPC: walletClient,
	}
}

func (p *paymentProcessorWallet) Proccess(ctx context.Context, params *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {

	out := models.ParamPaymentProcessOutput{
		Method:               params.Method,
		GatewayTransactionID: uuid.NewString(),
		Status:               models.PaymentFailed,
	}

	paramGet := proto.ParamGetWalletByAccountRequest{
		AccountID: params.AccountId,
	}

	wlt, err := p.walletGRPC.GetWalletByAccount(ctx, &paramGet)
	if err != nil {
		return &out, fmt.Errorf("p.walletGRPC.GetWalletByAccount: %w", err)
	}

	ref := params.ReferenceId
	if ref == "" {
		ref = out.GatewayTransactionID
	}

	paramDebit := proto.ParamDebitWalletRequest{
		WalletID:    wlt.WalletID,
		Amount:      params.Amount,
		ReferenceID: ref,
	}

	_, err = p.walletGRPC.Debit(ctx, &paramDebit)
	if err != nil {
		return &out, err
	}

	out.Status = models.PaymentPaid

	return &out, nil
}
