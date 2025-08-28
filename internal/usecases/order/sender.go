package order

import (
	"WBTECH_L0/internal/domain/dto"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

// Producer интерфейс для отправки сообщений в Kafka
type Producer interface {
	Produce(message, topic string) error
	Close()
}

// Sender структура для отправки заказов в Kafka
type Sender struct {
	producer Producer
}

// NewSender создает новый экземпляр отправителя
func NewSender(producer Producer) *Sender {
	return &Sender{
		producer: producer,
	}
}

// SendOrders отправляет заказы в указанный топик Kafka
func (s *Sender) SendOrders(orders []*dto.OrderDTO, topic string) error {
	logrus.Info("Запуск отправки заказов...")

	for i, order := range orders {
		// Проверяем, является ли заказ невалидным
		isInvalid := s.isInvalidOrder(order)
		if isInvalid {
			logrus.Warnf("Отправляем невалидный заказ #%d для тестирования", i+1)
		}

		// Сериализуем в JSON
		jsonData, err := json.Marshal(order)
		if err != nil {
			logrus.Errorf("Ошибка сериализации заказа #%d: %v", i+1, err)
			continue
		}

		// Отправляем в Kafka
		if err = s.producer.Produce(string(jsonData), topic); err != nil {
			logrus.Errorf("Ошибка отправки заказа #%d: %v", i+1, err)
		} else {
			logrus.Infof("Заказ #%d отправлен (ID: %s)", i+1, order.OrderUID)
		}

		// Небольшая пауза между отправками
		time.Sleep(100 * time.Millisecond)
	}

	logrus.Info("Все сообщения успешно отправлены.")
	return nil
}

// isInvalidOrder проверяет, является ли заказ невалидным
func (s *Sender) isInvalidOrder(order *dto.OrderDTO) bool {
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
