package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/dto"
)

// ProductRepository defines the persistence contract required by the product service.
type ProductRepository interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	List(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, int64, error)
	Update(ctx context.Context, p domain.Product) (domain.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, delta int32) (domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProductService defines the business logic operations for products.
type ProductService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error)
	ListProducts(ctx context.Context, query dto.ListProductsQuery) ([]dto.ProductResponse, int64, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	UpdateStock(ctx context.Context, id uuid.UUID, req dto.UpdateStockRequest) (dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := domain.Product{
		SKU:         strings.ToUpper(strings.TrimSpace(req.SKU)),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    isActive,
	}

	created, err := s.repo.Create(ctx, product)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateSKU) {
			return dto.ProductResponse{}, domain.DuplicateSKUError(product.SKU, err)
		}
		return dto.ProductResponse{}, apperror.Internal("Failed to create product", err)
	}

	return dto.ToProductResponse(created), nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return dto.ProductResponse{}, domain.ProductNotFoundError(err)
		}
		return dto.ProductResponse{}, apperror.Internal("Failed to fetch product", err)
	}

	return dto.ToProductResponse(product), nil
}

func (s *productService) ListProducts(ctx context.Context, query dto.ListProductsQuery) ([]dto.ProductResponse, int64, error) {
	filter := domain.ListProductsFilter{
		Search:   strings.TrimSpace(query.Search),
		IsActive: query.IsActive,
		Page:     query.Page,
		Limit:    query.Limit,
	}

	products, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal("Failed to list products", err)
	}

	return dto.ToProductListResponse(products), total, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product := domain.Product{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		IsActive:    *req.IsActive,
	}

	updated, err := s.repo.Update(ctx, product)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return dto.ProductResponse{}, domain.ProductNotFoundError(err)
		}
		return dto.ProductResponse{}, apperror.Internal("Failed to update product", err)
	}

	return dto.ToProductResponse(updated), nil
}

func (s *productService) UpdateStock(ctx context.Context, id uuid.UUID, req dto.UpdateStockRequest) (dto.ProductResponse, error) {
	updated, err := s.repo.UpdateStock(ctx, id, req.Delta)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			return dto.ProductResponse{}, domain.ProductNotFoundError(err)
		case errors.Is(err, domain.ErrInsufficientStock):
			return dto.ProductResponse{}, domain.InsufficientStockError(err)
		default:
			return dto.ProductResponse{}, apperror.Internal("Failed to update product stock", err)
		}
	}

	return dto.ToProductResponse(updated), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ProductNotFoundError(err)
		}
		return apperror.Internal("Failed to delete product", err)
	}
	return nil
}
