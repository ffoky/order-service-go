package payment

import (
	"WBTECH_L0/internal/domain"
	"context"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	row := FromModel(payment)

	query := r.SB.
		Insert(paymentsTable).
		Columns(paymentsTableColumns...).
		Values(row.Values()...)

	_, err := r.Pool.Execx(ctx, query)
	return err
}

func (r *Repository) GetPayment(ctx context.Context, transaction string) (*domain.Payment, error) {
	query := r.SB.
		Select(paymentsTableColumns...).
		From(paymentsTable).
		Where(squirrel.Eq{paymentsTableColumnTransaction: transaction})

	var row PaymentRow
	if err := r.Pool.Getx(ctx, &row, query); err != nil {
		return nil, err
	}

	return ToModel(&row), nil
}
