package provider

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisProvider struct {
	redisClient *redis.Client
}

func NewRedisProvider() (*redisProvider, error) {
	client := redis.NewClient(&redis.Options{
		Password: "password",
		Addr:     fmt.Sprintf("%s:%d", "localhost", 6379),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &redisProvider{
		redisClient: client,
	}, nil
}
func (w *redisProvider) Write(ctx context.Context, ttl time.Duration, key string, values ...interface{}) error {
	if err := w.redisClient.SAdd(ctx, key, values).Err(); err != nil {
		return fmt.Errorf("add to redis  %s: %v", key, err)
	}

	currentExpire, err := w.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error get ttl in redis for %s: %v", key, err)
	}

	if currentExpire == time.Duration(-1) {
		if err := w.redisClient.Expire(ctx, key, ttl).Err(); err != nil {
			return fmt.Errorf("expire redis %s: %v", key, err)
		}
	}

	return nil
}
