package domain

import "strings"

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

func validatePayment(payment *Payment) error {
	if strings.TrimSpace(payment.Transaction) == "" {
		return ErrInvalidTransaction
	}

	if strings.TrimSpace(payment.Currency) == "" {
		return ErrInvalidCurrency
	}

	if strings.TrimSpace(payment.Provider) == "" {
		return ErrInvalidProvider
	}

	if payment.Amount <= 0 {
		return ErrInvalidAmount
	}

	if payment.PaymentDt <= 0 {
		return ErrInvalidPaymentDt
	}

	if strings.TrimSpace(payment.Bank) == "" {
		return ErrInvalidBank
	}

	if payment.GoodsTotal < 0 {
		return ErrNegativeGoodsTotal
	}

	if payment.DeliveryCost < 0 {
		return ErrNegativeDeliveryCost
	}

	if payment.CustomFee < 0 {
		return ErrNegativeCustomFee
	}

	return nil
}
