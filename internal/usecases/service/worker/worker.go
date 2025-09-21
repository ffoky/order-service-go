package worker

import (
	"WBTECH_L0/internal/usecases/service/processor"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	id           int
	messagesChan <-chan []byte
	processor    *processor.Processor
}

func NewWorker(id int, messagesChan <-chan []byte, processor *processor.Processor) *Worker {
	return &Worker{id: id, messagesChan: messagesChan, processor: processor}
}

func (w *Worker) Start() {
	logrus.Infof("Worker %d started", w.id)
	defer logrus.Infof("Worker %d stopped", w.id)

	for message := range w.messagesChan {
		if message == nil {
			continue
		}

		if err := w.processor.ProcessMessage(message); err != nil {
			logrus.Errorf("Worker %d: final processing error (message sent to DLQ): %v", w.id, err)
		}
	}
}
