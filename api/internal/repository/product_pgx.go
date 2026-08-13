package repository

import (
	"context"
	"fmt"

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

	var query string
	var args []interface{}

	if filter.Query != "" {
		query = `
			SELECT id, category_id, title, slug, description, price, 10 as stock, created_at, updated_at
			FROM products
			WHERE (tsv @@ plainto_tsquery('simple', $1) OR title ILIKE '%' || $1 || '%')
			ORDER BY created_at DESC`
		args = append(args, filter.Query)
	} else {
		query = `
			SELECT id, category_id, title, slug, description, price, 10 as stock, created_at, updated_at
			FROM products
			ORDER BY created_at DESC`
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return products, fmt.Errorf("failed to query products: %w", err)
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
			return products, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}