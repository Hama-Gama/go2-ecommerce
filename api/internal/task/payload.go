package task

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeOrderConfirmation = "order:send_confirmation"
)

type OrderConfirmationPayload struct {
	OrderID     int64   `json:"order_id"`
	UserID      int64   `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
}

func NewOrderConfirmationTask(orderID, userID int64, totalAmount float64) (*asynq.Task, error) {
	payload, err := json.Marshal(OrderConfirmationPayload{
		OrderID:     orderID,
		UserID:      userID,
		TotalAmount: totalAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Создаем задачу с таймаутом и максимум 5 попытками повтора
	return asynq.NewTask(TypeOrderConfirmation, payload, asynq.MaxRetry(5)), nil
}
