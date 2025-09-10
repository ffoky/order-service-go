package processor

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/usecases"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	orderService usecases.Order
}

func NewProcessor(orderService usecases.Order) *Processor {
	return &Processor{orderService: orderService}
}

func (p *Processor) ProcessMessage(message []byte) error {
	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		logrus.Errorf("failed to unmarshal order: %v", err)
		return err
	}

	ctx := context.Background()
	_, err := p.orderService.Create(ctx, &order)
	if err != nil {
		logrus.Errorf("failed to save order %s: %v", order.OrderUID, err)
		return err
	}

	logrus.Infof("Order %s successfully saved", order.OrderUID)
	return nil
}
