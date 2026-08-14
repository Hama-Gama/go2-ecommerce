package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/yourusername/go-ecommerce-api/internal/domain"
	"github.com/yourusername/go-ecommerce-api/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetProducts godoc
// @Summary      Получить список товаров
// @Description  Возвращает список всех товаров или результаты поиска по запросу
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        query query string false "Поисковый запрос"
// @Success      200 {object} domain.ProductListResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	filter := domain.ProductFilter{
		Query: query.Get("query"),
	}

	// 1. Парсим Category ID
	if catStr := query.Get("category_id"); catStr != "" {
		if catID, err := strconv.ParseInt(catStr, 10, 32); err == nil {
			cID := int32(catID)
			filter.CategoryID = &cID
		}
	}

	// 2. Парсим Min Price
	if minStr := query.Get("min_price"); minStr != "" {
		if minP, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.MinPrice = &minP
		}
	}

	// 3. Парсим Max Price
	if maxStr := query.Get("max_price"); maxStr != "" {
		if maxP, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.MaxPrice = &maxP
		}
	}

	// 4. Парсим Limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// 5. Парсим Cursor (для курсорной пагинации)
	if cursorIDStr := query.Get("cursor_id"); cursorIDStr != "" {
		if cID, err := strconv.ParseInt(cursorIDStr, 10, 64); err == nil {
			filter.CursorID = &cID
		}
	}
	if cursorTimeStr := query.Get("cursor_time"); cursorTimeStr != "" {
		if cTime, err := time.Parse(time.RFC3339, cursorTimeStr); err == nil {
			filter.CursorTime = &cTime
		}
	}

	// Получаем данные из сервисного слоя
	response, err := h.service.GetProducts(ctx, filter)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch products"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
