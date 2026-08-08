package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	"github.com/aclgo/simple-api-gateway/internal/payment/pix"
	"github.com/redis/go-redis/v9"
)

type paymentProcessorPix struct {
	PixAuthorization string
	repo pix.Repository
}

func NewpaymentProcessorPix(authorization string, repo pix.Repository) models.PaymentProcessor {
	return &paymentProcessorPix{
		PixAuthorization: authorization,
		repo: repo,
	}
}


func (p *paymentProcessorPix) Proccess(ctx context.Context, in *models.ParamPaymentProcessInput) (*models.ParamPaymentProcessOutput, error) {

	err := p.repo.Get(ctx, in.AccountId)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == nil {
		return nil, pix.ErrExceddedLimitGenPix
	}
	
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	reqBody := fmt.Sprintf(`%s`, "ok")

	req, err := http.NewRequestWithContext(ctx, "POST", "", strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type:", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.PixAuthorization))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(respBody))

	p.repo.Set(ctx, in.AccountId)

	return &models.ParamPaymentProcessOutput{}, nil
}
