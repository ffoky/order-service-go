package order

import (
	"WBTECH_L0/internal/domain"
	"database/sql"
	"time"
)

type OrderRow struct {
	OrderUID          string         `db:"order_uid"`
	DeliveryID        int64          `db:"delivery_id"`
	PaymentID         string         `db:"payment_id"`
	TrackNumber       string         `db:"track_number"`
	Entry             string         `db:"entry"`
	Locale            string         `db:"locale"`
	InternalSignature sql.NullString `db:"internal_signature"`
	CustomerID        string         `db:"customer_id"`
	DeliveryService   string         `db:"delivery_service"`
	Shardkey          string         `db:"shardkey"`
	SmID              int            `db:"sm_id"`
	DateCreated       time.Time      `db:"date_created"`
	OofShard          string         `db:"oof_shard"`
}

func (row *OrderRow) Values() []any {
	return []any{
		row.OrderUID,
		row.DeliveryID,
		row.PaymentID,
		row.TrackNumber,
		row.Entry,
		row.Locale,
		row.InternalSignature,
		row.CustomerID,
		row.DeliveryService,
		row.Shardkey,
		row.SmID,
		row.DateCreated,
		row.OofShard,
	}
}

func ToModel(r *OrderRow) *domain.Order {
	if r == nil {
		return nil
	}

	return &domain.Order{
		OrderUID:          r.OrderUID,
		TrackNumber:       r.TrackNumber,
		Entry:             r.Entry,
		Locale:            r.Locale,
		InternalSignature: r.InternalSignature.String,
		CustomerID:        r.CustomerID,
		DeliveryService:   r.DeliveryService,
		Shardkey:          r.Shardkey,
		SmID:              r.SmID,
		DateCreated:       r.DateCreated,
		OofShard:          r.OofShard,
	}
}

func FromModel(m *domain.Order, deliveryID int64, paymentID string) OrderRow {
	if m == nil {
		return OrderRow{}
	}

	return OrderRow{
		OrderUID:          m.OrderUID,
		DeliveryID:        deliveryID,
		PaymentID:         paymentID,
		TrackNumber:       m.TrackNumber,
		Entry:             m.Entry,
		Locale:            m.Locale,
		InternalSignature: sql.NullString{String: m.InternalSignature, Valid: m.InternalSignature != ""},
		CustomerID:        m.CustomerID,
		DeliveryService:   m.DeliveryService,
		Shardkey:          m.Shardkey,
		SmID:              m.SmID,
		DateCreated:       m.DateCreated,
		OofShard:          m.OofShard,
	}
}
