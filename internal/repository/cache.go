package repository

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Cache interface {
	Set(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
	Has(ctx context.Context, orderUID string) bool
	LoadAll(ctx context.Context, orders []*domain.Order) error
	GetAll(ctx context.Context) map[string]*domain.Order
	Size() int
	Clear()
}
