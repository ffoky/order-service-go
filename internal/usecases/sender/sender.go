package sender

import (
	"WBTECH_L0/internal/domain"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type Producer interface {
	Produce(message, topic string) error
	Close()
}

type Sender struct {
	producer Producer
}

func NewSender(producer Producer) *Sender {
	return &Sender{
		producer: producer,
	}
}

func (s *Sender) SendOrders(orders []*domain.Order, topic string) error {
	logrus.Info("Sending orders...")

	for i, order := range orders {
		isInvalid := s.isInvalidOrder(order)
		if isInvalid {
			logrus.Warnf("Sending invalid order #%d", i+1)
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			logrus.Errorf("Error proccessing order #%d: %v", i+1, err)
			continue
		}

		if err = s.producer.Produce(string(jsonData), topic); err != nil {
			logrus.Errorf("Error occured when sending order #%d: %v", i+1, err)
		} else {
			logrus.Infof("Order #%d sent with (ID: %s)", i+1, order.OrderUID)
		}

		time.Sleep(100 * time.Millisecond)
	}

	logrus.Info("All messages sent")
	return nil
}

func (s *Sender) isInvalidOrder(order *domain.Order) bool {
	if order.OrderUID == "" {
		return true
	}
	if order.Payment.Amount < 0 {
		return true
	}
	if len(order.Items) == 0 {
		return true
	}
	if order.Delivery.Email == "invalid-email" {
		return true
	}
	return false
}
