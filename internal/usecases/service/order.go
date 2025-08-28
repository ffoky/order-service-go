package service

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/repository"
)

type Order struct {
	repo repository.Order
}

func NewOrder(repo repository.Order) *Order {
	return &Order{
		repo: repo,
	}
}

func (o *Order) Post(order domain.Order) error {
	return o.repo.Post(order)
}
