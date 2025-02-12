package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garciawell/pos-go-rate-limiter/configs"
	"github.com/garciawell/pos-go-rate-limiter/internal/database"
)

func RateLimiterMiddleware(redisClient *database.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			apiKey := r.Header.Get("API_KEY")

			rateLimiterToken, duration, err := GetRateLimitConfigs(configs.ENV.RateLimiterQtyToken, configs.ENV.RateLimiterTimeToken)
			if err != nil {
				fmt.Println("Erro ao obter configuração de limite por API Key", err)
				http.Error(w, "Erro interno", http.StatusInternalServerError)
				return
			}

			rateLimiterIp, durationIp, err := GetRateLimitConfigs(configs.ENV.RateLimiterQtyIp, configs.ENV.RateLimiterTimeIp)
			if err != nil {
				fmt.Println("Erro ao obter configuração de limite por IP", err)
				http.Error(w, "Erro interno", http.StatusInternalServerError)
				return
			}

			if apiKey != "" {
				if ExceededLimit(redisClient, ip, rateLimiterToken) {
					http.Error(w, "Too many requests", http.StatusTooManyRequests)
					return
				}
			} else {
				fmt.Println("API Key não informada", rateLimiterIp)
				if ExceededLimit(redisClient, apiKey, rateLimiterIp) || ExceededLimit(redisClient, ip, rateLimiterIp) {
					http.Error(w, "Too many requests", http.StatusTooManyRequests)
					return
				}
			}

			IncrementRequest(redisClient, apiKey, duration)
			IncrementRequest(redisClient, ip, durationIp)

			next.ServeHTTP(w, r)
		})
	}
}
