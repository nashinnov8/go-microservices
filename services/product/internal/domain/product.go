package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product represents the core domain entity for a catalog item.
type Product struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       int64
	Stock       int32
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListProductsFilter struct {
	Search   string
	IsActive *bool
	Page     int32
	Limit    int32
}
