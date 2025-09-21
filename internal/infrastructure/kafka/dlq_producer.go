package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type DLQMessage struct {
	OriginalMessage json.RawMessage `json:"original_message"`
	ErrorMessage    string          `json:"error_message"`
	Timestamp       time.Time       `json:"timestamp"`
	RetryCount      int             `json:"retry_count"`
	OriginalTopic   string          `json:"original_topic"`
}

type DLQProducer struct {
	producer *Producer
	dlqTopic string
}

func NewDLQProducer(brokerAddress []string, dlqTopic string) (*DLQProducer, error) {
	producer, err := NewProducer(brokerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create DLQ producer: %w", err)
	}

	return &DLQProducer{
		producer: producer,
		dlqTopic: dlqTopic,
	}, nil
}

func (d *DLQProducer) SendToDLQ(originalMessage []byte, errorMsg string, retryCount int, originalTopic string) error {
	dlqMsg := DLQMessage{
		OriginalMessage: json.RawMessage(originalMessage),
		ErrorMessage:    errorMsg,
		Timestamp:       time.Now(),
		RetryCount:      retryCount,
		OriginalTopic:   originalTopic,
	}

	msgBytes, err := json.Marshal(dlqMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", err)
	}

	if err := d.producer.Produce(string(msgBytes), d.dlqTopic); err != nil {
		logrus.Errorf("Failed to send message to DLQ: %v", err)
		return err
	}

	logrus.Warnf("Message sent to DLQ topic '%s' after %d retries. Error: %s", d.dlqTopic, retryCount, errorMsg)
	return nil
}

func (d *DLQProducer) Close() {
	if d.producer != nil {
		d.producer.Close()
	}
}
