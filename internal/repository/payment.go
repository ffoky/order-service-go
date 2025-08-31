package repository

import "WBTECH_L0/internal/domain"

type Payment interface {
	Post(payment domain.Payment) error
	Get(paymentId string) (*domain.Payment, error)
}
