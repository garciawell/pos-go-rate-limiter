package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garciawell/pos-go-rate-limiter/configs"
)

type LimiterInfo struct {
	Count          int
	BlockTimestamp time.Time
}

var limiterData = make(map[string]LimiterInfo)

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

		if limiterData[apiKey].Count >= rateLimiterToken {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		if limiterData[ip].Count >= rateLimiterIp {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		limiterInfo := limiterData[apiKey]
		limiterInfo.Count++
		limiterInfo.BlockTimestamp = time.Now()
		limiterData[apiKey] = limiterInfo

		fmt.Println("IP", ip)
		fmt.Println("Contador", limiterData[apiKey].Count)
		fmt.Println("Timestamp", limiterData[apiKey].BlockTimestamp)
		next.ServeHTTP(w, r)
	})
}
