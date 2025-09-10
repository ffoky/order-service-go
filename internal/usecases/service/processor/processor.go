package processor

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/usecases"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Processor отвечает за бизнес-логику обработки заказов.
// @description Принимает сообщение, десериализует заказ и сохраняет его через OrderService.
type Processor struct {
	orderService usecases.Order
}

// NewProcessor создает новый экземпляр Processor.
// @summary Конструктор процессора заказов
// @param orderService body usecases.Order true "Сервис заказов"
// @return Processor
func NewProcessor(orderService usecases.Order) *Processor {
	return &Processor{orderService: orderService}
}

// ProcessMessage обрабатывает одно сообщение с заказом.
// @summary Обработка сообщения
// @description Десериализует JSON в структуру заказа и вызывает сервис для его сохранения.
// @param message body []byte true "Сообщение с заказом"
// @return error "Ошибка при десериализации или сохранении заказа"
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
