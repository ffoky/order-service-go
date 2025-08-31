package order

import (
	"context"
	"errors"

	"WBTECH_L0/internal/domain"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	query := r.sb.
		Select(ordersTableColumns...).
		From(ordersTable).
		Where(squirrel.Eq{ordersTableColumnID: orderUID})

	var row OrderRow
	if err := r.pool.Getx(ctx, &row, query); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return ToModel(&row), nil
}
