package config

import (
	"context"
	"github.com/go-redis/redis/v9"
)

func NewRedis(config *Config) (*redis.Client, context.Context) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	return rdb, ctx
}
