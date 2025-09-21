package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createValidOrder(orderUID string) *domain.Order {
	return &domain.Order{
		OrderUID:    orderUID,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test Petrovich",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: domain.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "1",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func TestOrderService_Create(t *testing.T) {
	tests := []struct {
		name        string
		order       *domain.Order
		setupMocks  func(*mocks.Order, *mocks.Cache)
		expectError bool
		errorText   string
	}{
		{
			name:        "nil order",
			order:       nil,
			setupMocks:  func(repo *mocks.Order, cache *mocks.Cache) {},
			expectError: true,
			errorText:   "order cannot be nil",
		},
		{
			name: "validation error - empty OrderUID",
			order: &domain.Order{
				OrderUID: "",
			},
			setupMocks:  func(repo *mocks.Order, cache *mocks.Cache) {},
			expectError: true,
			errorText:   "validation failed",
		},
		{
			name:  "order exists in cache",
			order: createValidOrder("cached-order"),
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				cachedOrder := createValidOrder("cached-order")
				cache.EXPECT().Has(mock.Anything, "cached-order").Return(true)
				cache.EXPECT().Get(mock.Anything, "cached-order").Return(cachedOrder, nil)
			},
			expectError: false,
		},
		{
			name:  "successful creation",
			order: createValidOrder("new-order"),
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				expectedOrder := createValidOrder("new-order")
				cache.EXPECT().Has(mock.Anything, "new-order").Return(false)
				repo.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(expectedOrder, nil)
				cache.EXPECT().Set(mock.Anything, expectedOrder).Return(nil)
			},
			expectError: false,
		},
		{
			name:  "repository error",
			order: createValidOrder("repo-error"),
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				cache.EXPECT().Has(mock.Anything, "repo-error").Return(false)
				repo.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectError: true,
			errorText:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrder(t)
			mockCache := mocks.NewCache(t)
			service := NewOrderService(mockRepo, mockCache)

			tt.setupMocks(mockRepo, mockCache)

			result, err := service.Create(context.Background(), tt.order)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestOrderService_Get(t *testing.T) {
	tests := []struct {
		name        string
		orderUID    string
		setupMocks  func(*mocks.Order, *mocks.Cache)
		expectError bool
		errorText   string
	}{
		{
			name:        "empty OrderUID",
			orderUID:    "",
			setupMocks:  func(repo *mocks.Order, cache *mocks.Cache) {},
			expectError: true,
			errorText:   "order UID cannot be empty",
		},
		{
			name:     "get from cache",
			orderUID: "cached-uid",
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				cachedOrder := createValidOrder("cached-uid")
				cache.EXPECT().Get(mock.Anything, "cached-uid").Return(cachedOrder, nil)
			},
			expectError: false,
		},
		{
			name:     "get from repository",
			orderUID: "repo-uid",
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				repoOrder := createValidOrder("repo-uid")
				cache.EXPECT().Get(mock.Anything, "repo-uid").Return(nil, errors.New("not in cache"))
				repo.EXPECT().GetOrder(mock.Anything, "repo-uid").Return(repoOrder, nil)
				cache.EXPECT().Set(mock.Anything, repoOrder).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "order not found",
			orderUID: "nonexistent-uid",
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				cache.EXPECT().Get(mock.Anything, "nonexistent-uid").Return(nil, errors.New("not in cache"))
				repo.EXPECT().GetOrder(mock.Anything, "nonexistent-uid").Return(nil, errors.New("order not found"))
			},
			expectError: true,
			errorText:   "order not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrder(t)
			mockCache := mocks.NewCache(t)
			service := NewOrderService(mockRepo, mockCache)

			tt.setupMocks(mockRepo, mockCache)

			result, err := service.Get(context.Background(), tt.orderUID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.orderUID, result.OrderUID)
			}
		})
	}
}

func TestOrderService_InitializeCache(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mocks.Order, *mocks.Cache)
		expectError bool
		errorText   string
	}{
		{
			name: "successful initialization",
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				orders := []*domain.Order{
					createValidOrder("order-1"),
					createValidOrder("order-2"),
				}
				repo.EXPECT().GetOrders(mock.Anything, cacheInitLimit).Return(orders, nil)
				cache.EXPECT().LoadAll(mock.Anything, orders).Return(nil)
			},
			expectError: false,
		},
		{
			name: "repository error",
			setupMocks: func(repo *mocks.Order, cache *mocks.Cache) {
				repo.EXPECT().GetOrders(mock.Anything, cacheInitLimit).Return(nil, errors.New("database error"))
			},
			expectError: true,
			errorText:   "failed to load orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrder(t)
			mockCache := mocks.NewCache(t)
			service := NewOrderService(mockRepo, mockCache)

			tt.setupMocks(mockRepo, mockCache)

			err := service.InitializeCache(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderService_IsInCache(t *testing.T) {
	mockRepo := mocks.NewOrder(t)
	mockCache := mocks.NewCache(t)
	service := NewOrderService(mockRepo, mockCache)

	result := service.IsInCache(context.Background(), "")
	assert.False(t, result)

	mockCache.EXPECT().Has(mock.Anything, "test-uid").Return(true).Once()
	result = service.IsInCache(context.Background(), "test-uid")
	assert.True(t, result)
}
