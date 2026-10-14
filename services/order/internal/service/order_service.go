package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
	"github.com/nashinnov8/go-ecommerce/pkg/database"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/dto"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o domain.Order) (domain.Order, error)
	CreateOrderItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (domain.Order, error)
	ListOrdersByUserID(ctx context.Context, userID uuid.UUID, page, limit int32) ([]domain.Order, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error)
	ListUserOrders(ctx context.Context, userID uuid.UUID, page, limit int32) ([]dto.OrderResponse, error)
}

type orderService struct {
	repo      OrderRepository
	txManager database.TransactionManager
}

func NewOrderService(repo OrderRepository, txManager database.TransactionManager) OrderService {
	return &orderService{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	if len(req.Items) == 0 {
		return dto.OrderResponse{}, domain.EmptyOrderItemsError()
	}

	// 1. Calculate items subtotal and overall total amount
	var totalAmount int64
	itemsToCreate := make([]domain.OrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		productID, err := uuid.Parse(itemReq.ProductID)
		if err != nil {
			return dto.OrderResponse{}, apperror.BadRequest("INVALID_PRODUCT_ID", "Invalid product UUID", err)
		}

		subtotal := int64(itemReq.Quantity) * itemReq.UnitPrice
		totalAmount += subtotal

		itemsToCreate = append(itemsToCreate, domain.OrderItem{
			ProductID: productID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Subtotal:  subtotal,
		})
	}

	orderToCreate := domain.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      domain.OrderStatusPending,
	}

	var createdOrder domain.Order
	var createdItems []domain.OrderItem

	// 2. ATOMIC TRANSACTION: Both order and items are created inside the same tx
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		var err error

		// Step A: Insert Order
		createdOrder, err = s.repo.CreateOrder(txCtx, orderToCreate)
		if err != nil {
			return err // Trigger automatic rollback
		}

		// Step B: Insert each OrderItem with the generated OrderID
		for _, item := range itemsToCreate {
			item.OrderID = createdOrder.ID
			ci, err := s.repo.CreateOrderItem(txCtx, item)
			if err != nil {
				return err // Trigger automatic rollback
			}
			createdItems = append(createdItems, ci)
		}

		return nil // Trigger automatic commit
	})

	if err != nil {
		return dto.OrderResponse{}, apperror.Internal("Failed to create order transaction", err)
	}

	createdOrder.Items = createdItems
	return dto.ToOrderResponse(createdOrder), nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return dto.OrderResponse{}, domain.OrderNotFoundError(err)
		}
		return dto.OrderResponse{}, apperror.Internal("Failed to fetch order", err)
	}
	return dto.ToOrderResponse(order), nil
}

func (s *orderService) ListUserOrders(ctx context.Context, userID uuid.UUID, page, limit int32) ([]dto.OrderResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	orders, err := s.repo.ListOrdersByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, apperror.Internal("Failed to list orders", err)
	}

	res := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, dto.ToOrderResponse(o))
	}
	return res, nil
}
