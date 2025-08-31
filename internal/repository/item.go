package repository

import "WBTECH_L0/internal/domain"

type Item interface {
	Post(item domain.Item) error
	Get(chrtID int) (*domain.Item, error)
}
