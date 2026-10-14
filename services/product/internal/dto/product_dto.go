package dto

import (
	"time"

	"github.com/nashinnov8/go-ecommerce/services/product/internal/domain"
)

type CreateProductRequest struct {
	SKU         string `json:"sku" binding:"required,min=3,max=64"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int32  `json:"stock" binding:"gte=0"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	IsActive    *bool  `json:"is_active" binding:"required"`
}

type UpdateStockRequest struct {
	Delta int32 `json:"delta" binding:"required"`
}

type ListProductsQuery struct {
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
	Page     int32  `form:"page,default=1" binding:"gte=1"`
	Limit    int32  `form:"limit,default=20" binding:"gte=1,lte=100"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Stock       int32     `json:"stock"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID.String(),
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToProductListResponse(products []domain.Product) []ProductResponse {
	res := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		res = append(res, ToProductResponse(p))
	}
	return res
}
