package order

import (
	"WBTECH_L0/internal/domain/dto"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"time"
)

var (
	entries          = []string{"WBIL", "WBRU", "WBUS", "WBEU"}
	locales          = []string{"en", "ru", "de", "fr"}
	currencies       = []string{"USD", "EUR", "RUB"}
	providers        = []string{"wbpay", "paypal", "stripe", "alpha"}
	banks            = []string{"alpha", "sber", "tinkoff", "vtb"}
	deliveryServices = []string{"meest", "dhl", "fedex", "ups", "russianpost"}
	brands           = []string{"Vivienne Sabo", "L'Oreal", "Maybelline", "Revlon", "MAC", "Chanel"}
	itemNames        = []string{"Mascara", "Lipstick", "Foundation", "Eyeshadow", "Blush", "Concealer"}
	sizes            = []string{"0", "XS", "S", "M", "L", "XL"}
)

// Generator структура для генерации заказов
type Generator struct{}

// NewGenerator создает новый экземпляр генератора
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateOrders генерирует указанное количество заказов
func (g *Generator) GenerateOrders(count int) []*dto.OrderDTO {
	orders := make([]*dto.OrderDTO, 0, count)

	for i := 0; i < count; i++ {
		var order *dto.OrderDTO

		// Иногда генерируем невалидные заказы (5% случаев)
		if rand.Intn(20) < 1 {
			order = g.generateInvalidOrder()
		} else {
			order = g.generateOrder()
		}

		orders = append(orders, order)
	}

	return orders
}

// generateValidOrderUID создает валидный ID заказа (16 случайных символов + "test")
func (g *Generator) generateValidOrderUID() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 16)
	for i := 0; i < 16; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result) + "test"
}

// generateTrackNumber создает номер трека в стиле "WBILMTESTTRACK"
func (g *Generator) generateTrackNumber() string {
	return fmt.Sprintf("WBIL%sTRACK", gofakeit.LetterN(6))
}

// generateRID создает RID для товара
func (g *Generator) generateRID() string {
	return gofakeit.LetterN(20) + "test"
}

// generateInvalidOrder создает невалидный заказ для тестирования обработки ошибок
func (g *Generator) generateInvalidOrder() *dto.OrderDTO {
	order := g.generateOrder()

	// Делаем заказ невалидным случайным образом
	switch rand.Intn(4) {
	case 0:
		order.OrderUID = "" // пустой UID
	case 1:
		order.Payment.Amount = -100 // отрицательная сумма
	case 2:
		order.Items = []dto.ItemDTO{} // пустой массив товаров
	case 3:
		order.Delivery.Email = "invalid-email" // невалидный email
	}

	return order
}

// generateOrder создает случайный заказ
func (g *Generator) generateOrder() *dto.OrderDTO {
	gofakeit.Seed(time.Now().UnixNano())

	orderUID := g.generateValidOrderUID()
	trackNumber := g.generateTrackNumber()

	// Генерируем товары (от 1 до 5 товаров)
	itemCount := rand.Intn(5) + 1
	items := make([]dto.ItemDTO, itemCount)
	goodsTotal := 0

	for i := 0; i < itemCount; i++ {
		price := rand.Intn(5000) + 100 // цена от 100 до 5100
		sale := rand.Intn(70)          // скидка от 0 до 70%
		totalPrice := price - (price * sale / 100)
		goodsTotal += totalPrice

		items[i] = dto.ItemDTO{
			ChrtID:      rand.Intn(10000000) + 1000000, // 7-значное число
			TrackNumber: trackNumber,
			Price:       price,
			RID:         g.generateRID(),
			Name:        gofakeit.RandomString(itemNames),
			Sale:        sale,
			Size:        gofakeit.RandomString(sizes),
			TotalPrice:  totalPrice,
			NmID:        rand.Intn(10000000) + 1000000,
			Brand:       gofakeit.RandomString(brands),
			Status:      []int{200, 201, 202, 404}[rand.Intn(4)],
		}
	}

	deliveryCost := rand.Intn(2000) + 500
	amount := goodsTotal + deliveryCost

	order := &dto.OrderDTO{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       gofakeit.RandomString(entries),
		Delivery: dto.DeliveryDTO{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: dto.PaymentDTO{
			Transaction:  orderUID, // transaction совпадает с order_uid
			RequestID:    "",       // оставляем пустым как в примере
			Currency:     gofakeit.RandomString(currencies),
			Provider:     gofakeit.RandomString(providers),
			Amount:       amount,
			PaymentDt:    time.Now().Unix(),
			Bank:         gofakeit.RandomString(banks),
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    rand.Intn(500), // от 0 до 500
		},
		Items:             items,
		Locale:            gofakeit.RandomString(locales),
		InternalSignature: "", // оставляем пустым как в примере
		CustomerID:        gofakeit.Username(),
		DeliveryService:   gofakeit.RandomString(deliveryServices),
		Shardkey:          fmt.Sprintf("%d", rand.Intn(10)), // от 0 до 9
		SmID:              rand.Intn(1000),                  // случайное число
		DateCreated:       time.Now(),
		OofShard:          fmt.Sprintf("%d", rand.Intn(10)), // от 0 до 9
	}

	return order
}
