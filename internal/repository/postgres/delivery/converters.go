package delivery

import (
	"WBTECH_L0/internal/domain"
)

type DeliveryRow struct {
	DeliveryID int64  `db:"delivery_id"`
	Name       string `db:"name"`
	Phone      string `db:"phone"`
	Zip        string `db:"zip"`
	City       string `db:"city"`
	Address    string `db:"address"`
	Region     string `db:"region"`
	Email      string `db:"email"`
}

// Values возвращает все значения полей (для SELECT с RETURNING).
func (row *DeliveryRow) Values() []any {
	return []any{
		row.DeliveryID,
		row.Name,
		row.Phone,
		row.Zip,
		row.City,
		row.Address,
		row.Region,
		row.Email,
	}
}

// ValuesWithoutID возвращает значения полей без delivery_id (для INSERT).
func (row *DeliveryRow) ValuesWithoutID() []any {
	return []any{
		row.Name,
		row.Phone,
		row.Zip,
		row.City,
		row.Address,
		row.Region,
		row.Email,
	}
}

// ToModel конвертирует DeliveryRow в доменную модель domain.Delivery.
func ToModel(r *DeliveryRow) *domain.Delivery {
	if r == nil {
		return nil
	}
	return &domain.Delivery{
		Name:    r.Name,
		Phone:   r.Phone,
		Zip:     r.Zip,
		City:    r.City,
		Address: r.Address,
		Region:  r.Region,
		Email:   r.Email,
	}
}

// FromModel конвертирует доменную модель в DeliveryRow (для INSERT/UPDATE).
func FromModel(m *domain.Delivery) DeliveryRow {
	if m == nil {
		return DeliveryRow{}
	}
	return DeliveryRow{
		Name:    m.Name,
		Phone:   m.Phone,
		Zip:     m.Zip,
		City:    m.City,
		Address: m.Address,
		Region:  m.Region,
		Email:   m.Email,
	}
}
