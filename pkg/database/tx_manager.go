package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txContextKey struct{}

// TransactionManager defines the contract for executing business logic inside a transaction.
// The service layer depends ONLY on this interface, remaining 100% database-agnostic.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(txCtx context.Context) error) error
}

type pgxTxManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &pgxTxManager{pool: pool}
}

// WithTx starts a transaction, injects it into context, executes fn, and commits/rollbacks automatically.
func (m *pgxTxManager) WithTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	// If a transaction already exists in context, reuse it (supports nested WithTx calls)
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txContextKey{}, tx)

	// Defer rollback: if tx is already committed, Rollback returns pgx.ErrTxClosed (harmless)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(txCtx); err != nil {
		return err // deferred rollback will execute
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// TxFromContext retrieves a pgx.Tx from context if one was injected by TransactionManager.
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx, ok
}
