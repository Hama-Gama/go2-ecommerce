package domain

import (
	"context"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type Order struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderRequest struct {
	UserID int64                    `json:"user_id"` // Пока хардкодим/берём из параметров, пока нет JWT Auth
	Items  []CreateOrderItemRequest `json:"items"`
}

type OrderRepository interface {
	CreateOrderWithInventory(ctx context.Context, req CreateOrderRequest) (*Order, error)
}
