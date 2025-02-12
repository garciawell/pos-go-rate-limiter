package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garciawell/pos-go-rate-limiter/configs"
	"github.com/garciawell/pos-go-rate-limiter/types"
	"github.com/garciawell/pos-go-rate-limiter/utils"
)

var limiterData = make(map[string]types.LimiterInfo)

func RateLimitMiddleware(next http.Handler, ENVS *configs.Conf) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		rateLimiterToken, err := strconv.Atoi(ENVS.RateLimiterQtyToken)
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}
		rateLimiterIp, err := strconv.Atoi(ENVS.RateLimiterQtyToken)
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}

		duration, err := time.ParseDuration(ENVS.RateLimiterTimeToken)
		if err != nil {
			fmt.Println("Erro na conversão da data:", err)
			return
		}

		if limiterData[apiKey].Count >= rateLimiterToken {
			utils.ValidateMiddleware(limiterData[apiKey], duration)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		if limiterData[ip].Count >= rateLimiterIp {
			utils.ValidateMiddleware(limiterData[ip], duration)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		limiterInfo := limiterData[apiKey]
		limiterInfo.Count++
		limiterData[apiKey] = limiterInfo

		fmt.Println("IP", ip)
		fmt.Println("Contador", limiterData[apiKey].Count)
		fmt.Println("Timestamp", limiterData[apiKey].BlockTimestamp)
		next.ServeHTTP(w, r)
	})
}
