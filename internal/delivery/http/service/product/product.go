package product

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/internal/product"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
)

type productService struct {
	producUC product.Product
	logger   logger.Logger
}

func NewProductService(productUC product.Product, logger logger.Logger) *productService {
	return &productService{
		producUC: productUC,
		logger:   logger,
	}
}

func (s *productService) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var param product.ParamCreateInput

		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusMethodNotAllowed)
			return
		}

		create, err := s.producUC.Create(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, create, http.StatusOK)

	}
}
func (s *productService) Find(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		productID := r.PathValue("product_id")

		param := product.ParamFindInput{
			ProductId: productID,
		}

		if err := param.Validate(); err != nil {
			if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
				response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

				service.JSON(w, response, http.StatusBadRequest)
				return
			}
		}

		find, err := s.producUC.Find(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, find, http.StatusOK)

	}
}
func (s *productService) FindAll(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var param product.ParamFindAllInput

		if err := param.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		all, err := s.producUC.FindAll(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, all, http.StatusOK)

	}
}
func (s *productService) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var param product.ParamUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusMethodNotAllowed)
			return
		}

		updated, err := s.producUC.Update(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, updated, http.StatusOK)

	}
}
func (s *productService) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		productID := r.PathValue("product_id")

		param := product.ParamDeleteInput{Id: productID}

		if err := param.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		message, err := s.producUC.Delete(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, message, http.StatusOK)

	}
}
