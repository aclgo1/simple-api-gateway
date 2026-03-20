package repository

import (
	"context"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/wallet/pix"
	"github.com/redis/go-redis/v9"
)

type pixRepository struct {
	redis *redis.Client
}

func NewPixRepository(rds *redis.Client) pix.Repository {
	return &pixRepository{
		redis: rds,
	}
}

func (r *pixRepository) Get(ctx context.Context, key string) error {
	return r.redis.Get(ctx, pix.FormatPixKeyRepository(key)).Err()
}
func (r *pixRepository) Set(ctx context.Context, key string) error {
	return r.redis.Set(ctx, pix.FormatPixKeyRepository(key), nil, time.Minute*30).Err()
}
