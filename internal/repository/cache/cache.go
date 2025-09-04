package cache

import (
	"WBTECH_L0/internal/domain"
	"context"
	"fmt"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]*domain.Order
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*domain.Order),
	}
}

func (c *Cache) Set(ctx context.Context, order *domain.Order) error {
	if order == nil || order.OrderUID == "" {
		return fmt.Errorf("invalid order data")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[order.OrderUID] = order
	return nil
}

func (c *Cache) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, exists := c.items[orderUID]
	if !exists {
		return nil, fmt.Errorf("order not found in cache")
	}

	return order, nil
}

func (c *Cache) Has(ctx context.Context, orderUID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.items[orderUID]
	return exists
}

func (c *Cache) LoadAll(ctx context.Context, orders []*domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		if order != nil && order.OrderUID != "" {
			c.items[order.OrderUID] = order
		}
	}

	return nil
}

func (c *Cache) GetAll(ctx context.Context) map[string]*domain.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*domain.Order)
	for k, v := range c.items {
		result[k] = v
	}

	return result
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*domain.Order)
}
