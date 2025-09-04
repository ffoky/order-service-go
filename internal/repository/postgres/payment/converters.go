package payment

import (
	"WBTECH_L0/internal/domain"
	"database/sql"
)

type PaymentRow struct {
	Transaction  string         `db:"transaction"`
	RequestID    sql.NullString `db:"request_id"`
	Currency     string         `db:"currency"`
	Provider     string         `db:"provider"`
	Amount       float64        `db:"amount"`
	PaymentDt    int64          `db:"payment_dt"`
	Bank         string         `db:"bank"`
	DeliveryCost float64        `db:"delivery_cost"`
	GoodsTotal   float64        `db:"goods_total"`
	CustomFee    float64        `db:"custom_fee"`
}

func (row *PaymentRow) Values() []any {
	return []any{
		row.Transaction,
		row.RequestID,
		row.Currency,
		row.Provider,
		row.Amount,
		row.PaymentDt,
		row.Bank,
		row.DeliveryCost,
		row.GoodsTotal,
		row.CustomFee,
	}
}

func ToModel(r *PaymentRow) *domain.Payment {
	if r == nil {
		return nil
	}

	var requestID string
	if r.RequestID.Valid {
		requestID = r.RequestID.String
	}

	return &domain.Payment{
		Transaction:  r.Transaction,
		RequestID:    requestID,
		Currency:     r.Currency,
		Provider:     r.Provider,
		Amount:       r.Amount,
		PaymentDt:    r.PaymentDt,
		Bank:         r.Bank,
		DeliveryCost: r.DeliveryCost,
		GoodsTotal:   r.GoodsTotal,
		CustomFee:    r.CustomFee,
	}
}

func FromModel(m *domain.Payment) PaymentRow {
	if m == nil {
		return PaymentRow{}
	}

	return PaymentRow{
		Transaction:  m.Transaction,
		RequestID:    sql.NullString{String: m.RequestID, Valid: m.RequestID != ""},
		Currency:     m.Currency,
		Provider:     m.Provider,
		Amount:       m.Amount,
		PaymentDt:    m.PaymentDt,
		Bank:         m.Bank,
		DeliveryCost: m.DeliveryCost,
		GoodsTotal:   m.GoodsTotal,
		CustomFee:    m.CustomFee,
	}
}
