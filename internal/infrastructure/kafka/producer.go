package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	flushTimeout = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokerAddress []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokerAddress, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new producer: %w", err)
	}

	producer := &Producer{producer: p}
	go producer.handleDeliveryReports()

	return producer, nil
}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   nil,
	}

	return p.producer.Produce(kafkaMsg, nil)
}

func (p *Producer) handleDeliveryReports() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				logrus.Errorf("Delivery failed: %v", ev.TopicPartition.Error)
			} else {
				logrus.Debugf("Message delivered to partition %d at offset %v",
					ev.TopicPartition.Partition, ev.TopicPartition.Offset)
			}
		}
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
