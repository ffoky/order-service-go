package delivery

const deliveriesTable = "deliveries"

const (
	deliveriesTableColumnID      = "delivery_id"
	deliveriesTableColumnName    = "name"
	deliveriesTableColumnPhone   = "phone"
	deliveriesTableColumnZip     = "zip"
	deliveriesTableColumnCity    = "city"
	deliveriesTableColumnAddress = "address"
	deliveriesTableColumnRegion  = "region"
	deliveriesTableColumnEmail   = "email"
)

var deliveriesTableColumns = []string{
	deliveriesTableColumnID,
	deliveriesTableColumnName,
	deliveriesTableColumnPhone,
	deliveriesTableColumnZip,
	deliveriesTableColumnCity,
	deliveriesTableColumnAddress,
	deliveriesTableColumnRegion,
	deliveriesTableColumnEmail,
}
