package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yourusername/go-ecommerce-api/internal/domain"
	domainMocks "github.com/yourusername/go-ecommerce-api/internal/domain/mocks"
	"github.com/yourusername/go-ecommerce-api/internal/service"
	taskMocks "github.com/yourusername/go-ecommerce-api/internal/task/mocks"
)

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockOrderRepo := new(domainMocks.OrderRepository)
	mockTaskDistributor := new(taskMocks.TaskDistributor)

	orderService := service.NewOrderService(mockOrderRepo, mockTaskDistributor)
	ctx := context.Background()

	input := domain.CreateOrderRequest{
		UserID: 1,
		Items: []domain.CreateOrderItemRequest{
			{ProductID: 10, Quantity: 2},
		},
	}

	createdOrder := &domain.Order{
		ID:          100,
		UserID:      1,
		Status:      domain.OrderStatusPending,
		TotalAmount: 200.0,
	}

	// 1. Ожидаем вызов сохранения заказа и списания со склада в репозитории
	mockOrderRepo.
		On("CreateOrderWithInventory", ctx, input).
		Return(createdOrder, nil)

	// 2. Ожидаем отправку асинхронной задачи подтверждения заказа в Asynq
	mockTaskDistributor.
		On("DistributeTaskOrderConfirmation", ctx, mock.Anything, mock.Anything).
		Return(nil)

	// Выполняем тестирование
	result, err := orderService.CreateOrder(ctx, input)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.ID)
	assert.Equal(t, 200.0, result.TotalAmount)

	mockOrderRepo.AssertExpectations(t)
	mockTaskDistributor.AssertExpectations(t)
}
