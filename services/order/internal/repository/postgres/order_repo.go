package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nashinnov8/go-ecommerce/pkg/database"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/repository/sqlc"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// getDB returns the active transaction from context if one exists; otherwise returns the connection pool.
func (r *OrderRepository) getDB(ctx context.Context) sqlc.DBTX {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *OrderRepository) CreateOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	q := sqlc.New(r.getDB(ctx))
	row, err := q.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Status:      string(o.Status),
	})
	if err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		ID:          row.ID,
		UserID:      row.UserID,
		TotalAmount: row.TotalAmount,
		Status:      domain.OrderStatus(row.Status),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *OrderRepository) CreateOrderItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	q := sqlc.New(r.getDB(ctx))
	row, err := q.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
		OrderID:   item.OrderID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Subtotal:  item.Subtotal,
	})
	if err != nil {
		return domain.OrderItem{}, err
	}

	return domain.OrderItem{
		ID:        row.ID,
		OrderID:   row.OrderID,
		ProductID: row.ProductID,
		Quantity:  row.Quantity,
		UnitPrice: row.UnitPrice,
		Subtotal:  row.Subtotal,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	q := sqlc.New(r.getDB(ctx))
	orderRow, err := q.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return domain.Order{}, err
	}

	itemRows, err := q.ListOrderItemsByOrderID(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}

	items := make([]domain.OrderItem, 0, len(itemRows))
	for _, row := range itemRows {
		items = append(items, domain.OrderItem{
			ID:        row.ID,
			OrderID:   row.OrderID,
			ProductID: row.ProductID,
			Quantity:  row.Quantity,
			UnitPrice: row.UnitPrice,
			Subtotal:  row.Subtotal,
			CreatedAt: row.CreatedAt,
		})
	}

	return domain.Order{
		ID:          orderRow.ID,
		UserID:      orderRow.UserID,
		TotalAmount: orderRow.TotalAmount,
		Status:      domain.OrderStatus(orderRow.Status),
		CreatedAt:   orderRow.CreatedAt,
		UpdatedAt:   orderRow.UpdatedAt,
		Items:       items,
	}, nil
}

func (r *OrderRepository) ListOrdersByUserID(ctx context.Context, userID uuid.UUID, page, limit int32) ([]domain.Order, error) {
	q := sqlc.New(r.getDB(ctx))
	offset := (page - 1) * limit
	rows, err := q.ListOrdersByUserID(ctx, sqlc.ListOrdersByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]domain.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, domain.Order{
			ID:          row.ID,
			UserID:      row.UserID,
			TotalAmount: row.TotalAmount,
			Status:      domain.OrderStatus(row.Status),
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return orders, nil
}
