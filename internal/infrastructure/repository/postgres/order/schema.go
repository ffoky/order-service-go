package order

const ordersTable = "orders"

const (
	ordersTableColumnID              = "order_uid"
	ordersTableColumnDeliveryID      = "delivery_id"
	ordersTableColumnPaymentID       = "payment_id"
	ordersTableColumnTrackNumber     = "track_number"
	ordersTableColumnEntry           = "entry"
	ordersTableColumnLocale          = "locale"
	ordersTableColumnInternalSig     = "internal_signature"
	ordersTableColumnCustomerID      = "customer_id"
	ordersTableColumnDeliveryService = "delivery_service"
	ordersTableColumnShardkey        = "shardkey"
	ordersTableColumnSmID            = "sm_id"
	ordersTableColumnDateCreated     = "date_created"
	ordersTableColumnOofShard        = "oof_shard"
)

var ordersTableColumns = []string{
	ordersTableColumnID,
	ordersTableColumnDeliveryID,
	ordersTableColumnPaymentID,
	ordersTableColumnTrackNumber,
	ordersTableColumnEntry,
	ordersTableColumnLocale,
	ordersTableColumnInternalSig,
	ordersTableColumnCustomerID,
	ordersTableColumnDeliveryService,
	ordersTableColumnShardkey,
	ordersTableColumnSmID,
	ordersTableColumnDateCreated,
	ordersTableColumnOofShard,
}
