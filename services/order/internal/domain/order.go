package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
)

type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TotalAmount int64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem
}

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	UnitPrice int64
	Subtotal  int64
	CreatedAt time.Time
}
