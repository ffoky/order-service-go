package items

const itemsTable = "items"

const (
	itemsTableColumnChrtID      = "chrt_id"
	itemsTableColumnTrackNumber = "track_number"
	itemsTableColumnPrice       = "price"
	itemsTableColumnRID         = "rid"
	itemsTableColumnName        = "name"
	itemsTableColumnSale        = "sale"
	itemsTableColumnSize        = "size"
	itemsTableColumnTotalPrice  = "total_price"
	itemsTableColumnNmID        = "nm_id"
	itemsTableColumnBrand       = "brand"
	itemsTableColumnStatus      = "status"
)

var itemsTableColumns = []string{
	itemsTableColumnChrtID,
	itemsTableColumnTrackNumber,
	itemsTableColumnPrice,
	itemsTableColumnRID,
	itemsTableColumnName,
	itemsTableColumnSale,
	itemsTableColumnSize,
	itemsTableColumnTotalPrice,
	itemsTableColumnNmID,
	itemsTableColumnBrand,
	itemsTableColumnStatus,
}

const orderItemsTable = "order_items"

const (
	orderItemsTableColumnOrderUID   = "order_uid"
	orderItemsTableColumnItemChrtID = "item_chrt_id"
)
