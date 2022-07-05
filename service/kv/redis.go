package kv

import (
	"context"
	"github.com/OhanaFS/ohana/config"
	"github.com/go-redis/redis/v9"
	"time"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedis(cfg *config.Config) KV {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	return &Redis{rdb}
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, key, value string, ttl time.Duration) error {

	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {

	return r.rdb.Expire(ctx, key, 0).Err()
}
