package repository

import (
	"context"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/payment/pix"
	"github.com/redis/go-redis/v9"
)

type pixRepository struct {
	redis *redis.Client
	timeoutLockGenPix time.Duration
}

func NewPixRepository(timeoutLockGenPix time.Duration, rds *redis.Client) pix.Repository {
	return &pixRepository{
		timeoutLockGenPix: timeoutLockGenPix,
		redis: rds,
	}
}

func (r *pixRepository) Get(ctx context.Context, key string) error {
	return r.redis.Get(ctx, pix.FormatPixKeyRepository(key)).Err()
}
func (r *pixRepository) Set(ctx context.Context, key string) error {
	return r.redis.Set(ctx, pix.FormatPixKeyRepository(key), nil, r.timeoutLockGenPix).Err()
}
