package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/dto"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/service"
	httptransport "github.com/nashinnov8/go-ecommerce/services/product/internal/transport/http"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/transport/http/handler"
)

type mockProductRepo struct {
	products map[uuid.UUID]domain.Product
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[uuid.UUID]domain.Product)}
}

func (m *mockProductRepo) Create(_ context.Context, p domain.Product) (domain.Product, error) {
	for _, existing := range m.products {
		if existing.SKU == p.SKU {
			return domain.Product{}, domain.ErrDuplicateSKU
		}
	}
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	m.products[p.ID] = p
	return p, nil
}

func (m *mockProductRepo) GetByID(_ context.Context, id uuid.UUID) (domain.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	return p, nil
}

func (m *mockProductRepo) List(_ context.Context, _ domain.ListProductsFilter) ([]domain.Product, int64, error) {
	res := make([]domain.Product, 0, len(m.products))
	for _, p := range m.products {
		res = append(res, p)
	}
	return res, int64(len(res)), nil
}

func (m *mockProductRepo) Update(_ context.Context, p domain.Product) (domain.Product, error) {
	existing, ok := m.products[p.ID]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	existing.Name = p.Name
	existing.Description = p.Description
	existing.Price = p.Price
	existing.IsActive = p.IsActive
	existing.UpdatedAt = time.Now()
	m.products[p.ID] = existing
	return existing, nil
}

func (m *mockProductRepo) UpdateStock(_ context.Context, id uuid.UUID, delta int32) (domain.Product, error) {
	existing, ok := m.products[id]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	if existing.Stock+delta < 0 {
		return domain.Product{}, domain.ErrInsufficientStock
	}
	existing.Stock += delta
	existing.UpdatedAt = time.Now()
	m.products[id] = existing
	return existing, nil
}

func (m *mockProductRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := m.products[id]; !ok {
		return domain.ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func TestCreateAndGetProduct(t *testing.T) {
	repo := newMockProductRepo()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := httptransport.NewRouter(log, "test", h)

	body, _ := json.Marshal(dto.CreateProductRequest{
		SKU:         "ip16-pro-256",
		Name:        "iPhone 16 Pro 256GB",
		Description: "Apple flagship smartphone",
		Price:       28990000,
		Stock:       50,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var createResp struct {
		Success bool                `json:"success"`
		Data    dto.ProductResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !createResp.Success || createResp.Data.SKU != "IP16-PRO-256" {
		t.Fatalf("unexpected response payload: %+v", createResp)
	}

	// Fetch the created product by ID
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+createResp.Data.ID, nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d: %s", getRec.Code, getRec.Body.String())
	}
}
