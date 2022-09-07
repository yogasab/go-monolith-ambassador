package services

import (
	"context"

	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type RedisService interface {
	GetValue(ctx context.Context, key string) (string, error)
	SetValue(ctx context.Context, key string, value interface{}) (bool, error)
}

type redisService struct {
	redisRepository repositories.RedisRepository
}

func NewRedisService(redisRepository repositories.RedisRepository) RedisService {
	return &redisService{redisRepository: redisRepository}
}

func (s *redisService) GetValue(ctx context.Context, key string) (string, error) {
	result, err := s.redisRepository.GetValue(ctx, key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *redisService) SetValue(ctx context.Context, key string, value interface{}) (bool, error) {
	if err := s.redisRepository.SetValue(ctx, key, value); err != nil {
		return false, err
	}
	return true, nil
}
