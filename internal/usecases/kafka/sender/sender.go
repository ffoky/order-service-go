package sender

import (
	"WBTECH_L0/internal/domain"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sync"
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

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for i, order := range orders {
		wg.Add(1)
		go func(orderIndex int, ord *domain.Order) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			s.sendOrder(orderIndex+1, ord, topic)
			time.Sleep(10 * time.Millisecond)
		}(i, order)
	}

	wg.Wait()
	logrus.Info("All messages sent")
	return nil
}

func (s *Sender) sendOrder(orderNum int, order *domain.Order, topic string) {
	isInvalid := s.isInvalidOrder(order)
	if isInvalid {
		logrus.Warnf("Sending invalid order #%d", orderNum)
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		logrus.Errorf("Error processing order #%d: %v", orderNum, err)
		return
	}

	if err = s.producer.Produce(string(jsonData), topic); err != nil {
		logrus.Errorf("Error occurred when sending order #%d: %v", orderNum, err)
	} else {
		logrus.Infof("Order #%d sent with ID: %s", orderNum, order.OrderUID)
	}
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
