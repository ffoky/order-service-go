package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/infrastructure/repository"
	"WBTECH_L0/internal/usecases"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

const cacheInitLimit = 10

type OrderService struct {
	repo  repository.Order
	cache usecases.Cache
}

func NewOrderService(repo repository.Order, cache usecases.Cache) *OrderService {
	return &OrderService{repo: repo, cache: cache}
}

func (o *OrderService) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if order == nil {
		logrus.Error("Received nil order")
		return nil, fmt.Errorf("order cannot be nil")
	}

	if err := domain.ValidateOrder(order); err != nil {
		logrus.WithFields(logrus.Fields{
			"order_uid": order.OrderUID,
			"error":     err.Error(),
		}).Error("Order validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if o.cache != nil && o.cache.Has(ctx, order.OrderUID) {
		logrus.Infof("Order %s already exists in cache", order.OrderUID)
		return o.cache.Get(ctx, order.OrderUID)
	}

	createdOrder, err := o.repo.CreateOrder(ctx, order)
	if err != nil {
		logrus.Errorf("Failed to create order %s: %v", order.OrderUID, err)
		return nil, err
	}

	if o.cache != nil {
		if err := o.cache.Set(ctx, createdOrder); err != nil {
			logrus.Warnf("Failed to cache order %s: %v", createdOrder.OrderUID, err)
		}
	}

	logrus.Infof("Order %s successfully createdOrder", createdOrder.OrderUID)
	return createdOrder, nil
}

func (o *OrderService) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	if orderUID == "" {
		return nil, fmt.Errorf("order UID cannot be empty")
	}

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
		if err := o.cache.Set(ctx, order); err != nil {
			logrus.Warnf("Failed to cache order %s: %v", order.OrderUID, err)
		}
	}

	return order, nil
}

func (o *OrderService) InitializeCache(ctx context.Context) error {
	if o.cache == nil {
		return nil
	}

	orders, err := o.repo.GetOrders(ctx, cacheInitLimit)
	if err != nil {
		return fmt.Errorf("failed to load orders: %w", err)
	}

	if err := o.cache.LoadAll(ctx, orders); err != nil {
		return fmt.Errorf("failed to cache orders: %w", err)
	}

	logrus.Infof("Cache initialized with %d orders", len(orders))
	return nil
}

func (o *OrderService) IsInCache(ctx context.Context, orderUID string) bool {
	if o.cache == nil || orderUID == "" {
		return false
	}
	return o.cache.Has(ctx, orderUID)
}
