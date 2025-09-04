package domain

import (
	"fmt"
	"strings"
)

type Item struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	RID         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int64   `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
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
