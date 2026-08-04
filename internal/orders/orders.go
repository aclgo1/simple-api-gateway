package orders

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Orders interface {
	Create(context.Context, *OrderCreateInput) (*OrderCreateOutput, error)
	CreateWithSaga(ctx context.Context, in *OrderCreateInput)(*OrderCreateOutput,error)
	CreateSubscriptionOrExtend(ctx context.Context,
		params *ParamsCreateOrderSubscriptionInput)(*ParamsCreateOrderSubscriptionOutput,error)
	FindById(context.Context, *OrderFindByIdInput) (*OrderFindByIdOutput, error)
	FindByAccount(context.Context, *OrderByAccountInput) ([]*OrderByAccountOutput, error)
	FindByProduct(context.Context, *OrderByProductInput) (*OrderByProductOutput, error)
}

type OrderCreateInput struct {
	AccountId   string   `json:"account_id"`
	ProductsIDS []string `json:"products_ids"`
}

func (o *OrderCreateInput) Validate() error {
	if o.AccountId == "" {
		return errors.New("accountId empty")
	}

	if _, err := uuid.Parse(o.AccountId); err != nil {
		return errors.New("invalid uuid account")
	}

	for i := range o.ProductsIDS {
		if _, err := uuid.Parse(o.ProductsIDS[i]); err != nil {
			return errors.New("invalid uuid product")
		}
	}

	return nil
}

type OrderCreateOutput struct {
	OrderId     string    `json:"order_id"`
	AccountId   string    `json:"account_id"`
	ProductsIDS []string  `json:"products_ids"`
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
	ProductsIDS []string  `json:"products_ids"`
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
	ProductsIDS []string  `json:"products_ids"`
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
	ProductsIDS []string  `json:"products_ids"`
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
	SubscriptionId string `json:"subscription_id"` 
	UserId string`json:"user_id"`
	Plan string`json:"plan"`
	Days int `json:"days"`
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
	OrderId string `json:"order_id"`
	AccountId string `json:"account_id"`
	CreatedAt string `json:"created_at"`
	SubscriptionOutput *SubscriptionOutput `"json:"subscription"`
}
