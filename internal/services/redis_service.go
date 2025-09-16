package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type redisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) RedisService {
	return &redisService{client: client}
}

func (r *redisService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisService) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return res, err
}

func (r *redisService) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
