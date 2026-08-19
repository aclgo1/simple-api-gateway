package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/aclgo/simple-api-gateway/internal/orders"
	"github.com/aclgo/simple-api-gateway/internal/payment/pix"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoOrders "github.com/aclgo/simple-api-gateway/proto-service/orders"
	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"
	"github.com/google/uuid"
)

type paymentProcessorPix struct {
	PixAuthorization  string
	repo              pix.Repository
	clientOrdersGRPC  protoOrders.ServiceOrderClient
	clientBalanceGrpc protoBalance.WalletServiceClient
	clientUserGrpc    protoUser.SubscriptionServiceClient
}

func NewpaymentProcessorPix(authorization string, repo pix.Repository, clientOrdersGRPC protoOrders.ServiceOrderClient,
	clientBalanceGrpc protoBalance.WalletServiceClient,
	clientUserGrpc protoUser.SubscriptionServiceClient) *paymentProcessorPix {
	return &paymentProcessorPix{
		PixAuthorization:  authorization,
		repo:              repo,
		clientOrdersGRPC:  clientOrdersGRPC,
		clientBalanceGrpc: clientBalanceGrpc,
		clientUserGrpc:    clientUserGrpc,
	}
}

func (p *paymentProcessorPix) Proccess(ctx context.Context, in *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {

	out := models.ParamPaymentProcessOutput{
		Method:               in.Method,
		Status:               models.PaymentFailed,
		GatewayTransactionID: uuid.NewString(),
	}

	// err := p.repo.Get(ctx, in.AccountId)

	// if err != nil && err != redis.Nil {
	// 	return &out, err
	// }

	// if err == nil {
	// 	return &out, pix.ErrExceddedLimitGenPix
	// }

	// client := &http.Client{
	// 	Timeout: time.Second * 30,
	// }

	// reqBody := fmt.Sprintf(`%s`, "ok")

	// req, err := http.NewRequestWithContext(ctx, "POST", "", strings.NewReader(reqBody))
	// if err != nil {
	// 	return &out, err
	// }

	// req.Header.Add("Content-Type:", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.PixAuthorization))

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return &out, err
	// }

	// respBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return &out, err
	// }

	// fmt.Println(string(respBody))

	// p.repo.Set(ctx, in.AccountId)

	out.GatewayTransactionID = uuid.NewString()
	out.Status = models.PaymentPaid

	return &out, nil
}

func (p *paymentProcessorPix) Webhook(ctx context.Context, in *models.ParamPixWebHookInput) error {

	po := protoOrders.ParamFindOrderByGatewayTransactionIdRequest{
		GatewayTransactionId: in.GatewayTransactionId,
	}

	resp, err := p.clientOrdersGRPC.FindOrderByGatewayTransactionId(ctx, &po)
	if err != nil {
		return fmt.Errorf("p.clientOrdersGRPC.FindOrderByGatewayTransactionId: %w", err)
	}

	order := resp.Order

	if order.Status == protoOrders.OrderStatus_PAID {
		return nil
	}

	switch order.Type {
	case protoOrders.OrderType_BALANCE_DEPOSIT:
		if err := p.processDepositBalance(ctx, order); err != nil {
			return fmt.Errorf("processing balance deposity: %w", err)
		}
	case protoOrders.OrderType_PREMIUM_SUBSCRIPTION:
		if err := p.processSubscription(ctx, order); err != nil {
			return fmt.Errorf("processing premiun subscription: %w", err)
		}
	}

	pupd := protoOrders.ParamUpdateOrderStatusRequest{
		OrderId: order.OrderID,
		Status:  protoOrders.OrderStatus_PAID,
	}

	_, err = p.clientOrdersGRPC.UpdateOrderStatus(ctx, &pupd)
	if err != nil {
		return fmt.Errorf("failed to update order to status paid: %w", err)
	}

	return nil
}

func (p *paymentProcessorPix) processSubscription(ctx context.Context, order *protoOrders.Orders) error {
	var meta orders.ParamsSaveSubscriptionMetadata

	if err := json.Unmarshal(order.Metadata, &meta); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	psub := protoUser.CreateOrExtendSubscriptionRequest{
		UserId: meta.UserId,
		Plan:   meta.Plan,
		Days:   int64(meta.Days),
	}

	_, err := p.clientUserGrpc.CreateOrExtend(ctx, &psub)
	if err != nil {
		return fmt.Errorf("p.clientUserGrpc.CreateOrExtend: %w", err)
	}

	return nil
}

func (p *paymentProcessorPix) processDepositBalance(ctx context.Context, order *protoOrders.Orders) error {

	pf := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: order.AccountID,
	}

	find, err := p.clientBalanceGrpc.GetWalletByAccount(ctx, &pf)
	if err != nil {
		return fmt.Errorf("p.clientBalanceGrpc.GetWalletByAccount: %w", err)
	}

	pb := protoBalance.ParamCreditWalletRequest{
		WalletID:    find.WalletID,
		Amount:      order.Amount,
		ReferenceID: order.GatewayTransactionID,
	}

	_, err = p.clientBalanceGrpc.Credit(ctx, &pb)
	if err != nil {
		return fmt.Errorf("p.clientBalanceGrpc.Credit: %w", err)
	}
	return nil
}
