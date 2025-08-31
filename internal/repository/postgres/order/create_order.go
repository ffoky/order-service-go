package order

import (
	"context"
	"strings"

	"WBTECH_L0/internal/domain"
)

func (r *Repository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	row := FromModel(order)

	query := r.sb.
		Insert(ordersTable).
		Columns(ordersTableColumns...).
		Values(row.Values()...).
		Suffix("RETURNING " + strings.Join(ordersTableColumns, ","))

	var outRow OrderRow
	if err := r.pool.Getx(ctx, &outRow, query); err != nil {
		return nil, err
	}

	return ToModel(&outRow), nil
}
