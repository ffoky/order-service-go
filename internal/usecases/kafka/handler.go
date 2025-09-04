package processor

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/usecases"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	order usecases.Order
}

func NewHandler(order usecases.Order) *Handler {
	return &Handler{
		order: order,
	}
}

func (h *Handler) HandleMessage(message []byte, topic kafka.TopicPartition, cn int) error {
	logrus.Infof("Consumer №%d, Message from kafka offset %d on partition %d",
		cn, topic.Offset, topic.Partition)

	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		logrus.Errorf("failed to unmarshal order: %v", err)
		return err
	}

	ctx := context.Background()
	_, err := h.order.Create(ctx, &order)
	if err != nil {
		logrus.Errorf("failed to save order %s: %v", order.OrderUID, err)
		return err
	}

	logrus.Infof("Order %s successfully saved", order.OrderUID)
	return nil
}
