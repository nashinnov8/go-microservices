package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/repository/sqlc"
)

const pgUniqueViolationCode = "23505"

type ProductRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	row, err := r.queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		IsActive:    p.IsActive,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return domain.Product{}, domain.ErrDuplicateSKU
		}
		return domain.Product{}, err
	}

	return toDomain(row), nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	row, err := r.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		return domain.Product{}, err
	}
	return toDomain(row), nil
}

func (r *ProductRepository) List(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, int64, error) {
	offset := (filter.Page - 1) * filter.Limit

	total, err := r.queries.CountProducts(ctx, sqlc.CountProductsParams{
		Column1: filter.IsActive,
		Column2: filter.Search,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.queries.ListProducts(ctx, sqlc.ListProductsParams{
		Column1: filter.IsActive,
		Column2: filter.Search,
		Limit:   filter.Limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	products := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, toDomain(row))
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	row, err := r.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		IsActive:    p.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		return domain.Product{}, err
	}
	return toDomain(row), nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, delta int32) (domain.Product, error) {
	row, err := r.queries.UpdateProductStock(ctx, sqlc.UpdateProductStockParams{
		ID:    id,
		Stock: delta,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Distinguish between non-existent product and insufficient stock
			if _, getErr := r.queries.GetProductByID(ctx, id); errors.Is(getErr, pgx.ErrNoRows) {
				return domain.Product{}, domain.ErrProductNotFound
			}
			return domain.Product{}, domain.ErrInsufficientStock
		}
		return domain.Product{}, err
	}
	return toDomain(row), nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.queries.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func toDomain(row sqlc.Product) domain.Product {
	return domain.Product{
		ID:          row.ID,
		SKU:         row.Sku,
		Name:        row.Name,
		Description: row.Description,
		Price:       row.Price,
		Stock:       row.Stock,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
