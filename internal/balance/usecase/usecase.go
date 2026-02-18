package usecase

// import (
// 	"context"
// 	"fmt"

// 	"github.com/aclgo/simple-api-gateway/internal/balance"
// 	"github.com/aclgo/simple-api-gateway/pkg/logger"
// 	proto "github.com/aclgo/simple-api-gateway/proto-service/balance"
// )

// type balanceUC struct {
// 	clientBalanceGPRC proto.WalletServiceClient
// 	logger            logger.Logger
// }

// func NewBalanceUC(clientBalanceGRPC proto.WalletServiceClient, logger logger.Logger) balance.Balance {
// 	return &balanceUC{
// 		clientBalanceGPRC: clientBalanceGRPC,
// 		logger:            logger,
// 	}
// }

// func (u *balanceUC) Credit(ctx context.Context, in *balance.ParamCreditInput) (*balance.ParamCreditOutput, error) {

// 	ip := proto.ParamCreditWalletRequest{
// 		WalletID: in.WalletId,
// 		Amount:   in.Amount,
// 	}

// 	resp, err := u.clientBalanceGPRC.Credit(ctx, &ip)
// 	if err != nil {
// 		return nil, fmt.Errorf("u.clientBalanceGRPC.Credit: %w", err)
// 	}

// 	out := balance.ParamCreditOutput{
// 		WalletID:  resp.WalletID,
// 		AccountID: resp.AccountID,
// 		Balance:   resp.Balance,
// 		CreatedAT: resp.CreatedAT.AsTime(),
// 		UpdatedAT: resp.UpdatedAT.AsTime(),
// 	}

// 	return &out, nil
// }
