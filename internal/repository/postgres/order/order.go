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
		// 1. Сохраняем Delivery
		deliveryID, err := r.deliveryRepo.CreateDelivery(txCtx, &domainOrder.Delivery) // Исправляем на CreateDelivery
		if err != nil {
			return fmt.Errorf("failed to create delivery: %w", err)
		}

		// 2. Сохраняем Payment
		err = r.paymentRepo.CreatePayment(txCtx, &domainOrder.Payment) // Исправляем на CreatePayment
		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		// 3. Сохраняем Items
		err = r.itemsRepo.CreateItems(txCtx, domainOrder.Items)
		if err != nil {
			return fmt.Errorf("failed to create items: %w", err)
		}

		// 4. Подготавливаем и сохраняем запись Order
		orderRow := FromModel(domainOrder, deliveryID, domainOrder.Payment.Transaction)
		query := r.SB.
			Insert(ordersTable).
			Columns(ordersTableColumns...).
			Values(orderRow.Values()...).
			Suffix("RETURNING " + ordersTableColumnID) // Возвращаем только ID для подтверждения

		var orderUID string
		if err := r.Pool.Getx(txCtx, &orderUID, query); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 5. Сохраняем связи между заказом и товарами
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
		// 1. Достаем основную запись заказа
		query := r.SB.
			Select(ordersTableColumns...). // В columns теперь входят delivery_id и payment_id
			From(ordersTable).
			Where(squirrel.Eq{ordersTableColumnID: orderUID})

		var orderRow OrderRow // OrderRow теперь содержит DeliveryID и PaymentID
		if err := r.Pool.Getx(txCtx, &orderRow, query); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New("order not found")
			}
			return fmt.Errorf("failed to get order: %w", err)
		}

		// 2. Загружаем Delivery
		delivery, err := r.deliveryRepo.GetDelivery(txCtx, orderRow.DeliveryID) // Исправляем на GetDelivery
		if err != nil {
			return fmt.Errorf("failed to get delivery: %w", err)
		}

		// 3. Загружаем Payment
		payment, err := r.paymentRepo.GetPayment(txCtx, orderRow.PaymentID)
		if err != nil {
			return fmt.Errorf("failed to get payment: %w", err)
		}

		// 4. Загружаем Items
		itemsList, err := r.itemsRepo.GetItemsByOrderUID(txCtx, orderUID)
		if err != nil {
			return fmt.Errorf("failed to get items: %w", err)
		}

		// 5. Собираем полный агрегат
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
