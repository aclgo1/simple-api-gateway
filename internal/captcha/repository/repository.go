package repository

import (
	"time"

	"github.com/aclgo/simple-api-gateway/internal/captcha"
	"github.com/mojocn/base64Captcha"
)

func NewRepository() captcha.Repository {
	return base64Captcha.NewMemoryStore(10240, 10*time.Minute)
}
