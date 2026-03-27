package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/passwordmanager/backend/internal/common/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *RedisClient) SetRefreshToken(ctx context.Context, tokenID string, userID string, expiration time.Duration) error {
	return r.Set(ctx, "refresh:"+tokenID, userID, expiration)
}

func (r *RedisClient) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	return r.Get(ctx, "refresh:"+tokenID)
}

func (r *RedisClient) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	return r.Delete(ctx, "refresh:"+tokenID)
}

func (r *RedisClient) RevokeAllUserTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("refresh:*:%s", userID)
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.Client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (r *RedisClient) IncrementLoginAttempts(ctx context.Context, email string) (int64, error) {
	key := "login_attempts:" + email
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	r.Client.Expire(ctx, key, 15*time.Minute)
	return count, nil
}

func (r *RedisClient) GetLoginAttempts(ctx context.Context, email string) (int64, error) {
	key := "login_attempts:" + email
	count, err := r.Client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

func (r *RedisClient) ResetLoginAttempts(ctx context.Context, email string) error {
	key := "login_attempts:" + email
	return r.Client.Del(ctx, key).Err()
}
