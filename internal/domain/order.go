package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Order представляет заказ в системе
// @Description Полная информация о заказе в системе WBTECH L0
type Order struct {
	OrderUID          string    `json:"order_uid" example:"b563feb7b2b84b6test"`
	TrackNumber       string    `json:"track_number" example:"WBILMTESTTRACK"`
	Entry             string    `json:"entry" example:"WBIL"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale" example:"en"`
	InternalSignature string    `json:"internal_signature" example:""`
	CustomerID        string    `json:"customer_id" example:"test"`
	DeliveryService   string    `json:"delivery_service" example:"meest"`
	Shardkey          string    `json:"shardkey" example:"9"`
	SmID              int       `json:"sm_id" example:"99"`
	DateCreated       time.Time `json:"date_created" example:"2021-11-26T06:22:19Z"`
	OofShard          string    `json:"oof_shard" example:"1"`
}

func ValidateOrder(order *Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	if err := validateOrderFields(order); err != nil {
		return err
	}

	if err := validateDelivery(&order.Delivery); err != nil {
		return fmt.Errorf("delivery validation: %w", err)
	}

	if err := validatePayment(&order.Payment); err != nil {
		return fmt.Errorf("payment validation: %w", err)
	}

	if err := validateItems(order.Items); err != nil {
		return fmt.Errorf("items validation: %w", err)
	}

	return nil
}

func validateOrderFields(order *Order) error {
	if strings.TrimSpace(order.OrderUID) == "" {
		return ErrInvalidOrderUID
	}

	if strings.TrimSpace(order.TrackNumber) == "" {
		return ErrInvalidTrackNumber
	}

	if strings.TrimSpace(order.Entry) == "" {
		return ErrInvalidEntry
	}

	if strings.TrimSpace(order.Locale) == "" {
		return ErrInvalidLocale
	}

	if strings.TrimSpace(order.CustomerID) == "" {
		return ErrInvalidCustomerID
	}

	if order.DateCreated.IsZero() {
		return ErrInvalidDateCreated
	}

	return nil
}
