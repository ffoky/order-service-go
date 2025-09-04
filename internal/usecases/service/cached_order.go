// internal/usecases/service/cached_order.go
package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

const cacheInitLimit = 10

type CachedOrder struct {
	repo  repository.Order
	cache repository.Cache
}

func NewCachedOrder(repo repository.Order, cache repository.Cache) *CachedOrder {
	return &CachedOrder{
		repo:  repo,
		cache: cache,
	}
}

func (o *CachedOrder) InitializeCache(ctx context.Context) error {
	logrus.Infof("Initializing cache with first %d orders from database...", cacheInitLimit)

	orders, err := o.repo.GetOrders(ctx, cacheInitLimit)
	if err != nil {
		return fmt.Errorf("failed to load orders from database: %w", err)
	}

	if err := o.cache.LoadAll(ctx, orders); err != nil {
		return fmt.Errorf("failed to load orders to cache: %w", err)
	}

	logrus.Infof("Cache initialized with %d orders", len(orders))
	return nil
}

func (o *CachedOrder) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	createdOrder, err := o.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order in database: %w", err)
	}

	if err := o.cache.Set(ctx, createdOrder); err != nil {
		logrus.Errorf("Failed to cache order %s: %v", createdOrder.OrderUID, err)
	}

	logrus.Infof("Order %s created and cached", createdOrder.OrderUID)
	return createdOrder, nil
}

func (o *CachedOrder) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	if order, err := o.cache.Get(ctx, orderUID); err == nil {
		logrus.Debugf("Order %s retrieved from cache", orderUID)
		return order, nil
	}

	logrus.Debugf("Order %s not found in cache, fetching from database", orderUID)
	order, err := o.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from database: %w", err)
	}

	if err := o.cache.Set(ctx, order); err != nil {
		logrus.Errorf("Failed to cache order %s: %v", order.OrderUID, err)
	}

	logrus.Debugf("Order %s retrieved from database and cached", orderUID)
	return order, nil
}

func (o *CachedOrder) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_size": o.cache.Size(),
	}
}

func (o *CachedOrder) IsInCache(ctx context.Context, orderUID string) bool {
	return o.cache.Has(ctx, orderUID)
}
