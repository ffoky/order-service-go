package usecases

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Order interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
	InitializeCache(ctx context.Context) error
	GetCacheStats() map[string]any
	IsInCache(ctx context.Context, orderUID string) bool
}
