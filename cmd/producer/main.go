package main

import (
	"WBTECH_L0/internal/repository/kafka"
	"WBTECH_L0/internal/usecases/order"
	"github.com/sirupsen/logrus"
)

const (
	topic = "order"
)

var address = []string{"localhost:9092"}

func main() {
	// Создание зависимостей
	producer, err := kafka.NewProducer(address)
	if err != nil {
		logrus.Fatal(err)
	}
	defer producer.Close()

	generator := order.NewGenerator()
	sender := order.NewSender(producer)

	// Основная логика
	orders := generator.GenerateOrders(20)
	if err := sender.SendOrders(orders, topic); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Программа завершена.")
}
