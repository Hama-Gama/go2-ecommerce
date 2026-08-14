package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yourusername/go-ecommerce-api/internal/domain"
	"github.com/yourusername/go-ecommerce-api/internal/repository"
	"github.com/yourusername/go-ecommerce-api/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder godoc
// @Summary      Создать новый заказ
// @Description  Принимает товары и сбрасывает фоновую задачу отправки уведомления в Asynq
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        payload body domain.CreateOrderRequest true "Данные заказа"
// @Success      201 {object} domain.Order
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req)
	if err != nil {
		log.Printf("❌ CreateOrder Error: %v", err) // Печатаем подробную ошибку в консоль сервера

		w.Header().Set("Content-Type", "application/json")
		
		if errors.Is(err, repository.ErrInsufficientStock) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if errors.Is(err, repository.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
