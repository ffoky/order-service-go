package generator

import (
	"WBTECH_L0/internal/domain"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
)

func (g *Generator) generateValidOrderUID() string {
	result := make([]byte, orderUIDLength)
	for i := 0; i < orderUIDLength; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result) + "test"
}

func (g *Generator) generateTrackNumber() string {
	return fmt.Sprintf("WBIL%sTRACK", gofakeit.LetterN(6))
}

func (g *Generator) generateRID() string {
	return gofakeit.LetterN(ridLength) + "test"
}

func (g *Generator) generateInvalidOrder() *domain.Order {
	order := g.generateOrder()

	switch rand.Intn(invalidCaseCount) {
	case 0:
		order.OrderUID = ""
	case 1:
		order.Payment.Amount = -100
	case 2:
		order.Items = []domain.Item{}
	case 3:
		order.Delivery.Email = "invalid-email"
	}

	return order
}
