package usecases

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Cache interface {
	Get(ctx context.Context, key string) (*domain.Order, error)
	Set(ctx context.Context, order *domain.Order) error
	Has(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string)
	LoadAll(ctx context.Context, orders []*domain.Order) error
}
