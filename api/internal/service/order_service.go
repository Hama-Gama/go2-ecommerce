package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/yourusername/go-ecommerce-api/internal/domain"
	"github.com/yourusername/go-ecommerce-api/internal/task"
)

type OrderService struct {
	repo            domain.OrderRepository
	taskDistributor task.TaskDistributor
}

func NewOrderService(repo domain.OrderRepository, taskDistributor task.TaskDistributor) *OrderService {
	return &OrderService{
		repo:            repo,
		taskDistributor: taskDistributor,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	if req.UserID <= 0 {
		return nil, errors.New("invalid user_id")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("invalid product_id: %d", item.ProductID)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity for product %d must be greater than zero", item.ProductID)
		}
	}

	order, err := s.repo.CreateOrderWithInventory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("service failed to create order: %w", err)
	}

	// Отправляем асинхронную задачу в Redis очередь
	err = s.taskDistributor.DistributeTaskOrderConfirmation(ctx, task.OrderConfirmationPayload{
		OrderID:     order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
	})
	if err != nil {
		log.Printf("⚠️  Failed to enqueue order confirmation task: %v", err)
		// Не фейлим весь HTTP запрос, если упала только очередь
	}

	return order, nil
}
