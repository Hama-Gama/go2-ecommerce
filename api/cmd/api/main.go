package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hibiken/asynq"

	"github.com/yourusername/go-ecommerce-api/internal/handler"
	"github.com/yourusername/go-ecommerce-api/internal/repository"
	"github.com/yourusername/go-ecommerce-api/internal/service"
	"github.com/yourusername/go-ecommerce-api/internal/task"
)

func main() {
	ctx := context.Background()

	// 1. Подключение к PostgreSQL
	pgCfg := repository.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgrespassword",
		DBName:   "ecommerce_db",
		SSLMode:  "disable",
	}

	pool, err := repository.NewPostgresPool(ctx, pgCfg)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	// 2. Опции подключения к Redis для Asynq
	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	// 3. Инициализация Asynq Distributor (Client) & Processor (Worker)
	taskDistributor := task.NewRedisTaskDistributor(redisOpt)
	taskProcessor := task.NewRedisTaskProcessor(redisOpt)

	// Запускаем Asynq Worker в отдельной горутине
	go func() {
		if err := taskProcessor.Start(); err != nil {
			log.Fatalf("Failed to start task processor: %v", err)
		}
	}()

	// 4. Инициализация слоев Products
	productRepo := repository.NewProductRepository(pool)
	productSvc := service.NewProductService(productRepo)
	productHnd := handler.NewProductHandler(productSvc)

	// 5. Инициализация слоев Orders (передаем taskDistributor)
	orderRepo := repository.NewOrderRepository(pool)
	orderSvc := service.NewOrderService(orderRepo, taskDistributor)
	orderHnd := handler.NewOrderHandler(orderSvc)

	// 6. Настройка роутера Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/api/v1/products", productHnd.GetProducts)
	r.Post("/api/v1/orders", orderHnd.CreateOrder)

	fmt.Println("🚀 API Server & Asynq Worker are running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
