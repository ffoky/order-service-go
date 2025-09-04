package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	sessionTimeout = 7000 //ms
	noTimeout      = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(brokerAddress []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(brokerAddress, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
	}, nil
}

func (c *Consumer) StartReading(messagesChan chan<- []byte) {
	logrus.Info("Kafka consumer started reading messages...")

	for {
		if c.stop {
			break
		}

		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Errorf("Error reading message from Kafka: %v", err)
			continue
		}

		if kafkaMsg == nil {
			continue
		}

		select {
		case messagesChan <- kafkaMsg.Value:
			logrus.Debugf("Message sent to workers channel from partition %d, offset %d",
				kafkaMsg.TopicPartition.Partition, kafkaMsg.TopicPartition.Offset)
		default:
			logrus.Warn("Workers channel is full, message dropped")
		}

		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			logrus.Errorf("Error storing message: %v", err)
		}
	}

	logrus.Info("Kafka consumer stopped reading messages")
}

func (c *Consumer) Stop() error {
	logrus.Info("Stopping Kafka consumer...")
	c.stop = true

	if _, err := c.consumer.Commit(); err != nil {
		logrus.Errorf("Error committing offset: %v", err)
		return err
	}

	logrus.Info("Committed offset")
	return c.consumer.Close()
}
