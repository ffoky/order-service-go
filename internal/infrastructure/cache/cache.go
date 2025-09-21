package cache

import (
	"WBTECH_L0/internal/domain"
	"container/list"
	"context"
	"fmt"
	"sync"
)

type CacheItem struct {
	Key   string
	Order *domain.Order
}

type LRUCache struct {
	mu       sync.RWMutex
	items    map[string]*list.Element
	queue    *list.List
	capacity int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *LRUCache) Get(ctx context.Context, key string) (*domain.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.queue.MoveToFront(elem)
		return elem.Value.(*CacheItem).Order, nil
	}

	return nil, fmt.Errorf("order not found in cache")
}

func (c *LRUCache) Set(ctx context.Context, order *domain.Order) error {
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

	c.setInternal(order.OrderUID, order)
	return nil
}

func (c *LRUCache) setInternal(key string, order *domain.Order) {
	if elem, exists := c.items[key]; exists {
		c.queue.MoveToFront(elem)
		elem.Value.(*CacheItem).Order = order
		return
	}

	if c.queue.Len() >= c.capacity {
		c.evict()
	}

	item := &CacheItem{
		Key:   key,
		Order: order,
	}

	elem := c.queue.PushFront(item)
	c.items[key] = elem
}

func (c *LRUCache) evict() {
	elem := c.queue.Back()
	if elem == nil {
		return
	}

	item := c.queue.Remove(elem).(*CacheItem)
	delete(c.items, item.Key)
}

func (c *LRUCache) Has(ctx context.Context, key string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.items[key]
	return exists
}

func (c *LRUCache) Delete(ctx context.Context, key string) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.queue.Remove(elem)
		delete(c.items, key)
	}
}

func (c *LRUCache) LoadAll(ctx context.Context, orders []*domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, order := range orders {
		if i%100 == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
		if order != nil && order.OrderUID != "" {
			if _, exists := c.items[order.OrderUID]; !exists {
				c.setInternal(order.OrderUID, order)
			}
		}
	}

	return nil
}
