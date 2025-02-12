package main

import (
	"fmt"
	"net/http"

	"github.com/garciawell/pos-go-rate-limiter/configs"
	"github.com/garciawell/pos-go-rate-limiter/internal/database"
	middlewareInternal "github.com/garciawell/pos-go-rate-limiter/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func init() {
	confEnv, err := configs.LoadConfig(".")
	configs.ENV = confEnv
	if err != nil {
		panic(err)
	}
}

func main() {
	conn := database.NewRedisClient()
	repo := database.NewRepo(conn)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middlewareInternal.RateLimiterMiddleware(repo))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	defer conn.Close()

	fmt.Println("Server running on port 8080!")
	http.ListenAndServe(":8080", r)
}
