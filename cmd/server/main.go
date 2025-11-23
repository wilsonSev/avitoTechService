package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/wilsonSev/avitoTechService/internal/api"
	"github.com/wilsonSev/avitoTechService/internal/services"
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

func main() {
	ctx := context.Background()

	// настройка пула соединений с БД
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	pool, err := storage.NewPool(ctx, dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer pool.Close()

	// Репозитории
	userRepo := storage.NewUserRepo(pool)
	teamRepo := storage.NewTeamRepo(pool)
	prRepo := storage.NewPRRepo(pool)

	// Сервисы
	userService := services.NewUserService(userRepo, prRepo)
	teamService := services.NewTeamService(teamRepo, userRepo)
	prService := services.NewPRService(prRepo, userRepo)

	// Хендлеры
	teamHandler := api.NewTeamHandler(teamService)
	userHandler := api.NewUserHandler(userService)
	prHandler := api.NewPRHandler(prService)

	router := api.NewRouter(teamHandler, userHandler, prHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	log.Println("Server started at :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
