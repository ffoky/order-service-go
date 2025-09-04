package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const cacheInitLimit = 10

type OrderService struct {
	repo  repository.Order
	cache repository.Cache
}

func NewOrderService(repo repository.Order, cache repository.Cache) *OrderService {
	return &OrderService{repo: repo, cache: cache}
}

func (o *OrderService) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	created, err := o.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	if o.cache != nil {
		_ = o.cache.Set(ctx, created, 0)
	}
	return created, nil
}

func (o *OrderService) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	if o.cache != nil {
		if order, err := o.cache.Get(ctx, orderUID); err == nil {
			return order, nil
		}
	}
	order, err := o.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	if o.cache != nil {
		_ = o.cache.Set(ctx, order, 0)
	}
	return order, nil
}

func (o *OrderService) InitializeCache(ctx context.Context, ttl time.Duration) error {
	if o.cache == nil {
		return nil
	}
	orders, err := o.repo.GetOrders(ctx, cacheInitLimit)
	if err != nil {
		return fmt.Errorf("failed to load orders: %w", err)
	}
	if err := o.cache.LoadAll(ctx, orders, ttl); err != nil {
		return fmt.Errorf("failed to cache orders: %w", err)
	}
	logrus.Infof("Cache initialized with %d orders", len(orders))
	return nil
}

func (o *OrderService) GetCacheStats() map[string]interface{} {
	if o.cache == nil {
		return map[string]interface{}{"cache_size": 0}
	}
	return map[string]interface{}{
		"cache_size": o.cache.Size(),
	}
}

func (o *OrderService) IsInCache(ctx context.Context, orderUID string) bool {
	if o.cache == nil {
		return false
	}
	return o.cache.Has(ctx, orderUID)
}

func (o *OrderService) RefreshCache(ctx context.Context, orderUID string, ttl time.Duration) error {
	if o.cache == nil {
		return fmt.Errorf("cache not enabled")
	}
	return o.cache.Refresh(ctx, orderUID, ttl)
}
