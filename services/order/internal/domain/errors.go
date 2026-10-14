package domain

import (
	"errors"

	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrEmptyOrderItems  = errors.New("order must contain at least one item")
	ErrInvalidItemPrice = errors.New("item price and quantity must be greater than zero")
)

func OrderNotFoundError(err error) *apperror.AppError {
	return apperror.NotFound("ORDER_NOT_FOUND", "Order not found", err)
}

func EmptyOrderItemsError() *apperror.AppError {
	return apperror.BadRequest("EMPTY_ORDER_ITEMS", "Order must contain at least one item", ErrEmptyOrderItems)
}
