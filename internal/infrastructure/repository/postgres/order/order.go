package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"WBTECH_L0/internal/domain"
)

func (r *Repository) CreateOrder(ctx context.Context, domainOrder *domain.Order) (*domain.Order, error) {
	var createdOrder *domain.Order

	err := r.tm.WithTx(ctx, func(txCtx context.Context) error {
		deliveryID, err := r.deliveryRepo.CreateDelivery(txCtx, &domainOrder.Delivery)
		if err != nil {
			return fmt.Errorf("failed to create delivery: %w", err)
		}

		err = r.paymentRepo.CreatePayment(txCtx, &domainOrder.Payment)
		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		err = r.itemsRepo.CreateItems(txCtx, domainOrder.Items)
		if err != nil {
			return fmt.Errorf("failed to create items: %w", err)
		}

		orderRow := FromModel(domainOrder, deliveryID, domainOrder.Payment.Transaction)
		query := r.SB.
			Insert(ordersTable).
			Columns(ordersTableColumns...).
			Values(orderRow.Values()...).
			Suffix("RETURNING " + ordersTableColumnID)

		var orderUID string
		if err := r.Pool.Getx(txCtx, &orderUID, query); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		var itemChrtIDs []int
		for _, item := range domainOrder.Items {
			itemChrtIDs = append(itemChrtIDs, item.ChrtID)
		}
		err = r.itemsRepo.CreateOrderItemLinks(txCtx, domainOrder.OrderUID, itemChrtIDs)
		if err != nil {
			return fmt.Errorf("failed to create order-item links: %w", err)
		}

		createdOrder = domainOrder
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdOrder, nil
}

func (r *Repository) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	var resultOrder *domain.Order

	err := r.tm.WithTx(ctx, func(txCtx context.Context) error {
		query := r.SB.
			Select(ordersTableColumns...).
			From(ordersTable).
			Where(squirrel.Eq{ordersTableColumnID: orderUID})

		var orderRow OrderRow
		if err := r.Pool.Getx(txCtx, &orderRow, query); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New("order not found")
			}
			return fmt.Errorf("failed to get order: %w", err)
		}

		delivery, err := r.deliveryRepo.GetDelivery(txCtx, orderRow.DeliveryID)
		if err != nil {
			return fmt.Errorf("failed to get delivery: %w", err)
		}

		payment, err := r.paymentRepo.GetPayment(txCtx, orderRow.PaymentID)
		if err != nil {
			return fmt.Errorf("failed to get payment: %w", err)
		}

		itemsList, err := r.itemsRepo.GetItemsByOrderUID(txCtx, orderUID)
		if err != nil {
			return fmt.Errorf("failed to get items: %w", err)
		}

		resultOrder = &domain.Order{
			OrderUID:          orderRow.OrderUID,
			TrackNumber:       orderRow.TrackNumber,
			Entry:             orderRow.Entry,
			Delivery:          *delivery,
			Payment:           *payment,
			Items:             itemsList,
			Locale:            orderRow.Locale,
			InternalSignature: orderRow.InternalSignature.String,
			CustomerID:        orderRow.CustomerID,
			DeliveryService:   orderRow.DeliveryService,
			Shardkey:          orderRow.Shardkey,
			SmID:              orderRow.SmID,
			DateCreated:       orderRow.DateCreated,
			OofShard:          orderRow.OofShard,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resultOrder, nil
}

func (r *Repository) GetOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	var orders []*domain.Order

	err := r.tm.WithTx(ctx, func(txCtx context.Context) error {
		query := r.SB.
			Select(ordersTableColumns...).
			From(ordersTable).
			Limit(uint64(limit))

		var orderRows []OrderRow
		if err := r.Pool.Selectx(txCtx, &orderRows, query); err != nil {
			return fmt.Errorf("failed to get all orders: %w", err)
		}

		for _, orderRow := range orderRows {
			delivery, err := r.deliveryRepo.GetDelivery(txCtx, orderRow.DeliveryID)
			if err != nil {
				return fmt.Errorf("failed to get delivery for order %s: %w", orderRow.OrderUID, err)
			}

			payment, err := r.paymentRepo.GetPayment(txCtx, orderRow.PaymentID)
			if err != nil {
				return fmt.Errorf("failed to get payment for order %s: %w", orderRow.OrderUID, err)
			}

			itemsList, err := r.itemsRepo.GetItemsByOrderUID(txCtx, orderRow.OrderUID)
			if err != nil {
				return fmt.Errorf("failed to get items for order %s: %w", orderRow.OrderUID, err)
			}

			order := &domain.Order{
				OrderUID:          orderRow.OrderUID,
				TrackNumber:       orderRow.TrackNumber,
				Entry:             orderRow.Entry,
				Delivery:          *delivery,
				Payment:           *payment,
				Items:             itemsList,
				Locale:            orderRow.Locale,
				InternalSignature: orderRow.InternalSignature.String,
				CustomerID:        orderRow.CustomerID,
				DeliveryService:   orderRow.DeliveryService,
				Shardkey:          orderRow.Shardkey,
				SmID:              orderRow.SmID,
				DateCreated:       orderRow.DateCreated,
				OofShard:          orderRow.OofShard,
			}

			orders = append(orders, order)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}
