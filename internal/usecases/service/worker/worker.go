package worker

import (
	"WBTECH_L0/internal/usecases/service/processor"
	"github.com/sirupsen/logrus"
)

// Worker обрабатывает сообщения из кафки и сохраняет заказы через OrderService.
// @description Воркер отвечает за приём сообщений, обработку json и вызов создания заказа.
type Worker struct {
	id           int
	messagesChan <-chan []byte
	processor    *processor.Processor
}

// NewWorker создает нового воркера.
// @summary Создание воркера
// @param id query int true "ID воркера"
// @param messagesChan body []byte true "Канал сообщений Кафки"
// @return Worker
func NewWorker(id int, messagesChan <-chan []byte, processor *processor.Processor) *Worker {
	return &Worker{id: id, messagesChan: messagesChan, processor: processor}
}

// Start запускает цикл обработки сообщений из канала.
// @summary Запуск воркера
// @description запускает обработчик сообщений из Кафки
func (w *Worker) Start() {
	logrus.Infof("Worker %d started", w.id)
	for message := range w.messagesChan {
		if message == nil {
			continue
		}
		if err := w.processor.ProcessMessage(message); err != nil {
			logrus.Errorf("Worker %d failed to process message: %v", w.id, err)
		}
	}
	logrus.Infof("Worker %d stopped", w.id)
}
