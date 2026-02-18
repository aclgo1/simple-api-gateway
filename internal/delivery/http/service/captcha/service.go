package captcha

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aclgo/simple-api-gateway/internal/captcha"
	"github.com/aclgo/simple-api-gateway/internal/delivery/http/service"
)

type captchaService struct {
	CaptchaUC captcha.CaptchaInterface
}

func NewCaptchaService(captchaUC captcha.CaptchaInterface) *captchaService {
	return &captchaService{
		CaptchaUC: captchaUC,
	}
}

func (s captchaService) GenCaptcha(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		param := captcha.ParamDriverDigitInput{}
		out, err := s.CaptchaUC.GenerateDriverDigit(ctx, &param)
		if err != nil {
			log.Fatal(err)
			return
		}

		service.JSON(w, out, http.StatusCreated)

	}
}

func (s captchaService) ValidateDriverDigit(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		param := captcha.ParamVerifyDriverDigitInput{}

		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			log.Fatal(err)
			return
		}

		out, err := s.CaptchaUC.VerrifyDriverDigit(ctx, &param)
		if err != nil {
			log.Fatal(err)
			return
		}

		service.JSON(w, out, http.StatusCreated)

	}
}
