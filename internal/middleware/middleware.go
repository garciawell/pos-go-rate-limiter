package middleware

import (
	"fmt"
	"net/http"
)

var limiterData = make(map[string]int)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter, ok := limiterData[ip]
		if !ok {
			limiter = 0
		}
		limiter++
		fmt.Println("IP", ip)
		fmt.Println("Contador2", limiterData[ip])
		next.ServeHTTP(w, r)
	})
}
