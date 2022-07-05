package config

import (
	"context"
	"github.com/go-redis/redis/v9"
)

func NewRedis(config *Config) (*redis.Client, context.Context) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.address,
		Password: config.Redis.password,
		DB:       config.Redis.db,
	})

	return rdb, ctx
}
