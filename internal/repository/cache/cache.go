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
	insertOrder       []string
	defaultExpiration time.Duration
	maxSize           int
	requestCount      int
	maxRequestsPerSec int
	lastResetTime     time.Time
}

func NewCache(defaultExpiration time.Duration) *Cache {
	return &Cache{
		items:             make(map[string]CacheItem),
		insertOrder:       make([]string, 0),
		defaultExpiration: defaultExpiration,
		maxSize:           1000,
		maxRequestsPerSec: 10000,
		lastResetTime:     time.Now(),
	}
}

func (c *Cache) checkRateLimit() bool {
	now := time.Now()
	if now.Sub(c.lastResetTime) >= time.Second {
		c.requestCount = 0
		c.lastResetTime = now
	}

	c.requestCount++
	return c.requestCount <= c.maxRequestsPerSec
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

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.checkRateLimit() {
		return fmt.Errorf("rate limit exceeded")
	}

	if _, exists := c.items[order.OrderUID]; exists {
		return nil
	}

	if duration == 0 {
		duration = c.defaultExpiration
	}

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	if len(c.items) >= c.maxSize {
		c.deleteOldest()
	}

	c.items[order.OrderUID] = CacheItem{
		Order:      order,
		Expiration: expiration,
	}
	c.insertOrder = append(c.insertOrder, order.OrderUID)

	return nil
}

func (c *Cache) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.mu.RLock()
	if !c.checkRateLimit() {
		c.mu.RUnlock()
		return nil, fmt.Errorf("rate limit exceeded")
	}

	item, exists := c.items[orderUID]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("order not found in cache")
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.Delete(ctx, orderUID)
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
	defer c.mu.RUnlock()

	item, exists := c.items[orderUID]
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
	defer c.mu.Unlock()

	if _, exists := c.items[orderUID]; exists {
		delete(c.items, orderUID)

		for i, uid := range c.insertOrder {
			if uid == orderUID {
				c.insertOrder = append(c.insertOrder[:i], c.insertOrder[i+1:]...)
				break
			}
		}
	}
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
			if _, exists := c.items[order.OrderUID]; !exists {
				if len(c.items) >= c.maxSize {
					c.deleteOldest()
				}

				c.items[order.OrderUID] = CacheItem{
					Order:      order,
					Expiration: expiration,
				}
				c.insertOrder = append(c.insertOrder, order.OrderUID)
			}
		}
	}

	return nil
}

func (c *Cache) deleteOldest() {
	if len(c.insertOrder) > 0 {
		oldest := c.insertOrder[0]
		delete(c.items, oldest)
		c.insertOrder = c.insertOrder[1:]
	}
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
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
	c.insertOrder = make([]string, 0)
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
