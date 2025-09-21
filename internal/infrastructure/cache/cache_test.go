package cache

import (
	"WBTECH_L0/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testOrder(uid string) *domain.Order {
	return &domain.Order{
		OrderUID:    uid,
		TrackNumber: "TEST" + uid,
		Entry:       "WBIL",
		CustomerID:  "customer_" + uid,
		DateCreated: time.Now(),
		Delivery: domain.Delivery{
			Name:  "Test User",
			Email: "test@test.com",
		},
		Payment: domain.Payment{
			Currency: "USD",
			Amount:   100,
		},
		Items: []domain.Item{{
			Name:  "Test Item",
			Price: 100,
		}},
	}
}

func TestLRUCache_Basic(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	order1 := testOrder("order1")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)

	result, err := cache.Get(ctx, "order1")
	assert.NoError(t, err)
	assert.Equal(t, "order1", result.OrderUID)
}

func TestLRUCache_NotFound(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	result, err := cache.Get(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found in cache")
}

func TestLRUCache_Has(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	assert.False(t, cache.Has(ctx, "order1"))

	order1 := testOrder("order1")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)

	assert.True(t, cache.Has(ctx, "order1"))
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	order1 := testOrder("order1")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)
	assert.True(t, cache.Has(ctx, "order1"))

	cache.Delete(ctx, "order1")
	assert.False(t, cache.Has(ctx, "order1"))
}

func TestLRUCache_InvalidData(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	err := cache.Set(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid order data")

	order := &domain.Order{OrderUID: ""}
	err = cache.Set(ctx, order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid order data")
}

func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	order1 := testOrder("order1")
	order2 := testOrder("order2")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)
	err = cache.Set(ctx, order2)
	assert.NoError(t, err)

	assert.True(t, cache.Has(ctx, "order1"))
	assert.True(t, cache.Has(ctx, "order2"))

	order3 := testOrder("order3")
	err = cache.Set(ctx, order3)
	assert.NoError(t, err)

	assert.False(t, cache.Has(ctx, "order1"))
	assert.True(t, cache.Has(ctx, "order2"))
	assert.True(t, cache.Has(ctx, "order3"))
}

func TestLRUCache_LRUOrder(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	order1 := testOrder("order1")
	order2 := testOrder("order2")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)
	err = cache.Set(ctx, order2)
	assert.NoError(t, err)

	_, err = cache.Get(ctx, "order1")
	assert.NoError(t, err)

	order3 := testOrder("order3")
	err = cache.Set(ctx, order3)
	assert.NoError(t, err)

	assert.True(t, cache.Has(ctx, "order1"))
	assert.False(t, cache.Has(ctx, "order2"))
	assert.True(t, cache.Has(ctx, "order3"))
}

func TestLRUCache_LoadAll(t *testing.T) {
	cache := NewLRUCache(5)
	ctx := context.Background()

	orders := []*domain.Order{
		testOrder("order1"),
		testOrder("order2"),
		testOrder("order3"),
	}

	err := cache.LoadAll(ctx, orders)
	assert.NoError(t, err)

	assert.True(t, cache.Has(ctx, "order1"))
	assert.True(t, cache.Has(ctx, "order2"))
	assert.True(t, cache.Has(ctx, "order3"))
}

func TestLRUCache_LoadAllWithNils(t *testing.T) {
	cache := NewLRUCache(5)
	ctx := context.Background()

	orders := []*domain.Order{
		testOrder("order1"),
		nil,
		testOrder("order2"),
		&domain.Order{OrderUID: ""},
	}

	err := cache.LoadAll(ctx, orders)
	assert.NoError(t, err)

	assert.True(t, cache.Has(ctx, "order1"))
	assert.True(t, cache.Has(ctx, "order2"))
}

func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache(2)
	ctx := context.Background()

	order1 := testOrder("order1")
	err := cache.Set(ctx, order1)
	assert.NoError(t, err)

	updatedOrder := testOrder("order1")
	updatedOrder.TrackNumber = "UPDATED"
	err = cache.Set(ctx, updatedOrder)
	assert.NoError(t, err)

	result, err := cache.Get(ctx, "order1")
	assert.NoError(t, err)
	assert.Equal(t, "UPDATED", result.TrackNumber)
}

func TestLRUCache_ContextCancellation(t *testing.T) {
	cache := NewLRUCache(2)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	order := testOrder("order1")

	err := cache.Set(ctx, order)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	_, err = cache.Get(ctx, "order1")
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	result := cache.Has(ctx, "order1")
	assert.False(t, result)
}
