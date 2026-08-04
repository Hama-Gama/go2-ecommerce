package service

import (
	"context"
	"fmt"

	"github.com/yourusername/go-ecommerce-api/internal/domain"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context, filter domain.ProductFilter) (*domain.ProductListResponse, error) {
	// Значение по умолчанию для лимита
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}

	// Запрашиваем из репозитория
	products, err := s.repo.GetList(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service failed to get products: %w", err)
	}

	hasMore := false
	var nextCursor *domain.Cursor

	// Проверяем, получили ли мы дополнительный элемент (HasMore)
	if len(products) > filter.Limit {
		hasMore = true
		// Отрезаем лишний элемент
		products = products[:filter.Limit]
		
		// Берём последний элемент из оставшихся для формирования курсора
		lastProduct := products[len(products)-1]
		nextCursor = &domain.Cursor{
			ID:        lastProduct.ID,
			CreatedAt: lastProduct.CreatedAt,
		}
	}

	// Если массив пустой, возвращаем пустой срез вместо nil для красивого JSON []
	if products == nil {
		products = make([]domain.Product, 0)
	}

	return &domain.ProductListResponse{
		Products:   products,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
