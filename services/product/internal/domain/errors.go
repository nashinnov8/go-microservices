package domain

import (
	"errors"

	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrDuplicateSKU      = errors.New("product SKU already exists")
	ErrInsufficientStock = errors.New("insufficient product stock")
)

func ProductNotFoundError(err error) *apperror.AppError {
	return apperror.NotFound("PRODUCT_NOT_FOUND", "Product not found", err)
}

func DuplicateSKUError(sku string, err error) *apperror.AppError {
	return apperror.Conflict("PRODUCT_DUPLICATE_SKU", "Product with SKU '"+sku+"' already exists", err)
}

func InsufficientStockError(err error) *apperror.AppError {
	return apperror.Conflict("INSUFFICIENT_STOCK", "Insufficient stock for the requested operation", err)
}
