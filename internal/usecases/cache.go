package usecases

import (
	"WBTECH_L0/internal/domain"
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, order *domain.Order, ttl time.Duration) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
	Has(ctx context.Context, orderUID string) bool
	LoadAll(ctx context.Context, orders []*domain.Order, ttl time.Duration) error
	GetAll(ctx context.Context) map[string]*domain.Order
	Size() int
	Clear()
	Delete(ctx context.Context, orderUID string)
	Refresh(ctx context.Context, orderUID string, ttl time.Duration) error
}
