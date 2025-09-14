package repository

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Delivery interface {
	CreateDelivery(ctx context.Context, delivery *domain.Delivery) (int64, error)
	GetDelivery(ctx context.Context, deliveryID int64) (*domain.Delivery, error)
}
