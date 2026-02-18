package captcha

import (
	"context"

	"github.com/mojocn/base64Captcha"
)

type CaptchaInterface interface {
	GenerateDriverDigit(context.Context, *ParamDriverDigitInput) (*ParamDriverDigitOutput, error)
	VerrifyDriverDigit(ctx context.Context, i *ParamVerifyDriverDigitInput,
	) (*ParamVerifyDriverDigitOutput, error)
}

type Repository interface {
	base64Captcha.Store
}

type ParamDriverDigitInput struct{}

type ParamDriverDigitOutput struct {
	Id          string `json:"id"`
	Base64Image string `json:"base64_image"`
}

type ParamVerifyDriverDigitInput struct {
	Id          string `json:"id"`
	VerifyValue string `json:"verify_value"`
}
type ParamVerifyDriverDigitOutput struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}
