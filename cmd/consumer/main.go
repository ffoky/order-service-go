package main

import (
	"WBTECH_L0/internal/repository/handler"
	"WBTECH_L0/internal/repository/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const (
	topic         = "order"
	consumerGroup = "my-consumer-group"
)

var address = []string{"localhost:9092"}

func main() {
	h := handler.NewHandler()
	c, err := kafka.NewConsumer(h, address, topic, consumerGroup, 1)
	if err != nil {
		logrus.Fatal(err)
	}
	go c.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logrus.Fatal(c.Stop())
}
