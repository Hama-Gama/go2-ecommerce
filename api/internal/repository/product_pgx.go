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

	// Фильтрация по категории
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *filter.CategoryID)
		argPos++
	}


		query := `
	SELECT p.id, p.category_id, p.title, p.slug, p.description, p.price, COALESCE(i.stock, 0), p.image_url, p.created_at, p.updated_at
	FROM products p
	LEFT JOIN inventories i ON p.id = i.product_id`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	// Ограничение количества
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
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


func (r *ProductRepository) Create(ctx context.Context, p *Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	insertProduct := `
		INSERT INTO products (category_id, title, slug, description, price, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, insertProduct,
		p.CategoryID, p.Title, p.Slug, p.Description, p.Price, p.ImageURL,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	insertInventory := `
		INSERT INTO inventories (product_id, stock)
		VALUES ($1, $2)`

	if _, err := tx.Exec(ctx, insertInventory, p.ID, p.Stock); err != nil {
		return fmt.Errorf("failed to insert inventory: %w", err)
	}

	return tx.Commit(ctx)
}