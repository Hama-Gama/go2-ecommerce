package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go-ecommerce-api/internal/domain"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetList(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	products := make([]domain.Product, 0)

	var (
		conditions []string
		args       []interface{}
		argPos     = 1
	)

	// Поиск по названию
	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argPos))
		args = append(args, "%"+filter.Query+"%")
		argPos++
	}

	// Фильтрация по категории (если передана в filter)
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *filter.CategoryID)
		argPos++
	}

	query := `
		SELECT id, category_id, title, slug, description, price, stock, created_at, updated_at
		FROM products`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	// Добавляем пагинацию (если заданы Limit и Offset)
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
		argPos++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.Title,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}