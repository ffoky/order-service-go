package usecases

import "WBTECH_L0/internal/domain"

// Реализовал только Post потому что остальные пока не надо
// возможно понадобится Get
type Order interface {
	Post(order domain.Order) error
}
