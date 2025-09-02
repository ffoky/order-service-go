package payment

import (
	"WBTECH_L0/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	postgres.BaseRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		BaseRepository: postgres.NewBaseRepository(pool),
	}
}
