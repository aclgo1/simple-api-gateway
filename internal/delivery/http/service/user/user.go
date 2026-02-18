package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/auth"
	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
	"github.com/aclgo/simple-api-gateway/internal/user"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
)

type userService struct {
	userUC user.UserUC
	logger logger.Logger
}

func NewuserService(user user.UserUC, logger logger.Logger) *userService {
	return &userService{
		userUC: user,
		logger: logger,
	}
}

func (s *userService) Register(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params user.ParamsUserRegister

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		created, err := s.userUC.Register(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		err = s.userUC.SendConfirm(
			ctx,
			&user.ParamsConfirm{
				To:           created.Email,
				IntervalSend: time.Hour,
			},
		)

		if err != nil {

			//SE ERROR SEND EMAIL VERIFICACAO, DELETA O USER CRIADO PARA NAO DA CONFLITO COM O EMAIL
			errCancel := s.userUC.Delete(ctx, &user.ParamsUserDelete{
				UserID: created.UserID,
			})

			if errCancel != nil {
				response := service.NewRestError(http.StatusText(http.StatusInternalServerError), service.ErrSendEmailAndCancelNewRegister.Error())
				service.JSON(w, response, http.StatusInternalServerError)
				return
			}

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

func (s *userService) Login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := user.ParamsUserLoginRequest{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		logged, err := s.userUC.Login(ctx, &params)
		if err != nil {

			var response *service.RestError

			if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
				response = service.NewRestError(http.StatusText(http.StatusNotFound), "user not exist")
			} else {
				response = service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			}

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": "user logged",
			"tokens":  logged,
		}

		service.JSON(w, resp, http.StatusOK)
	}
}

func (s *userService) Logout(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramsRefresh, ok := r.Context().Value(auth.KeyCtxParamsRefreshToken).(*auth.ParamsTwoTokens)
		if !ok {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), service.ErrNoParamsInCtx.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		params := user.ParamsUserLogout{
			AccessToken:  paramsRefresh.AccessToken,
			RefreshToken: paramsRefresh.RefreshToken,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := s.userUC.Logout(ctx, &params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": "user logout",
		}

		service.JSON(w, resp, http.StatusOK)
	}
}

func (s *userService) Find(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramsTtk, ok := r.Context().Value(auth.KeyCtxParamsToken).(*auth.ParamsToken)

		if !ok {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), service.ErrNoParamsInCtx.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		params := user.ParamsUserFindById{
			UserID: paramsTtk.UserID,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		found, err := s.userUC.FindById(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": "user found",
			"user":    found,
		}

		service.JSON(w, resp, http.StatusOK)
	}
}

func (s *userService) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxData, ok := r.Context().Value(auth.KeyCtxParamsUpdate).(*auth.ParamsUpdate)
		if !ok {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), service.ErrNoParamsInCtx.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		params := user.ParamsUserUpdate{
			UserID:   ctxData.IdUpdate,
			Name:     ctxData.Name,
			Lastname: ctxData.Lastname,
			Password: ctxData.Password,
			Email:    ctxData.Email,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		updated, err := s.userUC.Update(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]any{
			"message": "user updated",
			"user":    updated,
		}

		service.JSON(w, resp, http.StatusOK)

	}
}

func (s *userService) ValidToken(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.JSON(w, nil, http.StatusOK)
	}
}

func (s *userService) UserConfirm(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := r.PathValue("confirm_code")

		params := user.ParamsConfirmOK{ConfirmCode: code}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := s.userUC.SendConfirmOK(ctx, &params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		// resp := map[string]string{
		// 	"message": "user confirmed signup",
		// }

		http.Redirect(w, r, "/login", http.StatusOK)
	}
}

func (s *userService) UserResetPass(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.PathValue("email")
		cptId := r.PathValue("captcha_id")
		cptAwnser := r.PathValue("captcha_awnser")

		params := user.ParamsResetPass{
			Email:         email,
			CaptchaId:     cptId,
			CaptchaAwnser: cptAwnser,
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)
			return

		}

		err := s.userUC.ResetPass(ctx, &params)
		if err != nil {
			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())

			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]string{
			"message": "code to reset pass send to email",
		}

		service.JSON(w, resp, http.StatusOK)
	}
}

func (s *userService) UserNewPass(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := r.PathValue("code")

		params := user.ParamsNewPass{
			NewPassCode: code,
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := params.Validate(); err != nil {
			response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())

			service.JSON(w, response, http.StatusBadRequest)

			return
		}

		if err := s.userUC.NewPass(ctx, &params); err != nil {
			if err == (user.ErrInvalidCode{}) {
				response := service.NewRestError(http.StatusText(http.StatusBadRequest), err.Error())
				service.JSON(w, response, http.StatusBadRequest)

				return
			}

			response := service.NewRestError(http.StatusText(http.StatusInternalServerError), err.Error())
			service.JSON(w, response, http.StatusInternalServerError)

			return
		}

		resp := map[string]string{
			"message": "user updated pass",
		}

		service.JSON(w, resp, http.StatusOK)

	}
}

func (s *userService) RefreshTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshTokens, ok := r.Context().Value(auth.KeyCtxParamsRefreshToken).(*auth.ParamsTwoTokens)
		if !ok {
			resp := service.NewRestError(http.StatusText(http.StatusInternalServerError), service.ErrNoParamsInCtx.Error())

			service.JSON(w, resp, http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"message": "tokens refreshed",
			"tokens": map[string]any{
				"access_token":  refreshTokens.AccessToken,
				"refresh_token": refreshTokens.RefreshToken,
			},
		}

		service.JSON(w, resp, http.StatusOK)
	}
}
