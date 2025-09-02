package repository

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Payment interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPayment(ctx context.Context, transaction string) (*domain.Payment, error)
}
