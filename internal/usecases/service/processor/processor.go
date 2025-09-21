package processor

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/infrastructure/kafka"
	"WBTECH_L0/internal/usecases"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	orderService usecases.Order
	dlqProducer  *kafka.DLQProducer
	mainTopic    string
}

func NewProcessor(orderService usecases.Order, dlqProducer *kafka.DLQProducer, mainTopic string) *Processor {
	return &Processor{
		orderService: orderService,
		dlqProducer:  dlqProducer,
		mainTopic:    mainTopic,
	}
}

func (p *Processor) ProcessMessage(message []byte) error {
	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		return p.handleUnmarshalError(message, err)
	}

	ctx := context.Background()
	_, err := p.orderService.Create(ctx, &order)
	if err != nil {
		return p.handleCreateError(message, &order, err)
	}

	logrus.Infof("Order %s successfully processed", order.OrderUID)
	return nil
}

func (p *Processor) handleUnmarshalError(message []byte, err error) error {
	errorMsg := fmt.Sprintf("JSON unmarshal failed: %v", err)
	logrus.Errorf("Failed to parse message: %v", err)

	if p.dlqProducer != nil {
		if dlqErr := p.dlqProducer.SendToDLQ(message, errorMsg, 0, p.mainTopic); dlqErr != nil {
			logrus.Errorf("Failed to send unmarshal error to DLQ: %v", dlqErr)
		}
	} else {
		logrus.Warn("DLQ producer not available - message with unmarshal error will be lost")
	}

	return fmt.Errorf("message processing failed: %w", err)
}

func (p *Processor) handleCreateError(message []byte, order *domain.Order, err error) error {
	orderUID := "unknown"
	if order != nil && order.OrderUID != "" {
		orderUID = order.OrderUID
	}

	errorMsg := fmt.Sprintf("Order creation failed for %s: %v", orderUID, err)
	logrus.Errorf("Failed to create order %s: %v", orderUID, err)

	if p.dlqProducer != nil {
		if dlqErr := p.dlqProducer.SendToDLQ(message, errorMsg, 0, p.mainTopic); dlqErr != nil {
			logrus.Errorf("Failed to send creation error to DLQ: %v", dlqErr)
		}
	} else {
		logrus.Warn("DLQ producer not available - failed order will be lost")
	}

	return fmt.Errorf("order processing failed: %w", err)
}
