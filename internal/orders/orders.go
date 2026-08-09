package orders

import (
	"context"
	"errors"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/google/uuid"
)



var (
	ErrAmountInvalid =  errors.New("amount invalid")
	ErrPlanInvalid = errors.New("plan invalid")
)

type Orders interface {
	Create(context.Context, *OrderCreateInput) (*OrderCreateOutput, error)
	CreateWithSaga(ctx context.Context, in *OrderCreateInput)(*OrderCreateOutput,error)
	CreateSubscriptionOrExtend(ctx context.Context,
		params *ParamsCreateOrderSubscriptionInput)(*ParamsCreateOrderSubscriptionOutput,error)
	FindById(context.Context, *OrderFindByIdInput) (*OrderFindByIdOutput, error)
	FindByAccount(context.Context, *OrderByAccountInput) ([]*OrderByAccountOutput, error)
	FindByProduct(context.Context, *OrderByProductInput) ([]*OrderByProductOutput, error)
	AddBalance(ctx context.Context, params *ParamsAddBalanceInput)(*ParamsAddBalanceOutput,error)
}

type PaymentGateway interface {
	GeneratePayment(ctx context.Context, params *models.ParamPaymentProcessInput)(*models.ParamPaymentProcessOutput,error)
}

type SubscriptionInterface interface {
	ActivateSubscription(context.Context, *models.ParamsActivateSubscriptionInput)(*models.ParamsActivateSubscriptionOutput,error)
}

type ParamPaymentProcessInput struct {

}

type ParamPaymentProccessOutput struct{

}

type ProductItem struct {
	Id string`json:"product_id"`
}

type OrderCreateInput struct {
	AccountId   string   `json:"account_id"`
	ProductsIDS []ProductItem `json:"products"`
}

func (o *OrderCreateInput) Validate() error {
	if o.AccountId == "" {
		return errors.New("accountId empty")
	}

	if _, err := uuid.Parse(o.AccountId); err != nil {
		return errors.New("invalid uuid account")
	}

	for i := range o.ProductsIDS {
		if _, err := uuid.Parse(o.ProductsIDS[i].Id); err != nil {
			return errors.New("invalid uuid product")
		}
	}

	return nil
}

type OrderCreateOutput struct {
	OrderId     string    `json:"order_id"`
	AccountId   string    `json:"account_id"`
	ProductsIDS []ProductItem  `json:"products"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderFindByIdInput struct {
	OrderId string `json:"order_id"`
}

func (o *OrderFindByIdInput) Validate() error {
	if o.OrderId == "" {
		return errors.New("order id empty")
	}

	if _, err := uuid.Parse(o.OrderId); err != nil {
		return errors.New("invalid uuid order")
	}

	return nil
}

type OrderFindByIdOutput struct {
	OrderId     string    `json:"order_id"`
	AccountId   string    `json:"account_id"`
	ProductsIDS []ProductItem  `json:"products"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderByAccountInput struct {
	AccountId string `json:"account_id"`
}

func (o *OrderByAccountInput) Validate() error {
	if o.AccountId == "" {
		return errors.New("account id empty")
	}

	if _, err := uuid.Parse(o.AccountId); err != nil {
		return errors.New("invalid uuid account")
	}
	return nil
}

type OrderByAccountOutput struct {
	OrderId     string    `json:"order_id"`
	AccountId   string    `json:"account_id"`
	ProductsIDS []ProductItem  `json:"products"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderByProductInput struct {
	ProductId string `json:"product_id"`
}

func (o *OrderByProductInput) Validate() error {
	if o.ProductId == "" {
		return errors.New("product id empty")
	}

	if _, err := uuid.Parse(o.ProductId); err != nil {
		return errors.New("invalid uuid product")
	}
	return nil
}

type OrderByProductOutput struct {
	OrderId     string    `json:"order_id"`
	AccountId   string    `json:"account_id"`
	ProductsIDS []string  `json:"products"`
	CreatedAt   time.Time `json:"created_at"`
}

type SagaWorker interface {
	AppendTask(task *CompensationTask)
}

type CompensationTask struct {
	Compensations []func(context.Context)error
	OriginalErr error
}


type ParamsCreateOrderSubscriptionInput struct{
	MethodPayment string `json:"method_payment"`
	UserId string`json:"user_id"`
	Plan string`json:"plan"`
	Days int `json:"days"`
	CardToken string `json:"card_token"`
	CardExpiration string `json:"card_expiration"`
}

func(p *ParamsCreateOrderSubscriptionInput)Validate()error{
	return nil
}


type SubscriptionOutput struct{
	Id string `json:"subscription_id"`
	UserId string `json:"user_id"`
	Plan string `json:"plan"`
	Status string `json:"status"`
	StartsAt string `json:"starts_at"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ParamsCreateOrderSubscriptionOutput struct{
	OrderID string `json:"order_id"`
	Status string `json:"status"`
	SubscriptionData *models.ParamsActivateSubscriptionOutput `subscription_data`
	GatewayTransactionID string `json:"gateway_transaction_id"`
	PixQRCode string `json:"pix_qr_code"`
	PixExpiration time.Time `json:"pix_expiration"`
	BoletoURL string `json:"boleto_url"`
	BoletoBarcode string `json:"boleto_bar_code"`
	BoletoExpiration	time.Time `json:"boleto_expiration"`
}

type ParamsAddBalanceInput struct {
	MethodPayment string `json:"method_payment"`
	UserId string `json:"user_id"`
	Amount int64 `json:"amount"`
	CardToken string `json:"card_token"`
	CardExpiration string `json:"card_expiration"`
}

func(p *ParamsAddBalanceInput)Validate()error{
	return nil
}

type ParamsAddBalanceOutput struct {
	OrderID string
	Status string
	GatewayTransactionID string
	PixQRCode string
	PixExpiration time.Time
	BoletoURL string
	BoletoBarcode string
	BoletoExpiration	time.Time
}

type ParamsSaveSubscriptionMetadata struct {
    UserId        string `json:"user_id"`
    Plan          string `json:"plan"`
    Days          int    `json:"days"`
}