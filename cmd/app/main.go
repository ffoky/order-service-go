package main

import (
	http3 "WBTECH_L0/internal/controller/http"
	"WBTECH_L0/internal/repository/handler"
	"WBTECH_L0/internal/repository/kafka"
	"WBTECH_L0/internal/repository/postgres"
	"WBTECH_L0/internal/repository/postgres/delivery"
	"WBTECH_L0/internal/repository/postgres/items"
	"WBTECH_L0/internal/repository/postgres/order"
	"WBTECH_L0/internal/repository/postgres/payment"
	sender2 "WBTECH_L0/internal/usecases/sender"
	"WBTECH_L0/internal/usecases/service"
	http2 "WBTECH_L0/pkg/http"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	PostgresUser     = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresPort     = os.Getenv("POSTGRES_PORT")
	PostgresDB       = os.Getenv("POSTGRES_DB")
	PostgresHost     = os.Getenv("POSTGRES_HOST")
	kafkaHost        = os.Getenv("KAFKA_HEALTHCHECK_HOST")
	kafkaPort        = os.Getenv("KAFKA_PORT")
	topic            = "orders"
)

func newPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		PostgresUser, PostgresPassword, PostgresHost, PostgresPort, PostgresDB,
	)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool connect %w", err)
	}
	return pool, nil
}

func main() {
	ctx := context.Background()
	pool, err := newPostgresConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	tm := postgres.NewTransactionManager(pool)

	deliveryRepository := delivery.NewRepository(pool)
	paymentRepository := payment.NewRepository(pool)
	itemsRepository := items.NewRepository(pool)
	orderRepository := order.NewRepository(
		pool,
		tm,
		*deliveryRepository,
		*paymentRepository,
		*itemsRepository,
	)

	orderUseCase := service.NewOrder(orderRepository)

	h := handler.NewHandler(orderUseCase)
	address := []string{"kafka:29092"}

	producer, err := kafka.NewProducer(address)
	if err != nil {
		logrus.Fatal(err)
	}
	defer producer.Close()

	generator := sender2.NewGenerator()
	sender := sender2.NewSender(producer)

	orders := generator.GenerateOrders(10)
	if err := sender.SendOrders(orders, topic); err != nil {
		logrus.Fatal(err)
	}

	consumer, err := kafka.NewConsumer(h, address, "orders", "my-consumer-group", 1)
	if err != nil {
		log.Fatal(err)
	}
	go consumer.Start()

	r := chi.NewRouter()
	orderHandlers := http3.NewOrderHandler(orderUseCase)
	orderHandlers.WithOrderHandlers(r)

	go func() {
		logrus.Infof("HTTP server started on :8080")
		if err := http2.CreateAndRunServer(r, ":8080"); err != nil {
			logrus.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logrus.Println("Shutting down...")
	consumer.Stop()
}
