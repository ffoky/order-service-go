package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
	"context"
)

type BasicOrder struct {
	repo repository.Order
}

func NewBasicOrder(repo repository.Order) *BasicOrder {
	return &BasicOrder{
		repo: repo,
	}
}

func (o *BasicOrder) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	return o.repo.CreateOrder(ctx, order)
}

func (o *BasicOrder) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	return o.repo.GetOrder(ctx, orderUID)
}

func (o *BasicOrder) InitializeCache(ctx context.Context) error {
	return nil
}

func (o *BasicOrder) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_size": 0,
		"note":       "Basic service without cache",
	}
}

func (o *BasicOrder) IsInCache(ctx context.Context, orderUID string) bool {
	return false
}
