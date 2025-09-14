package delivery

import (
	"WBTECH_L0/internal/domain"
	"context"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) CreateDelivery(ctx context.Context, delivery *domain.Delivery) (int64, error) {
	row := FromModel(delivery)

	query := r.SB.
		Insert(deliveriesTable).
		Columns(deliveriesTableColumns[1:]...).
		Values(row.ValuesWithoutID()...).
		Suffix("RETURNING " + deliveriesTableColumnID)

	var deliveryID int64
	if err := r.Pool.Getx(ctx, &deliveryID, query); err != nil {
		return 0, err
	}

	return deliveryID, nil
}

func (r *Repository) GetDelivery(ctx context.Context, deliveryID int64) (*domain.Delivery, error) {
	query := r.SB.
		Select(deliveriesTableColumns...).
		From(deliveriesTable).
		Where(squirrel.Eq{deliveriesTableColumnID: deliveryID})

	var row DeliveryRow
	if err := r.Pool.Getx(ctx, &row, query); err != nil {
		return nil, err
	}

	return ToModel(&row), nil
}
