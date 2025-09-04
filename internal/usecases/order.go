package usecases

import (
	"WBTECH_L0/internal/domain"
	"context"
	"time"
)

type Order interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
	InitializeCache(ctx context.Context, ttl time.Duration) error
	GetCacheStats() map[string]any
	IsInCache(ctx context.Context, orderUID string) bool
	RefreshCache(ctx context.Context, orderUID string, ttl time.Duration) error
}
