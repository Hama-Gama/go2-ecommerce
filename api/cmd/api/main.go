package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"

	"github.com/yourusername/go-ecommerce-api/internal/handler"
	"github.com/yourusername/go-ecommerce-api/internal/repository"
	"github.com/yourusername/go-ecommerce-api/internal/service"
	"github.com/yourusername/go-ecommerce-api/internal/task"
	_ "github.com/yourusername/go-ecommerce-api/docs"
	"github.com/swaggo/http-swagger"
)

// Helper-функция для чтения переменных окружения с фолбэком
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// @title           Go E-Commerce API
// @version         1.0
// @description     API сервис интернет-магазина на Go (Chi, PostgreSQL, Redis, Asynq).
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	ctx := context.Background()

	// 1. Подключение к PostgreSQL (динамически из Env, чтобы корректно работать в Docker)
	pgCfg := repository.PostgresConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgrespassword"),
		DBName:   getEnv("DB_NAME", "ecommerce_db"),
		SSLMode:  "disable",
	}

	pool, err := repository.NewPostgresPool(ctx, pgCfg)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	// 2. Опции подключения к Redis для Asynq (динамически из Env)
	redisOpt := asynq.RedisClientOpt{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
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

	// CORS middleware для соединения с Next.js
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/api/v1/products", productHnd.GetProducts)
	r.Post("/api/v1/orders", orderHnd.CreateOrder)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	fmt.Println("🚀 API Server & Asynq Worker are running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}