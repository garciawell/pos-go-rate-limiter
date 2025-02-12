package database

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return client
}
