package cache

import (
	"WBTECH_L0/internal/domain"
	"context"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Order      *domain.Order
	Expiration int64
}

type Cache struct {
	mu                sync.RWMutex
	items             map[string]CacheItem
	defaultExpiration time.Duration
}

func NewCache(defaultExpiration time.Duration) *Cache {
	return &Cache{
		items:             make(map[string]CacheItem),
		defaultExpiration: defaultExpiration,
	}
}

func (c *Cache) Set(ctx context.Context, order *domain.Order, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if order == nil || order.OrderUID == "" {
		return fmt.Errorf("invalid order data")
	}
	if duration == 0 {
		duration = c.defaultExpiration
	}

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	c.items[order.OrderUID] = CacheItem{Order: order, Expiration: expiration}
	c.mu.Unlock()
	return nil
}

func (c *Cache) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.mu.RLock()
	item, exists := c.items[orderUID]
	c.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("order not found in cache")
	}
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, fmt.Errorf("order expired in cache")
	}
	return item.Order, nil
}

func (c *Cache) Has(ctx context.Context, orderUID string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	c.mu.RLock()
	item, exists := c.items[orderUID]
	c.mu.RUnlock()
	if !exists {
		return false
	}
	return item.Expiration <= 0 || time.Now().UnixNano() <= item.Expiration
}

func (c *Cache) Delete(ctx context.Context, orderUID string) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	c.mu.Lock()
	delete(c.items, orderUID)
	c.mu.Unlock()
}

func (c *Cache) LoadAll(ctx context.Context, orders []*domain.Order, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if ttl == 0 {
		ttl = c.defaultExpiration
	}
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	for _, order := range orders {
		if order != nil && order.OrderUID != "" {
			c.items[order.OrderUID] = CacheItem{
				Order:      order,
				Expiration: expiration,
			}
		}
	}
	return nil
}

func (c *Cache) GetAll(ctx context.Context) map[string]*domain.Order {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]*domain.Order)
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration <= 0 || now <= v.Expiration {
			result[k] = v.Order
		}
	}
	return result
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	now := time.Now().UnixNano()
	for _, v := range c.items {
		if v.Expiration <= 0 || now <= v.Expiration {
			count++
		}
	}
	return count
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]CacheItem)
	c.mu.Unlock()
}

func (c *Cache) Refresh(ctx context.Context, orderUID string, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[orderUID]
	if !ok {
		return fmt.Errorf("order not found in cache")
	}
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		item.Expiration = time.Now().Add(duration).UnixNano()
	}
	c.items[orderUID] = item
	return nil
}
