package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/wallet"
	"github.com/aclgo/simple-api-gateway/internal/wallet/pix"
)

type paymentProcessorPix struct {
}

func NewpaymentProcessorPix() pix.PaymentProcessor {
	return &paymentProcessorPix{}
}

type DataProccessResponse struct{}

func (p *paymentProcessorPix) Proccess(ctx context.Context, in *wallet.ParamPaymentProcessorInput) (any, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	reqBody := fmt.Sprintf(`%s`, "ok")

	req, err := http.NewRequestWithContext(ctx, "POST", "", strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type:", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data DataProccessResponse
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", data)

	return &pix.ParamsPixOutput{Teste: "meu pix"}, nil
}
