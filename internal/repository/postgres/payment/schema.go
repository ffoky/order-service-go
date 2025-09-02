package payment

const paymentsTable = "payments"

const (
	paymentsTableColumnTransaction  = "transaction"
	paymentsTableColumnRequestID    = "request_id"
	paymentsTableColumnCurrency     = "currency"
	paymentsTableColumnProvider     = "provider"
	paymentsTableColumnAmount       = "amount"
	paymentsTableColumnPaymentDt    = "payment_dt"
	paymentsTableColumnBank         = "bank"
	paymentsTableColumnDeliveryCost = "delivery_cost"
	paymentsTableColumnGoodsTotal   = "goods_total"
	paymentsTableColumnCustomFee    = "custom_fee"
)

var paymentsTableColumns = []string{
	paymentsTableColumnTransaction,
	paymentsTableColumnRequestID,
	paymentsTableColumnCurrency,
	paymentsTableColumnProvider,
	paymentsTableColumnAmount,
	paymentsTableColumnPaymentDt,
	paymentsTableColumnBank,
	paymentsTableColumnDeliveryCost,
	paymentsTableColumnGoodsTotal,
	paymentsTableColumnCustomFee,
}
