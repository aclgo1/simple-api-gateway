package usecase

import (
	"context"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/product"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	proto "github.com/aclgo/simple-api-gateway/proto-service/product"
)

type productUC struct {
	logger            logger.Logger
	clientProductGRPC proto.ProductServiceClient
}

func NewProductUC(logger logger.Logger, clientProductGPRC proto.ProductServiceClient) product.Product {
	return &productUC{
		logger:            logger,
		clientProductGRPC: clientProductGPRC,
	}
}

func (u *productUC) Create(ctx context.Context, in *product.ParamCreateInput) (*product.ParamCreateOutput, error) {

	protoParamCreate := proto.ProductInsertRequest{
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
	}

	create, err := u.clientProductGRPC.Insert(ctx, &protoParamCreate)
	if err != nil {
		return nil, fmt.Errorf("u.clientProductGRPC.Insert: %w", err)
	}

	out := product.ParamCreateOutput{
		Id:          create.Product.Id,
		Name:        create.Product.Name,
		Price:       create.Product.Price,
		Quantity:    create.Product.Quantity,
		Description: create.Product.Description,
		CreatedAt:   create.Product.CreatedAt.AsTime(),
		UpdatedAt:   create.Product.UpdatedAt.AsTime(),
	}

	return &out, nil
}
func (u *productUC) Find(ctx context.Context, in *product.ParamFindInput) (*product.ParamFindOutput, error) {
	protoParamFind := proto.ProductFindRequest{
		Id: in.ProductId,
	}

	find, err := u.clientProductGRPC.Find(ctx, &protoParamFind)
	if err != nil {
		return nil, fmt.Errorf("u.clientProductGRPC.Find: %w", err)
	}

	out := product.ParamFindOutput{
		Id:          find.Product.Id,
		Name:        find.Product.Name,
		Price:       find.Product.Price,
		Quantity:    find.Product.Quantity,
		Description: find.Product.Description,
		CreatedAt:   find.Product.CreatedAt.AsTime(),
		UpdatedAt:   find.Product.UpdatedAt.AsTime(),
	}

	return &out, nil
}
func (u *productUC) FindAll(ctx context.Context, in *product.ParamFindAllInput) ([]*product.ParamFindAllOutput, error) {

	ip := proto.ProductFindAllRequest{}

	all, err := u.clientProductGRPC.FindAll(ctx, &ip)
	if err != nil {
		return nil, fmt.Errorf("u.clientProductGRPC.FindAll: %w", err)
	}

	var products []*product.ParamFindAllOutput

	for i := range all.Products {
		product := product.ParamFindAllOutput{
			Id:          all.Products[i].Id,
			Name:        all.Products[i].Name,
			Price:       all.Products[i].Price,
			Quantity:    all.Products[i].Quantity,
			Description: all.Products[i].Description,
			CreatedAt:   all.Products[i].CreatedAt.AsTime(),
			UpdatedAt:   all.Products[i].UpdatedAt.AsTime(),
		}

		products = append(products, &product)
	}

	return products, nil
}
func (u *productUC) Update(ctx context.Context, in *product.ParamUpdateInput) (*product.ParamUpdateOutput, error) {
	protoParamUpdate := proto.ProductUpdateRequest{
		Id:          in.Id,
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
	}

	updated, err := u.clientProductGRPC.Update(ctx, &protoParamUpdate)
	if err != nil {
		return nil, fmt.Errorf("u.clientProductGRPC.Update: %w", err)
	}

	out := product.ParamUpdateOutput{
		Id:          updated.Product.Id,
		Name:        updated.Product.Name,
		Price:       updated.Product.Price,
		Quantity:    updated.Product.Quantity,
		Description: updated.Product.Description,
		CreatedAt:   updated.Product.CreatedAt.AsTime(),
		UpdatedAt:   updated.Product.UpdatedAt.AsTime(),
	}

	return &out, nil
}
func (u *productUC) Delete(ctx context.Context, in *product.ParamDeleteInput) (*product.ParamDeleteOutput, error) {
	protoParamDelete := proto.ProductDeleteRequest{
		Id: in.Id,
	}

	deletedMessage, err := u.clientProductGRPC.Delete(ctx, &protoParamDelete)
	if err != nil {
		return nil, fmt.Errorf("u.clientProductGRPC.Delete: %w", err)
	}

	out := product.ParamDeleteOutput{
		Message: deletedMessage.Msg,
	}

	return &out, nil
}
