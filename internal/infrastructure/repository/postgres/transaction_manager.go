package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

type TxKey struct{}

func (tm *TransactionManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if tm.IsTransactionActive(ctx, tx) {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logrus.WithError(rollbackErr).Warn("failed to rollback transaction")
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	if !tm.IsTransactionActive(ctx, tx) {
		return fmt.Errorf("transaction is no longer active")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (tm *TransactionManager) IsTransactionActive(ctx context.Context, tx pgx.Tx) bool {
	if tx == nil {
		return false
	}

	_, err := tx.Exec(ctx, "SELECT 1")
	return err == nil
}

func GetTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(TxKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
