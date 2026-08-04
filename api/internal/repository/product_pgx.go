package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go-ecommerce-api/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) GetList(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	fetchLimit := filter.Limit + 1

	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argID := 1

	queryBuilder.WriteString(`
		SELECT 
			p.id, p.category_id, p.title, p.slug, p.description, p.price, 
			COALESCE(i.stock, 0) as stock, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN inventories i ON p.id = i.product_id
		WHERE 1=1
	`)

	if filter.Query != "" {
		fmt.Fprintf(&queryBuilder, " AND p.tsv @@ websearch_to_tsquery('english', $%d)", argID)
		args = append(args, filter.Query)
		argID++
	}

	if filter.CategoryID != nil {
		fmt.Fprintf(&queryBuilder, " AND p.category_id = $%d", argID)
		args = append(args, *filter.CategoryID)
		argID++
	}

	if filter.MinPrice != nil {
		fmt.Fprintf(&queryBuilder, " AND p.price >= $%d", argID)
		args = append(args, *filter.MinPrice)
		argID++
	}
	if filter.MaxPrice != nil {
		fmt.Fprintf(&queryBuilder, " AND p.price <= $%d", argID)
		args = append(args, *filter.MaxPrice)
		argID++
	}

	if filter.CursorTime != nil && filter.CursorID != nil {
		fmt.Fprintf(&queryBuilder, " AND (p.created_at, p.id) < ($%d, $%d)", argID, argID+1)
		args = append(args, *filter.CursorTime, *filter.CursorID)
		argID += 2
	}

	fmt.Fprintf(&queryBuilder, " ORDER BY p.created_at DESC, p.id DESC LIMIT $%d", argID)
	args = append(args, fetchLimit)

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get products query: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.CategoryID, &p.Title, &p.Slug, &p.Description, 
			&p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}
