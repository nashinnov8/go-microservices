package dto

import (
	"time"

	"github.com/nashinnov8/go-ecommerce/services/order/internal/domain"
)

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
	UnitPrice int64  `json:"unit_price" binding:"required,gt=0"`
}

type OrderResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	TotalAmount int64               `json:"total_amount"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Items       []OrderItemResponse `json:"items,omitempty"`
}

type OrderItemResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
	Subtotal  int64     `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
}

func ToOrderResponse(order domain.Order) OrderResponse {
	itemsResp := make([]OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		itemsResp = append(itemsResp, OrderItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
			CreatedAt: item.CreatedAt,
		})
	}

	return OrderResponse{
		ID:          order.ID.String(),
		UserID:      order.UserID.String(),
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Items:       itemsResp,
	}
}
