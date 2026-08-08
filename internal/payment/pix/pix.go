package pix

import (
	"context"
	"errors"
	"fmt"
)

type Repository interface {
	Get(ctx context.Context, key string) error
	Set(ctx context.Context, key string) error
}

func FormatPixKeyRepository(id string) string {
	return fmt.Sprintf("pix_generated:%s", id)
}

type ParamsPixOutput struct {
	Teste string `json:"teste"`
}

var (
	ErrExceddedLimitGenPix = errors.New("excedded limit gen pix")
)