package worker

import (
	"WBTECH_L0/internal/domain"
	"WBTECH_L0/internal/usecases"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	id           int
	messagesChan <-chan []byte
	orderService usecases.Order
}

func NewWorker(id int, messagesChan <-chan []byte, orderService usecases.Order) *Worker {
	return &Worker{
		id:           id,
		messagesChan: messagesChan,
		orderService: orderService,
	}
}

func (w *Worker) Start() {
	logrus.Infof("Worker %d started", w.id)

	for message := range w.messagesChan {
		if message == nil {
			continue
		}

		if err := w.processMessage(message); err != nil {
			logrus.Errorf("Worker %d failed to process message: %v", w.id, err)
		}
	}

	logrus.Infof("Worker %d stopped", w.id)
}

func (w *Worker) processMessage(message []byte) error {
	logrus.Infof("Worker %d processing message", w.id)

	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		logrus.Errorf("Worker %d failed to unmarshal order: %v", w.id, err)
		return err
	}

	ctx := context.Background()
	_, err := w.orderService.Create(ctx, &order)
	if err != nil {
		logrus.Errorf("Worker %d failed to save order %s: %v", w.id, order.OrderUID, err)
		return err
	}

	logrus.Infof("Worker %d successfully processed order %s", w.id, order.OrderUID)
	return nil
}
