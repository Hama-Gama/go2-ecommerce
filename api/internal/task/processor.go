package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	ProcessTaskOrderConfirmation(ctx context.Context, t *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt) *RedisTaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // Обработка 10 задач параллельно
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("❌ Task failed: type=%s, err=%v", task.Type(), err)
			}),
		},
	)

	return &RedisTaskProcessor{server: server}
}

func (p *RedisTaskProcessor) ProcessTaskOrderConfirmation(ctx context.Context, t *asynq.Task) error {
	var payload OrderConfirmationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("📥 [Worker] Processing Order Confirmation for Order ID: %d, User ID: %d, Amount: $%.2f",
		payload.OrderID, payload.UserID, payload.TotalAmount)

	// Симуляция генерации PDF чека и отправки email
	// В реальном проекте тут вызов SMTP или Mailgun/SendGrid
	log.Printf("📄 [Worker] PDF Receipt generated for Order #%d", payload.OrderID)
	log.Printf("📧 [Worker] Email confirmation sent to User #%d", payload.UserID)

	return nil
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeOrderConfirmation, p.ProcessTaskOrderConfirmation)

	log.Println("⚙️  Asynq Worker Processor started...")
	return p.server.Run(mux)
}
