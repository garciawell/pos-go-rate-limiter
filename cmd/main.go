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

var ENVS *configs.Conf

func init() {
	configs, err := configs.LoadConfig(".")
	ENVS = configs
	if err != nil {
		panic(err)
	}

	conn := database.NewRedisClient()
	if err != nil {
		panic(err)
	}

	conn.Close()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middlewareInternal.RateLimitMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	fmt.Println("Server running on port 8080!")
	http.ListenAndServe(":8080", r)
}
