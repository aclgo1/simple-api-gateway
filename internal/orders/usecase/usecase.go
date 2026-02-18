package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/aclgo/simple-api-gateway/internal/orders"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoOrders "github.com/aclgo/simple-api-gateway/proto-service/orders"
	protoProduct "github.com/aclgo/simple-api-gateway/proto-service/product"
)

type orderUC struct {
	clientOrdersGRPC   protoOrders.ServiceOrderClient
	clientBalanceGPRC  protoBalance.WalletServiceClient
	clientProductsGRPC protoProduct.ProductServiceClient
	mu                 *sync.Mutex
	logger             logger.Logger
}

func NeworderUC(
	clientOrdersGRPC protoOrders.ServiceOrderClient,
	clientProductsGRPC protoProduct.ProductServiceClient,
	clientBalanceGRPC protoBalance.WalletServiceClient,
	mu *sync.Mutex,
	logger logger.Logger,
) *orderUC {
	return &orderUC{
		clientOrdersGRPC:   clientOrdersGRPC,
		clientBalanceGPRC:  clientBalanceGRPC,
		clientProductsGRPC: clientProductsGRPC,
		mu:                 mu,
		logger:             logger,
	}
}

func (u *orderUC) Create(ctx context.Context, in *orders.OrderCreateInput) (*orders.OrderCreateOutput, error) {

	var amountProducts float64

	for i := range in.ProductsIDS {
		paramProductProto := protoProduct.ProductFindRequest{
			Id: in.ProductsIDS[i],
		}

		product, err := u.clientProductsGRPC.Find(ctx, &paramProductProto)
		if err != nil {
			return nil, fmt.Errorf("u.clientProductsGRPC.Find: %w", err)
		}
		u.mu.Lock()
		amountProducts = amountProducts + product.Product.Price
		u.mu.Unlock()
	}

	paramProtoFindAccount := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: in.AccountId,
	}

	wallet, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &paramProtoFindAccount)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.GetWalletByAccount: %w", err)
	}

	if wallet.Balance < amountProducts {
		return nil, fmt.Errorf("amount of products is %.2f balance in account %.2f", amountProducts, wallet.Balance)
	}

	paramProtoDebit := protoBalance.ParamDebitWalletRequest{
		WalletID: wallet.WalletID,
		Amount:   amountProducts,
	}

	u.mu.Lock()

	_, err = u.clientBalanceGPRC.Debit(ctx, &paramProtoDebit)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.Debit: %w", err)
	}

	u.mu.Unlock()

	paramProtoCreateOrder := protoOrders.ParamCreateOrderRequest{
		AccountID:   in.AccountId,
		ProductsIDS: in.ProductsIDS,
	}

	orderCreate, err := u.clientOrdersGRPC.Create(ctx, &paramProtoCreateOrder)
	if err != nil {
		paramProtoCredit := protoBalance.ParamCreditWalletRequest{
			WalletID: wallet.WalletID,
			Amount:   amountProducts,
		}

		_, err := u.clientBalanceGPRC.Credit(ctx, &paramProtoCredit)
		if err != nil {
			fmt.Printf("failed to credit account: %v", err)
		}

		return nil, fmt.Errorf("u.clientOrdersGRPC.Create: %w", err)
	}

	out := orders.OrderCreateOutput{
		OrderId:     orderCreate.Order.OrderID,
		AccountId:   orderCreate.Order.AccountID,
		ProductsIDS: orderCreate.Order.ProductsIDS,
		CreatedAt:   orderCreate.Order.CreatedAT.AsTime(),
	}

	return &out, nil
}
func (u *orderUC) FindById(ctx context.Context, in *orders.OrderFindByIdInput) (*orders.OrderFindByIdOutput, error) {

	paramProto := protoOrders.ParamFindOrderRequest{
		OrderID: in.OrderId,
	}

	find, err := u.clientOrdersGRPC.Find(ctx, &paramProto)
	if err != nil {
		return nil, fmt.Errorf("u.clientOrdersGRPC.Find: %w", err)
	}

	out := orders.OrderFindByIdOutput{
		OrderId:     find.Order.OrderID,
		AccountId:   find.Order.AccountID,
		ProductsIDS: find.Order.ProductsIDS,
		CreatedAt:   find.Order.CreatedAT.AsTime(),
	}

	return &out, nil
}
func (u *orderUC) FindByAccount(ctx context.Context, in *orders.OrderByAccountInput) ([]*orders.OrderByAccountOutput, error) {

	protoParamFind := protoOrders.ParamFindOrderByAccountRequest{
		AccountID: in.AccountId,
	}

	ords, err := u.clientOrdersGRPC.FindOrderByAccount(ctx, &protoParamFind)
	if err != nil {
		return nil, fmt.Errorf("u.clientOrdersGRPC.FindOrderByAccount: %w", err)
	}

	var results []*orders.OrderByAccountOutput

	for i := range ords.Orders {
		ord := orders.OrderByAccountOutput{
			OrderId:     ords.Orders[i].OrderID,
			AccountId:   ords.Orders[i].AccountID,
			ProductsIDS: ords.Orders[i].ProductsIDS,
			CreatedAt:   ords.Orders[i].CreatedAT.AsTime(),
		}

		results = append(results, &ord)

	}

	return results, nil
}

func (u *orderUC) FindByProduct(ctx context.Context, in *orders.OrderByProductInput) (*orders.OrderByProductOutput, error) {

	protoParamFind := protoOrders.ParamFindOrderByProductRequest{
		ProductID: in.ProductId,
	}

	find, err := u.clientOrdersGRPC.FindOrderByProduct(ctx, &protoParamFind)
	if err != nil {
		return nil, fmt.Errorf("u.clientOrdersGRPC.FindOrderByProduct: %w", err)
	}

	out := orders.OrderByProductOutput{
		OrderId:     find.Order.OrderID,
		AccountId:   find.Order.AccountID,
		ProductsIDS: find.Order.ProductsIDS,
		CreatedAt:   find.Order.CreatedAT.AsTime(),
	}

	return &out, nil
}
