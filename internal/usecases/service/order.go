package orderservice

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const cacheInitLimit = 10

// OrderService отвечает за бизнес-логику работы с заказами.
// @description Сервис для управления заказами: создание, получение и работа с кэшем.
type OrderService struct {
	repo  repository.Order
	cache repository.Cache
}

// NewOrderService создает экземпляр сервиса.
// @summary Создание сервиса
// @param repo body repository.Order true "Репозиторий заказов"
// @param cache body repository.Cache true "Кэш"
// @return *OrderService
func NewOrderService(repo repository.Order, cache repository.Cache) *OrderService {
	return &OrderService{repo: repo, cache: cache}
}

// Create создает новый заказ
// @summary Создание заказа
// @description Валидирует и сохраняет заказ в БД и кэш.
// @param order body domain.Order true "Заказ"
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

	created, err := o.repo.CreateOrder(ctx, order)
	if err != nil {
		logrus.Errorf("Failed to create order %s: %v", order.OrderUID, err)
		return nil, err
	}

	if o.cache != nil {
		if err := o.cache.Set(ctx, created, 0); err != nil {
			logrus.Warnf("Failed to cache order %s: %v", created.OrderUID, err)
		}
	}

	logrus.Infof("Order %s successfully created", created.OrderUID)
	return created, nil
}

// Get возвращает заказ по UID.
// @summary Получение заказа
// @description Ищет заказ сначала в кэше, затем в БД.
// @param order_uid path string true "UID заказа"
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
		if err := o.cache.Set(ctx, order, 0); err != nil {
			logrus.Warnf("Failed to cache order %s: %v", order.OrderUID, err)
		}
	}

	return order, nil
}

// InitializeCache загружает заказы в кэш.
// @summary Инициализация кэша
// @description Загружает последние заказы из БД в кэш при старте сервиса.
// @param ttl query int false "Время жизни кэша в секундах"
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
		return map[string]interface{}{
			"cache_enabled": false,
			"cache_size":    0,
		}
	}

	return map[string]interface{}{
		"cache_enabled": true,
		"cache_size":    o.cache.Size(),
	}
}

// IsInCache проверяет наличие заказа в кэше.
// @summary Проверка заказа в кэше
// @param order_uid path string true "UID заказа"
func (o *OrderService) IsInCache(ctx context.Context, orderUID string) bool {
	if o.cache == nil || orderUID == "" {
		return false
	}
	return o.cache.Has(ctx, orderUID)
}

// RefreshCache обновляет заказ в кэше.
// @summary Обновление заказа в кэше
// @param order_uid path string true "UID заказа"
// @param ttl query int false "Время жизни кэша в секундах"
func (o *OrderService) RefreshCache(ctx context.Context, orderUID string, ttl time.Duration) error {
	if o.cache == nil {
		return fmt.Errorf("cache not enabled")
	}
	if orderUID == "" {
		return fmt.Errorf("order UID cannot be empty")
	}
	return o.cache.Refresh(ctx, orderUID, ttl)
}
