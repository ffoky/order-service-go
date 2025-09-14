package items

import (
	"WBTECH_L0/internal/domain"
	"database/sql"
)

type ItemRow struct {
	ChrtID      int            `db:"chrt_id"`
	TrackNumber string         `db:"track_number"`
	Price       float64        `db:"price"`
	RID         sql.NullString `db:"rid"`
	Name        string         `db:"name"`
	Sale        sql.NullInt32  `db:"sale"`
	Size        string         `db:"size"`
	TotalPrice  float64        `db:"total_price"`
	NmID        int64          `db:"nm_id"`
	Brand       string         `db:"brand"`
	Status      int            `db:"status"`
}

func (row *ItemRow) Values() []any {
	return []any{
		row.ChrtID,
		row.TrackNumber,
		row.Price,
		row.RID,
		row.Name,
		row.Sale,
		row.Size,
		row.TotalPrice,
		row.NmID,
		row.Brand,
		row.Status,
	}
}

func ToModel(r *ItemRow) *domain.Item {
	if r == nil {
		return nil
	}

	var rid string
	if r.RID.Valid {
		rid = r.RID.String
	}

	var sale int
	if r.Sale.Valid {
		sale = int(r.Sale.Int32)
	}

	return &domain.Item{
		ChrtID:      r.ChrtID,
		TrackNumber: r.TrackNumber,
		Price:       r.Price,
		RID:         rid,
		Name:        r.Name,
		Sale:        sale,
		Size:        r.Size,
		TotalPrice:  r.TotalPrice,
		NmID:        r.NmID,
		Brand:       r.Brand,
		Status:      r.Status,
	}
}

func FromModel(m *domain.Item) ItemRow {
	if m == nil {
		return ItemRow{}
	}

	return ItemRow{
		ChrtID:      m.ChrtID,
		TrackNumber: m.TrackNumber,
		Price:       m.Price,
		RID:         sql.NullString{String: m.RID, Valid: m.RID != ""},
		Name:        m.Name,
		Sale:        sql.NullInt32{Int32: int32(m.Sale), Valid: true},
		Size:        m.Size,
		TotalPrice:  m.TotalPrice,
		NmID:        m.NmID,
		Brand:       m.Brand,
		Status:      m.Status,
	}
}
