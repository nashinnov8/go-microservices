package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/domain"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/dto"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/service"
)

type mockTxManager struct {
	committed  bool
	rolledBack bool
}

func (m *mockTxManager) WithTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	if err := fn(ctx); err != nil {
		m.rolledBack = true
		return err
	}
	m.committed = true
	return nil
}

type mockOrderRepo struct {
	orders       map[uuid.UUID]domain.Order
	orderItems   map[uuid.UUID][]domain.OrderItem
	failOnItemNo int // Simulate DB failure on Nth item
	itemCount    int
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		orders:     make(map[uuid.UUID]domain.Order),
		orderItems: make(map[uuid.UUID][]domain.OrderItem),
	}
}

func (m *mockOrderRepo) CreateOrder(_ context.Context, o domain.Order) (domain.Order, error) {
	o.ID = uuid.New()
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	m.orders[o.ID] = o
	return o, nil
}

func (m *mockOrderRepo) CreateOrderItem(_ context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	m.itemCount++
	if m.failOnItemNo > 0 && m.itemCount == m.failOnItemNo {
		return domain.OrderItem{}, errors.New("simulated database disk failure on item")
	}

	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	m.orderItems[item.OrderID] = append(m.orderItems[item.OrderID], item)
	return item, nil
}

func (m *mockOrderRepo) GetOrderByID(_ context.Context, id uuid.UUID) (domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	o.Items = m.orderItems[id]
	return o, nil
}

func (m *mockOrderRepo) ListOrdersByUserID(_ context.Context, _ uuid.UUID, _, _ int32) ([]domain.Order, error) {
	res := make([]domain.Order, 0, len(m.orders))
	for _, o := range m.orders {
		res = append(res, o)
	}
	return res, nil
}

func TestCreateOrder_Success_CommitsTransaction(t *testing.T) {
	repo := newMockOrderRepo()
	txMgr := &mockTxManager{}
	svc := service.NewOrderService(repo, txMgr)

	userID := uuid.New()
	prod1 := uuid.New().String()
	prod2 := uuid.New().String()

	req := dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: prod1, Quantity: 2, UnitPrice: 100}, // subtotal 200
			{ProductID: prod2, Quantity: 1, UnitPrice: 300}, // subtotal 300
		},
	}

	resp, err := svc.CreateOrder(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalAmount != 500 {
		t.Errorf("expected total_amount 500, got %d", resp.TotalAmount)
	}

	if len(resp.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(resp.Items))
	}

	if !txMgr.committed {
		t.Errorf("expected transaction to be committed")
	}

	if txMgr.rolledBack {
		t.Errorf("transaction should not be rolled back")
	}
}

func TestCreateOrder_ItemFails_RollsBackTransaction(t *testing.T) {
	repo := newMockOrderRepo()
	repo.failOnItemNo = 2 // 2nd item fails

	txMgr := &mockTxManager{}
	svc := service.NewOrderService(repo, txMgr)

	userID := uuid.New()
	prod1 := uuid.New().String()
	prod2 := uuid.New().String()

	req := dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: prod1, Quantity: 1, UnitPrice: 100},
			{ProductID: prod2, Quantity: 1, UnitPrice: 200},
		},
	}

	_, err := svc.CreateOrder(context.Background(), userID, req)
	if err == nil {
		t.Fatalf("expected error when item creation fails, got nil")
	}

	if !txMgr.rolledBack {
		t.Errorf("expected transaction to be rolled back on failure")
	}

	if txMgr.committed {
		t.Errorf("transaction should NOT have committed on failure")
	}
}
