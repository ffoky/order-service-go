package repository

import "WBTECH_L0/internal/domain"

type Order interface {
	Post(order domain.Order) error
}
