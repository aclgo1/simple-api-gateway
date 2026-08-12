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

		var params orders.ParamsCreateOrderAction

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return
		}

		switch params.Action {
		case string(orders.AddBalance):
			var input orders.ParamsAddBalanceInput
			if err := json.Unmarshal(params.Payload, &input); err != nil {
				response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
				service.JSON(w, response, http.StatusBadRequest)
				return
			}

			added, err := s.ordersUC.AddBalance(r.Context(), &input)
			if err != nil {
				response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
				service.JSON(w, response, http.StatusInternalServerError)
				return
			}

			service.JSON(w, added, http.StatusOK)

		case string(orders.BuyProduct):
			var input orders.ParamBuyProductInput
			if err := json.Unmarshal(params.Payload, &input); err != nil {
				response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
				service.JSON(w, response, http.StatusBadRequest)
				return
			}

			buyed, err := s.ordersUC.CreateWithSaga(r.Context(), &input)
			if err != nil {
				response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
				service.JSON(w, response, http.StatusInternalServerError)
				return
			}

			service.JSON(w, buyed, http.StatusOK)

		case string(orders.NewSubscription):
			var input orders.ParamsCreateOrderSubscriptionInput
			if err := json.Unmarshal(params.Payload, &input); err != nil {
				response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
				service.JSON(w, response, http.StatusBadRequest)
				return
			}

			subscription, err := s.ordersUC.CreateSubscriptionOrExtend(r.Context(), &input)
			if err != nil {
				response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
				service.JSON(w, response, http.StatusInternalServerError)
				return
			}

			service.JSON(w, subscription, http.StatusOK)

		default:
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), "action not suported")
			service.JSON(w, response, http.StatusBadRequest)
		}
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
