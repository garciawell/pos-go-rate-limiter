package database

import (
	"github.com/garciawell/pos-go-rate-limiter/configs"
	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: configs.ENV.DbHost + ":" + "6379",
	})

	return client
}
