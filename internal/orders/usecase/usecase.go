package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/aclgo/simple-api-gateway/internal/orders"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoOrders "github.com/aclgo/simple-api-gateway/proto-service/orders"
	protoProduct "github.com/aclgo/simple-api-gateway/proto-service/product"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type orderUC struct {
	clientOrdersGRPC   protoOrders.ServiceOrderClient
	clientBalanceGPRC  protoBalance.WalletServiceClient
	clientProductsGRPC protoProduct.ProductServiceClient
	mu                 *sync.Mutex
	logger             logger.Logger
	workerSaga orders.SagaWorker
	gateway orders.PaymentGateway
	subscription orders.SubscriptionInterface
}

func NeworderUC(
	clientOrdersGRPC protoOrders.ServiceOrderClient,
	clientProductsGRPC protoProduct.ProductServiceClient,
	clientBalanceGRPC protoBalance.WalletServiceClient,
	mu *sync.Mutex,
	logger logger.Logger,
	workerSaga orders.SagaWorker,
	gateway orders.PaymentGateway,
	subscription orders.SubscriptionInterface,
) (*orderUC,error) {

	if gateway == nil {
		return nil, errors.New("not configured orders gateways payment")
	}
	return &orderUC{
		clientOrdersGRPC:   clientOrdersGRPC,
		clientBalanceGPRC:  clientBalanceGRPC,
		clientProductsGRPC: clientProductsGRPC,
		mu:                 mu,
		logger:             logger,
		workerSaga:workerSaga,
		gateway: gateway,
		subscription: subscription,
	},nil
}


//version create order v1 simple
func (u *orderUC) Create(ctx context.Context, in *orders.OrderCreateInput) (*orders.OrderCreateOutput, error) {

	var amountProducts int64

	for i := range in.ProductsIDS {
		paramProductProto := protoProduct.ProductFindRequest{
			Id: in.ProductsIDS[i].Id,
		}

		product, err := u.clientProductsGRPC.Find(ctx, &paramProductProto)
		if err != nil {
			return nil, fmt.Errorf("u.clientProductsGRPC.Find: %w", err)
		}
		amountProducts = amountProducts + product.Product.Price
	}

	paramProtoFindAccount := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: in.AccountId,
	}

	wallet, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &paramProtoFindAccount)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.GetWalletByAccount: %w", err)
	}

	if wallet.Balance < amountProducts {
		return nil, fmt.Errorf("amount of products is %d balance in account %d", amountProducts, wallet.Balance)
	}

	refrenceId := uuid.NewString()
	
	paramProtoDebit := protoBalance.ParamDebitWalletRequest{
		WalletID: wallet.WalletID,
		Amount:   amountProducts,
		ReferenceID: refrenceId,
	}

	_, err = u.clientBalanceGPRC.Debit(ctx, &paramProtoDebit)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGPRC.Debit: %w", err)
	}

	for i := range in.ProductsIDS {
		paramProductProto := protoProduct.ProductUpdateRequest{
			Id:         in.ProductsIDS[i].Id,
			HasOrdered: true,
		}

		_, err := u.clientProductsGRPC.Update(ctx, &paramProductProto)
		if err != nil {
			return nil, fmt.Errorf("u.clientProductsGRPC.Find: %w", err)
		}
	}

	metadata, err := json.Marshal(in.ProductsIDS)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w",err)
	}

	paramProtoCreateOrder := protoOrders.ParamCreateOrderRequest{
		AccountID:   in.AccountId,
		Type: protoOrders.OrderType_PRODUCT_PURCHASE,
		PaymentMethod: protoOrders.PaymentMethod_INTERNAL_BALANCE,
		Status: protoOrders.OrderStatus_PAID,
		Metadata: metadata,
	}

	orderCreate, err := u.clientOrdersGRPC.Create(ctx, &paramProtoCreateOrder)
	if err != nil {
		paramProtoCredit := protoBalance.ParamCreditWalletRequest{
			WalletID: wallet.WalletID,
			Amount:   amountProducts,
			ReferenceID: uuid.NewString(),
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
		CreatedAt:   orderCreate.Order.CreatedAT.AsTime(),
	}

	if err := json.Unmarshal(orderCreate.Order.Metadata, &out.ProductsIDS);err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w",err)
	}

	return &out, nil
}

//create order v2 using saga orchestration
func (u *orderUC) CreateWithSaga(ctx context.Context, in *orders.OrderCreateInput) (*orders.OrderCreateOutput, error) {
	var amountProducts int64
	for _, pID := range in.ProductsIDS {
		product, err := u.clientProductsGRPC.Find(ctx, &protoProduct.ProductFindRequest{Id: pID.Id})
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
		return nil, fmt.Errorf("insufficient funds: amount is %d, balance is %d", amountProducts, wallet.Balance)
	}

	var compensations []func(context.Context)error

	rollback := func(originalErr error)error{
		u.workerSaga.AppendTask(&orders.CompensationTask{
			OriginalErr: originalErr,
			Compensations: compensations,
		})
		return originalErr
	}

	referenceId := uuid.NewString()

	_, err = u.clientBalanceGPRC.Debit(ctx, &protoBalance.ParamDebitWalletRequest{
		WalletID: wallet.WalletID,
		Amount:   amountProducts,
		ReferenceID: referenceId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}

	referenceIdCreditCompensate := uuid.NewString()

	compensations = append(compensations, func(ctx context.Context) error {
		_, err := u.clientBalanceGPRC.Credit(ctx, &protoBalance.ParamCreditWalletRequest{
			WalletID: wallet.WalletID,
			Amount:   amountProducts,
			ReferenceID: referenceIdCreditCompensate,
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
			Id:         pID.Id,
			HasOrdered: true,
		})
		if err != nil {
			return nil, rollback(fmt.Errorf("failed to update product %s: %w", pID.Id, err))
		}
		successfullyUpdatedProducts = append(successfullyUpdatedProducts, pID.Id)
	}

	metadata, err := json.Marshal(in.ProductsIDS)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v",err)
	}

	paramProtoCreateOrder := protoOrders.ParamCreateOrderRequest{
		AccountID:   in.AccountId,
		Type: protoOrders.OrderType_PRODUCT_PURCHASE,
		PaymentMethod: protoOrders.PaymentMethod_INTERNAL_BALANCE,
		Status: protoOrders.OrderStatus_PAID,
		Metadata: metadata,
	}

	orderCreate, err := u.clientOrdersGRPC.Create(ctx, &paramProtoCreateOrder)
	if err != nil {
		return nil, rollback(fmt.Errorf("failed to create order: %w", err))
	}

	out := &orders.OrderCreateOutput{
		OrderId:     orderCreate.Order.OrderID,
		AccountId:   orderCreate.Order.AccountID,
		CreatedAt:   orderCreate.Order.CreatedAT.AsTime(),
	}

	if err := json.Unmarshal(orderCreate.Order.Metadata, &out.ProductsIDS);err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w",err)
	}


	return out,nil
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
		CreatedAt:   find.Order.CreatedAT.AsTime(),
	}

	if err := json.Unmarshal(find.Order.Metadata, &out.ProductsIDS);err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w",err)
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
			CreatedAt:   ords.Orders[i].CreatedAT.AsTime(),
		}

		if err := json.Unmarshal(ords.Orders[i].Metadata, &ord.ProductsIDS);err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w",err)
		}

		results = append(results, &ord)

	}

	return results, nil
}

func (u *orderUC) FindByProduct(ctx context.Context, in *orders.OrderByProductInput) ([]*orders.OrderByProductOutput, error) {

	protoParamFind := protoOrders.ParamFindOrderByProductRequest{
		ProductID: in.ProductId,
	}

	finds, err := u.clientOrdersGRPC.FindOrderByProduct(ctx, &protoParamFind)
	if err != nil {
		return nil, fmt.Errorf("u.clientOrdersGRPC.FindOrderByProduct: %w", err)
	}

	outs := make([]*orders.OrderByProductOutput, 0, len(finds.Orders))

	for i := range finds.Orders{
		out := orders.OrderByProductOutput {
			OrderId:     finds.Orders[i].OrderID,
			AccountId:   finds.Orders[i].AccountID,
			CreatedAt:   finds.Orders[i].CreatedAT.AsTime(),
		}

		if err := json.Unmarshal(finds.Orders[i].Metadata, &out.ProductsIDS); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w",err)
		}

		outs = append(outs, &out)
	}

	return outs, nil
}


func (u *orderUC) CreateSubscriptionOrExtend(ctx context.Context,
	params *orders.ParamsCreateOrderSubscriptionInput)(*orders.ParamsCreateOrderSubscriptionOutput,error){

	var amount int64
	switch params.Plan{
	case string(models.Plan_Week):
		amount = 1999
	case string(models.Plan_Month):
		amount = 3499
	case string(models.Plan_Year):
		amount = 28900
	case string(models.Plan_Undefined):
		amount = int64(params.Days) * 199
	default:
		return nil, orders.ErrPlanInvalid
	}

	pg := models.ParamPaymentProcessInput{
		Method: params.MethodPayment,
		AccountId: params.UserId,
		ReferenceId: "",
		Amount: amount,
		CardToken: params.CardToken,
		CardExpiration: params.CardExpiration,
	}

	payment, err := u.gateway.GeneratePayment(ctx, &pg)
	if err != nil {
		return nil,err
	}

	var out orders.ParamsCreateOrderSubscriptionOutput

	status := protoOrders.OrderStatus_PENDING

	switch payment.Status {
	case models.PaymentPaid:
		status = protoOrders.OrderStatus_PAID

		ps := models.ParamsActivateSubscriptionInput{
			AccountID: params.UserId,
			Plan: params.Plan,
			Days: params.Days,
		}


		act, err := u.subscription.ActivateSubscription(ctx, &ps)
		if err != nil {
			return nil, fmt.Errorf("u.subscription.Activate: %w",err)
		}

		out.SubscriptionData = act

	case models.PaymentFailed:
		status = protoOrders.OrderStatus_FAILED
	case models.PaymentPending:
		status = protoOrders.OrderStatus_PENDING
	case models.PaymentCancelled:
		status = protoOrders.OrderStatus_CANCELLED
	case models.PaymentRefunded:
		status = protoOrders.OrderStatus_REFUNDED
	case models.PaymentUnspecified:
		status = protoOrders.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	method := protoOrders.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	switch payment.Method {
	case models.PaymentMethodPix:
		method = protoOrders.PaymentMethod_PIX
	case models.PaymentMethodCard:
		method = protoOrders.PaymentMethod_CREDIT_CARD
	case models.PaymentMethodBoleto:
		method = protoOrders.PaymentMethod_BOLETO
	case models.PaymentMethodInternalBalance:
		method = protoOrders.PaymentMethod_INTERNAL_BALANCE
	}

	var pixExp *timestamppb.Timestamp
	if !payment.PixExpiration.IsZero() {
	    pixExp = timestamppb.New(payment.PixExpiration)
	}

	var boletoExp *timestamppb.Timestamp	
	if !payment.BoletoExpiration.IsZero()	 {
		boletoExp = timestamppb.New(payment.BoletoExpiration)
	}

	metadataObj := orders.ParamsSaveSubscriptionMetadata{
		UserId: params.UserId,
    	Plan :params.Plan,
   	 	Days:params.Days,
	}

	metadataJson, err := json.Marshal(metadataObj)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v\n",err)
	}
	
	paramsNewOrder := protoOrders.ParamCreateOrderRequest{
		AccountID:            params.UserId,
		Type:                 protoOrders.OrderType_PREMIUM_SUBSCRIPTION,
		Status:               status,
		Amount:               amount,
		PaymentMethod:        method,
		Metadata:             metadataJson,
		GatewayTransactionID: payment.GatewayTransactionID,
		PixQRCode:            payment.PixQRCode,
		PixExpiration:        pixExp,
		CardToken:            payment.CardToken,
		CardExpiration:       payment.CardExpiration,
		BoletoURL:            payment.BoletoURL,
		BoletoBarcode:        payment.BoletoBarcode,
		BoletoExpiration:     boletoExp,
	}


	newOrder, err := u.clientOrdersGRPC.Create(ctx, &paramsNewOrder)
	if err != nil {
		return nil,err
	}

	var outPixExp, outBoletoExp time.Time
	if newOrder.Order.PixExpiration != nil {
		outPixExp = newOrder.Order.PixExpiration.AsTime()
	}
	if newOrder.Order.BoletoExpiration != nil {
		outBoletoExp = newOrder.Order.BoletoExpiration.AsTime()
	}

	out = orders.ParamsCreateOrderSubscriptionOutput{
		OrderID:              newOrder.Order.OrderID,
		Status:               newOrder.Order.Status.String(),
		GatewayTransactionID: newOrder.Order.GatewayTransactionID,
		PixQRCode:            newOrder.Order.PixQRCode,
		PixExpiration: outPixExp,
		BoletoURL:            newOrder.Order.BoletoURL,
		BoletoBarcode:        newOrder.Order.BoletoBarcode,
		BoletoExpiration:outBoletoExp,
	}

	return &out,nil
}

func (u *orderUC) AddBalance(ctx context.Context, params *orders.ParamsAddBalanceInput) (*orders.ParamsAddBalanceOutput, error) {

	switch params.Amount {
	case 2500, 5000, 10000:
	default:
		return nil, orders.ErrAmountInvalid
	}

	mp := models.ParamPaymentProcessInput{
		Method:         params.MethodPayment,
		AccountId:      params.UserId,
		Amount:         params.Amount,
		CardToken:      params.CardToken,    
		CardExpiration: params.CardExpiration,
	}

	switch params.MethodPayment{
	case models.PaymentMethodPix, models.PaymentMethodCard, models.PaymentMethodBoleto:
	default:
		return nil, errors.New("method pay invalid")
	}

	payment, err := u.gateway.GeneratePayment(ctx, &mp)
	if err != nil {
		return nil, fmt.Errorf("u.gateway.GeneratePayment: %w", err)
	}

	status := protoOrders.OrderStatus_PENDING

	switch payment.Status {
	case models.PaymentPaid:
		status = protoOrders.OrderStatus_PAID

		pf := protoBalance.ParamGetWalletByAccountRequest{
			AccountID: params.UserId,
		}

		wlt, err := u.clientBalanceGPRC.GetWalletByAccount(ctx, &pf)
		if err != nil {
			return nil, fmt.Errorf("u.clientBalanceGPRC.GetWalletByAccount: %w",err)
		}

		pb := protoBalance.ParamCreditWalletRequest{
			Amount: params.Amount,
			WalletID: wlt.WalletID,
			ReferenceID: payment.GatewayTransactionID,
		}

		_, err = u.clientBalanceGPRC.Credit(ctx, &pb)
		if err != nil {
			return nil, fmt.Errorf("u.clientBalanceGPRC.Credit: %w",err)
		}

	case models.PaymentFailed:
		status = protoOrders.OrderStatus_FAILED
	case models.PaymentPending:
		status = protoOrders.OrderStatus_PENDING
	case models.PaymentCancelled:
		status = protoOrders.OrderStatus_CANCELLED
	case models.PaymentRefunded:
		status = protoOrders.OrderStatus_REFUNDED
	case models.PaymentUnspecified:
		status = protoOrders.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}


	method := protoOrders.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	switch payment.Method {
	case models.PaymentMethodPix:
		method = protoOrders.PaymentMethod_PIX
	case models.PaymentMethodCard:
		method = protoOrders.PaymentMethod_CREDIT_CARD
	case models.PaymentMethodBoleto:
		method = protoOrders.PaymentMethod_BOLETO
	}

	var pixExp *timestamppb.Timestamp
	if !payment.PixExpiration.IsZero() {
	    pixExp = timestamppb.New(payment.PixExpiration)
	}

	var boletoExp *timestamppb.Timestamp	
	if !payment.BoletoExpiration.IsZero()	 {
		boletoExp = timestamppb.New(payment.BoletoExpiration)
	}	

	paramsNewOrder := protoOrders.ParamCreateOrderRequest{
		AccountID:            params.UserId,
		Type:                 protoOrders.OrderType_BALANCE_DEPOSIT,
		Status:               status,
		Amount:               params.Amount,
		PaymentMethod:        method,
		Metadata:             []byte(`{}`),
		GatewayTransactionID: payment.GatewayTransactionID,
		PixQRCode:            payment.PixQRCode,
		PixExpiration:        pixExp,
		CardToken:            payment.CardToken,
		CardExpiration:       payment.CardExpiration,
		BoletoURL:            payment.BoletoURL,
		BoletoBarcode:        payment.BoletoBarcode,
		BoletoExpiration:     boletoExp,
	}

	newOrder, err := u.clientOrdersGRPC.Create(ctx, &paramsNewOrder)
	if err != nil {
		return nil, fmt.Errorf("u.clientOrdersGRPC.Create: %w", err)
	}

	var outPixExp, outBoletoExp time.Time
	if newOrder.Order.PixExpiration != nil {
		outPixExp = newOrder.Order.PixExpiration.AsTime()
	}
	if newOrder.Order.BoletoExpiration != nil {
		outBoletoExp = newOrder.Order.BoletoExpiration.AsTime()
	}


	out := orders.ParamsAddBalanceOutput{
		OrderID:              newOrder.Order.OrderID,
		Status:               newOrder.Order.Status.String(),
		GatewayTransactionID: newOrder.Order.GatewayTransactionID,
		PixQRCode:            newOrder.Order.PixQRCode,
		PixExpiration: outPixExp,
		BoletoURL:            newOrder.Order.BoletoURL,
		BoletoBarcode:        newOrder.Order.BoletoBarcode,
		BoletoExpiration: outBoletoExp,
	}

	return &out, nil
}