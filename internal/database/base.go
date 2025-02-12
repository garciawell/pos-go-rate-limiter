package database

import (
	"time"

	"github.com/go-redis/redis"
)

type Repo struct {
	RedisClient *redis.Client
}

type RepoInterface interface {
	GetBase(key string) (string, error)
	IncrBase(key string) (int64, error)
	SetBase(key string, value interface{}) (string, error)
	SetExBase(key string, expiration time.Duration) (bool, error)
}

func NewRepo(db *redis.Client) *Repo {
	return &Repo{
		RedisClient: db,
	}
}

func (r *Repo) GetBase(key string) (string, error) {
	return r.RedisClient.Get(key).Result()
}

func (r *Repo) IncrBase(key string) (int64, error) {
	return r.RedisClient.Incr(key).Result()
}

func (r *Repo) SetBase(key string, value interface{}) (string, error) {
	return r.RedisClient.Set(key, value, 0).Result()
}

func (r *Repo) SetExBase(key string, expiration time.Duration) (bool, error) {
	return r.RedisClient.Expire(key, expiration).Result()
}
