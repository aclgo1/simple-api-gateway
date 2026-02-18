package orders

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Orders interface {
	Create(context.Context, *OrderCreateInput) (*OrderCreateOutput, error)
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
