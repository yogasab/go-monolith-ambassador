package repositories

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	GetValue(ctx context.Context, key string) (string, error)
	SetValue(ctx context.Context, key string, value interface{}) error
}

type redisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) RedisRepository {
	return &redisRepository{rdb: rdb}
}

func (r *redisRepository) GetValue(ctx context.Context, key string) (string, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *redisRepository) SetValue(ctx context.Context, key string, value interface{}) error {
	if err := r.rdb.Set(ctx, key, value, 30*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}
