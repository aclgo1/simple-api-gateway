package product

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Product interface {
	Create(context.Context, *ParamCreateInput) (*ParamCreateOutput, error)
	Find(context.Context, *ParamFindInput) (*ParamFindOutput, error)
	FindAll(context.Context, *ParamFindAllInput) ([]*ParamFindAllOutput, error)
	Update(context.Context, *ParamUpdateInput) (*ParamUpdateOutput, error)
	Delete(context.Context, *ParamDeleteInput) (*ParamDeleteOutput, error)
}

type ParamCreateInput struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
}

func (p *ParamCreateInput) Validate() error {
	if p.Name == "" {
		return errors.New("name product empty")
	}

	if p.Price <= 0 {
		return errors.New("price product invalid")
	}

	if p.Quantity <= 0 {
		return errors.New("quantity product invalid")
	}

	if p.Description == "" {
		return errors.New("description product empty")
	}

	return nil
}

type ParamCreateOutput struct {
	Id          string    `json:"product_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParamFindInput struct {
	ProductId string `json:"product_id"`
}

func (p *ParamFindInput) Validate() error {
	return nil
}

type ParamFindOutput struct {
	Id          string    `json:"product_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParamFindAllInput struct {
}

func (p *ParamFindAllInput) Validate() error {
	return nil
}

type ParamFindAllOutput struct {
	Id          string    `json:"product_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParamUpdateInput struct {
	Id          string  `json:"product_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
}

func (p *ParamUpdateInput) Validate() error {
	if p.Id == "" {
		return errors.New("product id empty")
	}

	if _, err := uuid.Parse(p.Id); err != nil {
		return errors.New("product uuid invalid")
	}

	return nil
}

type ParamUpdateOutput struct {
	Id          string    `json:"product_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParamDeleteInput struct {
	Id string `json:"product_id"`
}

func (p *ParamDeleteInput) Validate() error {
	return nil
}

type ParamDeleteOutput struct {
	Message string
}
