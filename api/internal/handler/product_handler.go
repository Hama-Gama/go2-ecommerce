package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	if catStr := query.Get("category_id"); catStr != "" {
		if catID, err := strconv.ParseInt(catStr, 10, 32); err == nil {
			cID := int32(catID)
			filter.CategoryID = &cID
		}
	}

	if minStr := query.Get("min_price"); minStr != "" {
		if minP, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.MinPrice = &minP
		}
	}

	if maxStr := query.Get("max_price"); maxStr != "" {
		if maxP, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.MaxPrice = &maxP
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

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

	response, err := h.service.GetProducts(ctx, filter)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch products"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"invalid multipart form"}`, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.ParseInt(r.FormValue("stock"), 10, 32)
	catID, err := strconv.ParseInt(r.FormValue("category_id"), 10, 32)
	if err != nil || catID == 0 {
		catID = 1
	}

	if title == "" || price < 0 {
		http.Error(w, `{"error":"title and valid price are required"}`, http.StatusBadRequest)
		return
	}

	var imageURL string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		os.MkdirAll("./uploads", os.ModePerm)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
		dstPath := filepath.Join("./uploads", filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, `{"error":"failed to save image"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, `{"error":"failed to write image"}`, http.StatusInternalServerError)
			return
		}
		imageURL = "/uploads/" + filename
	}

	slug := fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(title, " ", "-")), time.Now().Unix())

	product := &domain.Product{
		CategoryID:  int32(catID),
		Title:       title,
		Slug:        slug,
		Description: description,
		Price:       price,
		Stock:       int32(stock),
		ImageURL:    imageURL,
	}

	if err := h.service.CreateProduct(r.Context(), product); err != nil {
		http.Error(w, `{"error":"failed to create product"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
