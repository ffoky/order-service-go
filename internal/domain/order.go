package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
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
