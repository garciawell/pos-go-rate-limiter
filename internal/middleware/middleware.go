package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/garciawell/pos-go-rate-limiter/configs"
)

var limiterData = make(map[string]int)

func RateLimitMiddleware(next http.Handler, ENVS *configs.Conf) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		apiKey := r.Header.Get("API_KEY")

		rateLimiterQty, err := strconv.Atoi(ENVS.RateLimiterQtyIp)
		if err != nil {
			fmt.Println("Erro na conversão:", err)
			return
		}

		if limiterData[apiKey] > rateLimiterQty {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
			return
		}

		limiterData[apiKey]++

		fmt.Println("IP", ip)
		fmt.Println("Contador2", limiterData[apiKey])
		next.ServeHTTP(w, r)
	})
}
