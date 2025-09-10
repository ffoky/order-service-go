package domain

import "strings"

// Payment представляет информацию о платеже
// @Description Данные о платеже за заказ
type Payment struct {
	Transaction  string  `json:"transaction" example:"b563feb7b2b84b6test"`
	RequestID    string  `json:"request_id" example:""`
	Currency     string  `json:"currency" example:"USD"`
	Provider     string  `json:"provider" example:"wbpay"`
	Amount       float64 `json:"amount" example:"1817"`
	PaymentDt    int64   `json:"payment_dt" example:"1637907727"`
	Bank         string  `json:"bank" example:"alpha"`
	DeliveryCost float64 `json:"delivery_cost" example:"1500"`
	GoodsTotal   float64 `json:"goods_total" example:"317"`
	CustomFee    float64 `json:"custom_fee" example:"0"`
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
