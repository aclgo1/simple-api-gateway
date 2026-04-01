package usecase

import (
	"context"
	"fmt"
	"image/color"

	"github.com/aclgo/simple-api-gateway/internal/captcha"
	"github.com/mojocn/base64Captcha"
)

type captchaUC struct {
	repo captcha.Repository
}

func NewCaptchaUC(repo captcha.Repository) captcha.CaptchaInterface {
	return &captchaUC{
		repo: repo,
	}
}

func (c *captchaUC) GenerateDriverDigit(ctx context.Context, i *captcha.ParamDriverDigitInput,
) (*captcha.ParamDriverDigitOutput, error) {

	driver := &base64Captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      50,
		ShowLineOptions: 2,
		Length:          5,
		Source:          "1234567890",
		BgColor:         &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Fonts:           []string{"wqy-microhei.ttc"},
	}

	new := base64Captcha.NewCaptcha(driver, c.repo)
	id, b64s, _, err := new.Generate()
	if err != nil {
		return nil, fmt.Errorf("new.Generate: %w", err)
	}

	out := captcha.ParamDriverDigitOutput{
		Id:          id,
		Base64Image: b64s,
	}

	return &out, nil
}

func (c *captchaUC) VerrifyDriverDigit(ctx context.Context, i *captcha.ParamVerifyDriverDigitInput,
) (*captcha.ParamVerifyDriverDigitOutput, error) {
	if c.repo.Verify(i.Id, i.VerifyValue, true) {
		out := captcha.ParamVerifyDriverDigitOutput{
			Code: 1,
			Msg:  "ok",
		}

		return &out, nil
	}

	out := captcha.ParamVerifyDriverDigitOutput{
		Code: 0,
		Msg:  "failed",
	}

	return &out, nil
}
