package generator

const (
	orderUIDLength     = 16
	ridLength          = 20
	minItemCount       = 1
	maxItemCount       = 5
	minPrice           = 100
	maxPrice           = 5000
	maxSalePercent     = 70
	minChrtID          = 1000000
	maxChrtID          = 10000000
	minNmID            = 1000000
	maxNmID            = 10000000
	minDeliveryCost    = 500
	maxDeliveryCost    = 2000
	maxCustomFee       = 500
	minShardValue      = 0
	maxShardValue      = 9
	maxSmID            = 1000
	charset            = "abcdefghijklmnopqrstuvwxyz0123456789"
	invalidOrderChance = 1
	invalidCaseCount   = 4
)

var (
	entries          = []string{"WBIL", "WBRU", "WBUS", "WBEU"}
	locales          = []string{"en", "ru", "de", "fr"}
	currencies       = []string{"USD", "EUR", "RUB"}
	providers        = []string{"wbpay", "paypal", "stripe", "alpha"}
	banks            = []string{"alpha", "sber", "tinkoff", "vtb"}
	deliveryServices = []string{"fedex", "ups", "CDEK", "russianpost"}
	brands           = []string{"Vivienne Sabo", "L'Oreal", "Maybelline", "Revlon", "MAC", "Chanel"}
	itemNames        = []string{"Mascara", "Lipstick", "Foundation", "Eyeshadow", "Blush", "Concealer"}
	sizes            = []string{"0", "XS", "S", "M", "L", "XL"}
	statusCodes      = []int{200, 201, 202, 404}
)
