package rediscache

import (
	"context"
	"errors"
	"fazil-syed/gofinance/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type RedisClient struct {
	cfg    *config.RedisConfig
	client *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
	})
	r := &RedisClient{
		cfg:    cfg,
		client: client,
	}
	return r
}

func (r *RedisClient) GetData(ctx context.Context, key string) ([]byte, error) {
	cached, err := r.client.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, err
	}
	return cached, nil

}
func (r *RedisClient) SetData(ctx context.Context, key string, data []byte, ttl time.Duration) error {

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}
	return nil

}
