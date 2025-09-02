package items

import (
	"WBTECH_L0/internal/domain"
	"context"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) CreateItem(ctx context.Context, item *domain.Item) error {
	row := FromModel(item)

	query := r.SB.
		Insert(itemsTable).
		Columns(itemsTableColumns...).
		Values(row.Values()...)

	_, err := r.Pool.Execx(ctx, query)
	return err
}

func (r *Repository) CreateItems(ctx context.Context, items []domain.Item) error {
	if len(items) == 0 {
		return nil
	}

	query := r.SB.Insert(itemsTable).Columns(itemsTableColumns...)

	for _, item := range items {
		row := FromModel(&item)
		query = query.Values(row.Values()...)
	}

	_, err := r.Pool.Execx(ctx, query)
	return err
}

func (r *Repository) GetItem(ctx context.Context, chrtID int) (*domain.Item, error) {
	query := r.SB.
		Select(itemsTableColumns...).
		From(itemsTable).
		Where(squirrel.Eq{itemsTableColumnChrtID: chrtID})

	var row ItemRow
	if err := r.Pool.Getx(ctx, &row, query); err != nil {
		return nil, err
	}

	return ToModel(&row), nil
}

func (r *Repository) GetItemsByOrderUID(ctx context.Context, orderUID string) ([]domain.Item, error) {
	query := r.SB.
		Select(itemsTableColumns...).
		From(itemsTable).
		Join("order_items oi ON items.chrt_id = oi.item_chrt_id").
		Where(squirrel.Eq{"oi.order_uid": orderUID})

	var rows []ItemRow
	if err := r.Pool.Selectx(ctx, &rows, query); err != nil {
		return nil, err
	}

	items := make([]domain.Item, len(rows))
	for i, row := range rows {
		items[i] = *ToModel(&row)
	}

	return items, nil
}

func (r *Repository) CreateOrderItemLinks(ctx context.Context, orderUID string, itemChrtIDs []int) error {
	if len(itemChrtIDs) == 0 {
		return nil
	}

	query := r.SB.Insert(orderItemsTable).
		Columns(orderItemsTableColumnOrderUID, orderItemsTableColumnItemChrtID)

	for _, chrtID := range itemChrtIDs {
		query = query.Values(orderUID, chrtID)
	}

	_, err := r.Pool.Execx(ctx, query)
	return err
}
