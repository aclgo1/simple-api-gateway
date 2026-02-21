package authUC

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aclgo/simple-api-gateway/internal/auth"

	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"
)

type authUC struct {
	userSvcClient protoUser.UserServiceClient
}

func NewAuthUC(userSvcClient protoUser.UserServiceClient) *authUC {
	return &authUC{
		userSvcClient: userSvcClient,
	}
}

func (a *authUC) validateToken(ctx context.Context, token string) (*auth.ParamsToken, error) {
	resp, err := a.userSvcClient.ValidateToken(
		ctx,
		&protoUser.ValidateTokenRequest{Token: token},
	)

	if err != nil {
		return nil, err
	}

	return &auth.ParamsToken{
		UserID: resp.UserId,
		Role:   resp.UserRole,
	}, nil

}

func getAccessToken(r *http.Request) string {
	accessToken := r.Header.Get(auth.KeyAccessTokenHeader)
	if len(accessToken) < 7 || strings.ToLower(accessToken[:7]) != "bearer " {
		return ""
	}

	return accessToken[7:]
}

func getRefreshToken(r *http.Request) string {
	refreshToken := r.Header.Get(auth.KeyRefreshTokenHeader)
	if len(refreshToken) < 7 || strings.ToLower(refreshToken[:7]) != "bearer " {
		return ""
	}

	return refreshToken[7:]
}

func getPixWebHookToken(r *http.Request) string {
	pixToken := r.Header.Get("pix-token")
	return pixToken
}

func (a *authUC) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		if accessToken == "" {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		paramsToken, err := a.validateToken(context.Background(), accessToken)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		v := context.WithValue(r.Context(), auth.KeyCtxParamsToken, paramsToken)

		next.ServeHTTP(w, r.WithContext(v))

	}
}

func (a *authUC) ValidateTokenWs(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := strings.TrimSpace(r.PathValue(string(auth.KeyQueryTokenValue)))

		if accessToken == "" {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		paramsToken, err := a.validateToken(context.Background(), accessToken)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		v := context.WithValue(r.Context(), auth.KeyCtxParamsToken, paramsToken)

		next.ServeHTTP(w, r.WithContext(v))

	}
}

func (a *authUC) ValidateTwoToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		refreshToken := getRefreshToken(r)

		// _, err := a.validateToken(context.Background(), accessToken)
		// if err != nil {
		// 	resp := auth.Response{
		// 		Error:   http.StatusText(http.StatusUnauthorized),
		// 		Message: err.Error(),
		// 	}

		// 	auth.Json(w, resp, http.StatusUnauthorized)

		// 	return
		// }

		refreshParams := auth.ParamsTwoTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		if err := refreshParams.Validate(); err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		protoParamRefresh := protoUser.RefreshTokensRequest{
			AccessToken:  refreshParams.AccessToken,
			RefreshToken: refreshParams.RefreshToken,
		}

		newRefreshParams, err := a.userSvcClient.RefreshTokens(r.Context(), &protoParamRefresh)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		outTokens := auth.ParamsTwoTokens{
			AccessToken:  newRefreshParams.AccessToken,
			RefreshToken: newRefreshParams.RefreshToken,
		}

		ctx := context.WithValue(r.Context(), auth.KeyCtxParamsRefreshToken, &outTokens)

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

func (a *authUC) ValidateUpdate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		if accessToken == "" {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		paramsToken, err := a.validateToken(context.Background(), accessToken)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		params := auth.ParamsUpdate{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusBadRequest)

			return
		}

		if params.IdUpdate == "" {
			params.IdUpdate = paramsToken.UserID
		}

		isOwner := params.IdUpdate == paramsToken.UserID
		isAdmin := paramsToken.Role == string(auth.ADMIN)
		isSuper := paramsToken.Role == string(auth.SUPERADMIN)

		if !isOwner && !isAdmin && !isSuper {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusForbidden),
				Message: "usuario nao autorizado fazer esse update",
			}

			auth.Json(w, resp, http.StatusForbidden)

			return
		}

		if !isAdmin && !isSuper {
			params.Balance = 0
		}

		v := context.WithValue(r.Context(), auth.KeyCtxParamsUpdate, &params)

		next.ServeHTTP(w, r.WithContext(v))
	}
}
func (a *authUC) ValidateCreateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		if accessToken == "" {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		paramsTokenLogged, err := a.validateToken(context.Background(), accessToken)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		if paramsTokenLogged.Role != string(auth.SUPERADMIN) && paramsTokenLogged.Role != string(auth.ADMIN) {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: "role not super-admin or admin",
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		params := auth.ParamsCreateAdmin{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusBadRequest),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusBadRequest)

			return
		}

		var finalRole string

		switch params.Role {
		case string(auth.ADMIN):
			if paramsTokenLogged.Role == string(auth.SUPERADMIN) {
				finalRole = string(auth.ADMIN)
			} else {
				finalRole = string(auth.CLIENT)
			}
		default:
			finalRole = string(auth.CLIENT)
		}

		params.Role = finalRole

		v := context.WithValue(r.Context(), auth.KeyCtxParamsCreateAdmin, &params)

		next.ServeHTTP(w, r.WithContext(v))

	}
}
func (a *authUC) ValidateIsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		if accessToken == "" {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: auth.ErrInvalidToken{}.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		paramsToken, err := a.validateToken(context.Background(), accessToken)
		if err != nil {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		isAdmin := paramsToken.Role == string(auth.ADMIN)
		isSuperAdmin := paramsToken.Role == string(auth.SUPERADMIN)

		if !isAdmin && !isSuperAdmin {
			resp := auth.Response{
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: "role not admin or super",
			}

			auth.Json(w, resp, http.StatusUnauthorized)

			return
		}

		v := context.WithValue(r.Context(), auth.KeyCtxParamsToken, paramsToken)

		next.ServeHTTP(w, r.WithContext(v))

	}
}

func (a *authUC) ValidateWebHookPix(ctxf context.Context, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		stop := context.AfterFunc(ctxf, func() { cancel() })
		defer stop()

		token := getPixWebHookToken(r)

		v := context.WithValue(ctx, "teste", token)

		r = r.WithContext(v)

		next.ServeHTTP(w, r)
	}
}
