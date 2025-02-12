package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garciawell/pos-go-rate-limiter/internal/database"
	"github.com/go-redis/redis"
)

func IncrementRequest(redisClient database.RepoInterface, key string, duration time.Duration) {
	if key == "" {
		return
	}
	_, err := redisClient.IncrBase(key)
	if err != nil {
		fmt.Printf("Erro ao incrementar contador para %s\n", key)
	}

	_, err = redisClient.SetExBase(key, duration)
	if err != nil {
		fmt.Printf("Erro ao definir expiração para %s\n", key)
	}
}

func ExceededLimit(redisClient database.RepoInterface, key string, limit int) bool {
	count, err := redisClient.GetBase(key)
	if err != nil && err != redis.Nil {
		fmt.Printf("Erro ao acessar Redis para %s\n", key)
	}
	fmt.Println("Contador para", count)
	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Printf("Erro ao converter contador para inteiro para %s\n", key)
		return false
	}
	return countInt >= limit
}

func GetRateLimitConfigs(qty, timeStr string) (int, time.Duration, error) {
	limit, err := strconv.Atoi(qty)
	if err != nil {
		return 0, 0, fmt.Errorf("falha ao converter limite: %w", err)
	}

	duration, err := time.ParseDuration(timeStr)
	if err != nil {
		return 0, 0, fmt.Errorf("falha ao converter duração: %w", err)
	}

	return limit, duration, nil
}
