package generator

import (
	"WBTECH_L0/internal/domain"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateOrders(count int) []*domain.Order {
	orders := make([]*domain.Order, 0, count)

	for i := 0; i < count; i++ {
		var order *domain.Order

		if rand.Intn(count) < invalidOrderChance {
			order = g.generateInvalidOrder()
		} else {
			order = g.generateOrder()
		}

		orders = append(orders, order)
	}

	return orders
}

func (g *Generator) generateOrder() *domain.Order {
	gofakeit.Seed(time.Now().UnixNano())

	orderUID := g.generateValidOrderUID()
	trackNumber := g.generateTrackNumber()

	itemCount := rand.Intn(maxItemCount-minItemCount+1) + minItemCount
	items := make([]domain.Item, itemCount)
	goodsTotal := 0.0

	for i := 0; i < itemCount; i++ {
		price := float64(rand.Intn(maxPrice-minPrice+1) + minPrice)
		sale := rand.Intn(maxSalePercent + 1)
		totalPrice := price - (price * float64(sale) / 100)
		goodsTotal += totalPrice

		items[i] = domain.Item{
			ChrtID:      rand.Intn(maxChrtID-minChrtID+1) + minChrtID,
			TrackNumber: trackNumber,
			Price:       price,
			RID:         g.generateRID(),
			Name:        gofakeit.RandomString(itemNames),
			Sale:        sale,
			Size:        gofakeit.RandomString(sizes),
			TotalPrice:  totalPrice,
			NmID:        int64(rand.Intn(maxNmID-minNmID+1) + minNmID),
			Brand:       gofakeit.RandomString(brands),
			Status:      statusCodes[rand.Intn(len(statusCodes))],
		}
	}

	deliveryCost := float64(rand.Intn(maxDeliveryCost-minDeliveryCost+1) + minDeliveryCost)
	amount := goodsTotal + deliveryCost

	order := &domain.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       gofakeit.RandomString(entries),
		Delivery: domain.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: domain.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     gofakeit.RandomString(currencies),
			Provider:     gofakeit.RandomString(providers),
			Amount:       amount,
			PaymentDt:    time.Now().Unix(),
			Bank:         gofakeit.RandomString(banks),
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    float64(rand.Intn(maxCustomFee + 1)),
		},
		Items:             items,
		Locale:            gofakeit.RandomString(locales),
		InternalSignature: "",
		CustomerID:        gofakeit.Username(),
		DeliveryService:   gofakeit.RandomString(deliveryServices),
		Shardkey:          fmt.Sprintf("%d", rand.Intn(maxShardValue-minShardValue+1)+minShardValue),
		SmID:              rand.Intn(maxSmID + 1),
		DateCreated:       time.Now(),
		OofShard:          fmt.Sprintf("%d", rand.Intn(maxShardValue-minShardValue+1)+minShardValue),
	}

	return order
}
