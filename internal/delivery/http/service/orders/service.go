package orders

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/auth"
	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/internal/orders"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
)

type ordersService struct {
	ordersUC orders.Orders
	logger   logger.Logger
}

func NewOrdersService(ordersUC orders.Orders, logger logger.Logger) *ordersService {
	return &ordersService{
		ordersUC: ordersUC,
		logger:   logger,
	}
}

func (s *ordersService) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var param orders.OrderCreateInput

		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		if err := param.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		ord, err := s.ordersUC.Create(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, ord, http.StatusCreated)

	}
}
func (s *ordersService) FindById(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("order_id")

		param := orders.OrderFindByIdInput{
			OrderId: id,
		}

		if err := param.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		find, err := s.ordersUC.FindById(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, find, http.StatusOK)

	}
}
func (s *ordersService) FindByAccount(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramTtk := r.Context().Value(auth.KeyCtxParamsToken).(*auth.ParamsToken)

		param := orders.OrderByAccountInput{
			AccountId: paramTtk.UserID,
		}

		if err := param.Validate(); err != nil {
			return
		}

		ordersByAccount, err := s.ordersUC.FindByAccount(ctx, &param)
		if err != nil {
			return
		}

		service.JSON(w, ordersByAccount, http.StatusOK)

	}
}
func (s *ordersService) FindByProduct(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("product_id")

		param := orders.OrderByProductInput{
			ProductId: id,
		}

		if err := param.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		find, err := s.ordersUC.FindByProduct(ctx, &param)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)
			return
		}

		service.JSON(w, find, http.StatusOK)

	}
}
