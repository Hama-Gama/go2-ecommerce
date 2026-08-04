package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	CategoryID  int32     `json:"category_id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductFilter struct {
	Query      string     `json:"query"`
	CategoryID *int32     `json:"category_id"`
	MinPrice   *float64   `json:"min_price"`
	MaxPrice   *float64   `json:"max_price"`
	CursorID   *int64     `json:"cursor_id"`
	CursorTime *time.Time `json:"cursor_time"`
	Limit      int        `json:"limit"`
}

type ProductListResponse struct {
	Products   []Product `json:"products"`
	NextCursor *Cursor   `json:"next_cursor,omitempty"`
	HasMore    bool      `json:"has_more"`
}

type Cursor struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductRepository interface {
	GetList(ctx context.Context, filter ProductFilter) ([]Product, error)
}
