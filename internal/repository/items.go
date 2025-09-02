package repository

import (
	"WBTECH_L0/internal/domain"
	"context"
)

type Items interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	CreateItems(ctx context.Context, items []domain.Item) error
	GetItem(ctx context.Context, chrtID int) (*domain.Item, error)
	GetItemsByOrderUID(ctx context.Context, orderUID string) ([]domain.Item, error)
	CreateOrderItemLinks(ctx context.Context, orderUID string, itemChrtIDs []int) error
}
