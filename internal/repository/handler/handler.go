package handler

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

// Вот вот здесь  может быть сохранение в базу данных

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message []byte, topic kafka.TopicPartition, cn int) error {
	logrus.Infof("Consumer №%d, Message from kafka with offset %d '%s on partition %d", topic.Offset, string(message), topic.Partition)
	return nil
}
