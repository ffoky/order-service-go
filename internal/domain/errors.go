package domain

import (
	"errors"
)

var (
	ErrInvalidOrderUID    = errors.New("order UID is required")
	ErrInvalidTrackNumber = errors.New("track number is required")
	ErrInvalidEntry       = errors.New("entry is required")
	ErrInvalidLocale      = errors.New("locale is required")
	ErrInvalidCustomerID  = errors.New("customer ID is required")
	ErrInvalidDateCreated = errors.New("date created is required")
	ErrEmptyItems         = errors.New("order must contain at least one item")

	ErrInvalidDeliveryName = errors.New("delivery name is required")
	ErrInvalidPhone        = errors.New("invalid phone format")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidZip          = errors.New("zip code is required")
	ErrInvalidCity         = errors.New("city is required")
	ErrInvalidAddress      = errors.New("address is required")

	ErrInvalidTransaction   = errors.New("transaction ID is required")
	ErrInvalidCurrency      = errors.New("currency is required")
	ErrInvalidProvider      = errors.New("provider is required")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrInvalidPaymentDt     = errors.New("payment date is required")
	ErrInvalidBank          = errors.New("bank is required")
	ErrNegativeGoodsTotal   = errors.New("goods total cannot be negative")
	ErrNegativeDeliveryCost = errors.New("delivery cost cannot be negative")
	ErrNegativeCustomFee    = errors.New("custom fee cannot be negative")

	ErrInvalidChrtID     = errors.New("chrt_id must be positive")
	ErrInvalidItemName   = errors.New("item name is required")
	ErrInvalidPrice      = errors.New("price must be positive")
	ErrInvalidTotalPrice = errors.New("total price must be positive")
	ErrInvalidNmID       = errors.New("nm_id must be positive")
	ErrInvalidBrand      = errors.New("brand is required")
	ErrInvalidRID        = errors.New("rid is required")
	ErrNegativeSale      = errors.New("sale cannot be negative")
)
