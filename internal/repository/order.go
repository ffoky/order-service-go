package repository

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Order interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
}
