package service

import (
	"context"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/admin"
	"github.com/aclgo/simple-api-gateway/internal/auth"
	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
)

type adminService struct {
	adminUC admin.AdminUC
	logger  logger.Logger
}

func NewadminService(adminUC admin.AdminUC, logger logger.Logger) *adminService {
	return &adminService{
		adminUC: adminUC,
		logger:  logger,
	}
}

func (s *adminService) Create(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxData, ok := r.Context().Value(auth.KeyCtxParamsCreateAdmin).(*auth.ParamsCreateAdmin)
		if !ok {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), service.ErrNoParamsInCtx.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		params := admin.ParamsCreateAdmin{
			Name:          ctxData.Name,
			Lastname:      ctxData.Lastname,
			Password:      ctxData.Password,
			Email:         ctxData.Email,
			Role:          ctxData.Role,
			Verified:      ctxData.Verified,
			Balance:       ctxData.Balance,
			CaptchaId:     ctxData.CaptchaId,
			CaptchaAwnser: ctxData.CaptchaAwnser,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		created, err := s.adminUC.Create(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": "user created",
			"user":    created,
		}

		service.JSON(w, resp, http.StatusOK)
	}
}

func (s *adminService) FindUsers(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		params := admin.ParamsSearch{
			Page:  page,
			Limit: limit,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		search, err := s.adminUC.Search(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		service.JSON(w, search, http.StatusOK)
	}
}

func (s *adminService) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idDelete := r.PathValue("user_id")

		i := admin.ParamsDeleteUser{
			UserId: idDelete,
		}

		if err := i.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		msg, err := s.adminUC.Delete(ctx, &i)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": msg,
		}

		service.JSON(w, resp, http.StatusOK)
	}
}
