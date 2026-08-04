package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go-ecommerce-api/internal/domain"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

var (
	ErrInsufficientStock = errors.New("insufficient stock for product")
	ErrProductNotFound   = errors.New("product not found")
)

func (r *OrderRepository) CreateOrderWithInventory(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	items := make([]domain.CreateOrderItemRequest, len(req.Items))
	copy(items, req.Items)
	sort.Slice(items, func(i, j int) bool {
		return items[i].ProductID < items[j].ProductID
	})

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var totalAmount float64
	orderItems := make([]domain.OrderItem, 0, len(items))

	for _, item := range items {
		var currentStock int32
		var price float64

		lockQuery := `SELECT stock FROM inventories WHERE product_id = $1 FOR UPDATE`
		err := tx.QueryRow(ctx, lockQuery, item.ProductID).Scan(&currentStock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%w: product_id %d", ErrProductNotFound, item.ProductID)
			}
			return nil, fmt.Errorf("failed to lock inventory for product %d: %w", item.ProductID, err)
		}

		priceQuery := `SELECT price FROM products WHERE id = $1`
		err = tx.QueryRow(ctx, priceQuery, item.ProductID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch price for product %d: %w", item.ProductID, err)
		}

		if currentStock < item.Quantity {
			return nil, fmt.Errorf("%w: product_id %d (available: %d, requested: %d)", 
				ErrInsufficientStock, item.ProductID, currentStock, item.Quantity)
		}

		updateStockQuery := `UPDATE inventories SET stock = stock - $1, updated_at = NOW() WHERE product_id = $2`
		_, err = tx.Exec(ctx, updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct stock for product %d: %w", item.ProductID, err)
		}

		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: price,
		})
	}

	var orderID int64
	var createdAt, updatedAt time.Time

	createOrderQuery := `
		INSERT INTO orders (user_id, status, total_amount) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, createOrderQuery, req.UserID, domain.OrderStatusPending, totalAmount).Scan(&orderID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	for i := range orderItems {
		orderItems[i].OrderID = orderID
		insertItemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id
		`
		err = tx.QueryRow(ctx, insertItemQuery, orderID, orderItems[i].ProductID, orderItems[i].Quantity, orderItems[i].UnitPrice).Scan(&orderItems[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &domain.Order{
		ID:          orderID,
		UserID:      req.UserID,
		Status:      domain.OrderStatusPending,
		TotalAmount: totalAmount,
		Items:       orderItems,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
