package auth

import "net/http"

type Auth interface {
	ValidateToken(next http.HandlerFunc) http.HandlerFunc
	ValidateTokenWs(next http.HandlerFunc) http.HandlerFunc
	ValidateTwoToken(next http.HandlerFunc) http.HandlerFunc
	ValidateUpdate(next http.HandlerFunc) http.HandlerFunc
	ValidateCreateAdmin(next http.HandlerFunc) http.HandlerFunc
	ValidateIsAdmin(next http.HandlerFunc) http.HandlerFunc
}

type (
	Level                string
	CtxParamsUpdate      string
	CtxParamsCreateAdmin string
	CtxParamsToken       string
	CtxParamsTwoTokens   string
	KeyQueryToken        string
)

var (
	SUPERADMIN Level = "super-admin"
	ADMIN      Level = "admin"
	CLIENT     Level = "client"

	KeyCtxParamsUpdate       CtxParamsUpdate      = "params-update"
	KeyCtxParamsCreateAdmin  CtxParamsCreateAdmin = "params-create-admin"
	KeyCtxParamsToken        CtxParamsToken       = "params-token"
	KeyCtxParamsRefreshToken CtxParamsTwoTokens   = "params-two-tokens"
	KeyAccessTokenHeader                          = "access-token"
	KeyRefreshTokenHeader                         = "refresh-token"
	KeyQueryTokenValue       KeyQueryToken        = "token"
)

type ParamsUpdate struct {
	ParamsToken
	IdUpdate string  `json:"user_id"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

type ParamsCreateAdmin struct {
	Name          string  `json:"name"`
	Lastname      string  `json:"lastname"`
	Password      string  `json:"password"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Verified      string  `json:"verified"`
	Balance       float64 `json:"balance"`
	CaptchaId     string  `json:"captcha_id"`
	CaptchaAwnser string  `json:"captcha_awnser"`
}

type ParamsToken struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type ParamsTwoTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrEmptyToken struct {
}

func (e ErrEmptyToken) Error() string {
	return "empty token"
}

func (p *ParamsTwoTokens) Validate() error {
	if p.AccessToken == "" {
		return ErrEmptyToken{}
	}

	if p.RefreshToken == "" {
		return ErrEmptyToken{}
	}

	return nil
}
