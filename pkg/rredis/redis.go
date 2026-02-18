package rredis

import (
	"context"
	"log"

	"github.com/aclgo/simple-api-gateway/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
		PoolSize: 10000,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("client.Ping().Err(): %v", err)
	}

	return client
}
