package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisEngine struct {
	client *redis.Client
}

func newRedis(opts *redis.Options) ICache {
	c := redis.NewClient(opts)
	return &redisEngine{
		client: c,
	}
}

func (e *redisEngine) Ping(ctx context.Context) (string, error) {
	return e.client.Ping(ctx).Result()
}

func (e *redisEngine) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return e.client.Set(ctx, key, value, expiration).Err()
}

func (e *redisEngine) BindJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := e.client.Get(ctx, key).Bytes()
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return false, err
	}

	return true, nil
}
