package usecase

import (
	"context"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/orders"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoOrders "github.com/aclgo/simple-api-gateway/proto-service/orders"
	protoProduct "github.com/aclgo/simple-api-gateway/proto-service/product"
)

type OrderCreateSagaUC struct {
	clientProductsGRPC protoProduct.ProductServiceClient
	clientOrdersGRPC   protoOrders.ServiceOrderClient
	clientBalanceGPRC  protoBalance.WalletServiceClient
	sagaWorkerCompensate orders.SagaWorker
}

func NewOrderCreateSaga(clientProductsGRPC protoProduct.ProductServiceClient,
	clientOrdersGRPC   protoOrders.ServiceOrderClient,
	clientBalanceGPRC  protoBalance.WalletServiceClient,
	sagaWorkerCompensate orders.SagaWorker,
	) *OrderCreateSagaUC{
	return &OrderCreateSagaUC{
		clientProductsGRPC: clientProductsGRPC,
		clientOrdersGRPC: clientOrdersGRPC,
		clientBalanceGPRC: clientBalanceGPRC,
		sagaWorkerCompensate: sagaWorkerCompensate,
	}
}

func(u *OrderCreateSagaUC)Execute(ctx context.Context, in *orders.OrderCreateInput)(*orders.OrderCreateOutput,error){
	var amountProducts float64
	for _, pID := range in.ProductsIDS {
		product, err := u.clientProductsGRPC.Find(ctx, &protoProduct.ProductFindRequest{Id: pID})
		if err != nil {
			return nil, fmt.Errorf("failed to find product %s: %w", pID, err)
		}
		amountProducts += product.Product.Price
	}

	wallet, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &protoBalance.ParamGetWalletByAccountRequest{AccountID: in.AccountId})
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.Balance < amountProducts {
		return nil, fmt.Errorf("insufficient funds: amount is %.2f, balance is %.2f", amountProducts, wallet.Balance)
	}

	var compensations []func(context.Context)error

	rollback := func(originalErr error)error{
		u.sagaWorkerCompensate.AppendTask(&orders.CompensationTask{
			OriginalErr: originalErr,
			Compensations: compensations,
		})
		return originalErr
	}

	_, err = u.clientBalanceGPRC.Debit(ctx, &protoBalance.ParamDebitWalletRequest{
		WalletID: wallet.WalletID,
		Amount:   amountProducts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}

	compensations = append(compensations, func(ctx context.Context) error {
		_, err := u.clientBalanceGPRC.Credit(ctx, &protoBalance.ParamCreditWalletRequest{
			WalletID: wallet.WalletID,
			Amount:   amountProducts,
		})
		return err
	})

	successfullyUpdatedProducts:= make([]string, 0,len(in.ProductsIDS))

	compensations = append(compensations, func(ctx context.Context) error {
		for _, pId := range successfullyUpdatedProducts{
			_, err := u.clientProductsGRPC.Update(ctx, &protoProduct.ProductUpdateRequest{
				Id:         pId,
				HasOrdered: false, 
			})
			if err != nil {
				return err 
			} 
		}

		return nil
	})

	for _, pID := range in.ProductsIDS {
		_, err := u.clientProductsGRPC.Update(ctx, &protoProduct.ProductUpdateRequest{
			Id:         pID,
			HasOrdered: true,
		})
		if err != nil {
			return nil, rollback(fmt.Errorf("failed to update product %s: %w", pID, err))
		}
		successfullyUpdatedProducts = append(successfullyUpdatedProducts, pID)
	}

	orderCreate, err := u.clientOrdersGRPC.Create(ctx, &protoOrders.ParamCreateOrderRequest{
		AccountID:   in.AccountId,
		ProductsIDS: in.ProductsIDS,
	})
	if err != nil {
		return nil, rollback(fmt.Errorf("failed to create order: %w", err))
	}

	out := &orders.OrderCreateOutput{
		OrderId:     orderCreate.Order.OrderID,
		AccountId:   orderCreate.Order.AccountID,
		ProductsIDS: orderCreate.Order.ProductsIDS,
		CreatedAt:   orderCreate.Order.CreatedAT.AsTime(),
	}

	return out,nil
}

