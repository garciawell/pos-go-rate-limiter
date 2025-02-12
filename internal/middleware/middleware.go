package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garciawell/pos-go-rate-limiter/configs"
	"github.com/garciawell/pos-go-rate-limiter/internal/database"
	"github.com/go-redis/redis"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		rateLimiterToken, err := strconv.Atoi(configs.ENV.RateLimiterQtyToken)
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}

		duration, err := time.ParseDuration(configs.ENV.RateLimiterTimeToken)
		if err != nil {
			fmt.Println("Erro na conversão da data:", err)
			return
		}

		count, err := database.RedisClient.Get(apiKey).Int()
		if err != nil && err != redis.Nil {
			fmt.Println("Erro ao acessar Redis:", err)
			http.Error(w, "Erro ao acessar Redis", http.StatusInternalServerError)
			return
		}
		if count >= rateLimiterToken {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		countIp, err := database.RedisClient.Get(ip).Int()
		if err != nil && err != redis.Nil {
			fmt.Println("Erro ao acessar Redis:", err)
			http.Error(w, "Erro ao acessar Redis", http.StatusInternalServerError)
			return
		}

		rateLimiterIp, err := strconv.Atoi(configs.ENV.RateLimiterQtyIp)
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}

		durationIp, err := time.ParseDuration(configs.ENV.RateLimiterTimeIp)
		if err != nil {
			fmt.Println("Erro na conversão da data:", err)
			return
		}

		if countIp >= rateLimiterIp {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		_, err = database.RedisClient.Incr(apiKey).Result()
		if err != nil {
			http.Error(w, "Erro ao registrar requisição", http.StatusInternalServerError)
			return
		}

		_, err = database.RedisClient.Expire(apiKey, duration).Result()
		if err != nil {
			http.Error(w, "Erro ao definir tempo de expiração", http.StatusInternalServerError)
			return
		}

		_, err = database.RedisClient.Incr(ip).Result()
		if err != nil {
			http.Error(w, "Erro ao registrar requisição", http.StatusInternalServerError)
			return
		}

		_, err = database.RedisClient.Expire(ip, durationIp).Result()
		if err != nil {
			http.Error(w, "Erro ao definir tempo de expiração", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
