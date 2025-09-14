package order

import (
	"WBTECH_L0/internal/infrastructure/repository/postgres"
	"WBTECH_L0/internal/infrastructure/repository/postgres/delivery"
	"WBTECH_L0/internal/infrastructure/repository/postgres/items"
	"WBTECH_L0/internal/infrastructure/repository/postgres/payment"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	postgres.BaseRepository
	deliveryRepo delivery.Repository
	paymentRepo  payment.Repository
	itemsRepo    items.Repository
	tm           *postgres.TransactionManager
}

func NewRepository(
	pool *pgxpool.Pool,
	tm *postgres.TransactionManager,
	deliveryRepo delivery.Repository,
	paymentRepo payment.Repository,
	itemsRepo items.Repository,
) *Repository {
	return &Repository{
		BaseRepository: postgres.NewBaseRepository(pool),
		tm:             tm,
		deliveryRepo:   deliveryRepo,
		paymentRepo:    paymentRepo,
		itemsRepo:      itemsRepo,
	}
}
