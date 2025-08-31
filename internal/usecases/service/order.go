package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
	"context"
)

type Order struct {
	repo repository.Order
}

func NewOrder(repo repository.Order) *Order {
	return &Order{
		repo: repo,
	}
}

func (o *Order) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	return o.repo.CreateOrder(ctx, order)
}

func (o *Order) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	return o.repo.GetOrder(ctx, orderUID)
}
