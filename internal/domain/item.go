package domain

import (
	"fmt"
	"strings"
)

// Item представляет товар в заказе
// @Description Информация о товаре в заказе
type Item struct {
	ChrtID      int     `json:"chrt_id" example:"9934930"`
	TrackNumber string  `json:"track_number" example:"WBILMTESTTRACK"`
	Price       float64 `json:"price" example:"453"`
	RID         string  `json:"rid" example:"ab4219087a764ae0btest"`
	Name        string  `json:"name" example:"Mascaras"`
	Sale        int     `json:"sale" example:"30"`
	Size        string  `json:"size" example:"0"`
	TotalPrice  float64 `json:"total_price" example:"317"`
	NmID        int64   `json:"nm_id" example:"2389212"`
	Brand       string  `json:"brand" example:"Vivienne Sabo"`
	Status      int     `json:"status" example:"202"`
}

func validateItems(items []Item) error {
	if len(items) == 0 {
		return ErrEmptyItems
	}

	for i, item := range items {
		if err := validateItem(&item); err != nil {
			return fmt.Errorf("item[%d]: %w", i, err)
		}
	}

	return nil
}

func validateItem(item *Item) error {
	if item.ChrtID <= 0 {
		return ErrInvalidChrtID
	}

	if strings.TrimSpace(item.Name) == "" {
		return ErrInvalidItemName
	}

	if item.Price <= 0 {
		return ErrInvalidPrice
	}

	if item.TotalPrice <= 0 {
		return ErrInvalidTotalPrice
	}

	if item.NmID <= 0 {
		return ErrInvalidNmID
	}

	if strings.TrimSpace(item.Brand) == "" {
		return ErrInvalidBrand
	}

	if strings.TrimSpace(item.RID) == "" {
		return ErrInvalidRID
	}

	if item.Sale < 0 {
		return ErrNegativeSale
	}

	return nil
}
