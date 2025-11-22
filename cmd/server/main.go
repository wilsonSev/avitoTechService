package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/wilsonSev/avitoTechService/internal/storage"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/wilsonSev/avitoTechService/internal/api"
)

func main() {
	// настройка пула соединений с БД
	_ = godotenv.Load()
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	dsn := "postgres://" + user + ":" + password + "@localhost:5432/app?sslmode=disable"

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := storage.NewPool(ctx, dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer pool.Close()

	// настройка сервера
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	h := api.NewHandlers()

	r.Post("/team/add", h.CreateTeam)

	http.ListenAndServe(":3000", r)
}
