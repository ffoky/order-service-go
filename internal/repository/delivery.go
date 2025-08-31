package repository

import "WBTECH_L0/internal/domain"

type Delivery interface {
	Post(delivery domain.Delivery) error
	Get(deliveryId int) (*domain.Delivery, error)
}
