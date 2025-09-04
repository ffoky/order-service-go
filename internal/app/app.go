package app

import (
	"WBTECH_L0/internal/usecases/generator"
	"os"
	"os/signal"
	"syscall"

	appConfig "WBTECH_L0/config"
	pkgHttp "WBTECH_L0/internal/controller/http"
	"WBTECH_L0/internal/repository/cache"
	"WBTECH_L0/internal/repository/kafka"
	"WBTECH_L0/internal/repository/kafka/handler"
	"WBTECH_L0/internal/repository/postgres"
	"WBTECH_L0/internal/repository/postgres/delivery"
	"WBTECH_L0/internal/repository/postgres/items"
	"WBTECH_L0/internal/repository/postgres/order"
	"WBTECH_L0/internal/repository/postgres/payment"
	sender "WBTECH_L0/internal/usecases/sender"
	"WBTECH_L0/internal/usecases/service"
	pkgPostgres "WBTECH_L0/pkg/postgres"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	ordersNumToGenerate = 10
	consumerNumber      = 1
)

func Run(cfg *appConfig.AppConfig) {
	ctx := context.Background()

	pgCfg := &pkgPostgres.Config{
		Host:     cfg.DatabaseConfig.Host,
		Port:     cfg.DatabaseConfig.Port,
		User:     cfg.DatabaseConfig.User,
		Password: cfg.DatabaseConfig.Password,
		DBName:   cfg.DatabaseConfig.DBName,
		SSLMode:  cfg.DatabaseConfig.SSLMode,
	}

	pool, err := pkgPostgres.NewConnection(ctx, pgCfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	logrus.Info("Database connection established successfully")

	tm := postgres.NewTransactionManager(pool)
	deliveryRepo := delivery.NewRepository(pool)
	paymentRepo := payment.NewRepository(pool)
	itemsRepo := items.NewRepository(pool)
	orderRepo := order.NewRepository(
		pool,
		tm,
		*deliveryRepo,
		*paymentRepo,
		*itemsRepo,
	)

	orderCache := cache.NewCache(cfg.CacheTTL)
	logrus.Info("Cache created successfully")

	orderUseCase := service.NewOrderService(orderRepo, orderCache)

	if err := orderUseCase.InitializeCache(ctx, cfg.CacheTTL); err != nil {
		logrus.Errorf("Failed to initialize cache: %v", err)
	} else {
		stats := orderUseCase.GetCacheStats()
		logrus.Infof("Cache initialized successfully with %v orders", stats["cache_size"])
	}

	kafkaHandler := handler.NewHandler(orderUseCase)

	producer, err := kafka.NewProducer(cfg.KafkaConfig.Brokers)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	newGenerator := generator.NewGenerator()
	newSender := sender.NewSender(producer)
	orders := newGenerator.GenerateOrders(ordersNumToGenerate)

	if err := newSender.SendOrders(orders, cfg.KafkaConfig.Topic); err != nil {
		logrus.Errorf("Failed to send orders to Kafka: %v", err)
	} else {
		logrus.Info("Test orders sent to Kafka successfully")
	}

	consumer, err := kafka.NewConsumer(
		kafkaHandler,
		cfg.KafkaConfig.Brokers,
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.GroupID,
		consumerNumber,
	)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go consumer.Start()

	r := chi.NewRouter()
	orderHandlers := pkgHttp.NewOrderHandler(orderUseCase)
	orderHandlers.WithOrderHandlers(r)

	fs := http.FileServer(http.Dir("./static"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	go func() {
		logrus.Infof("HTTP server starting on %s", cfg.HTTPConfig.Address)
		logrus.Infof("Open http://%s in your browser to see the demo page", "localhost:8081")

		if err := http.ListenAndServe(cfg.HTTPConfig.Address, r); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("HTTP server failed: %v", err)
		}
	}()

	stats := orderUseCase.GetCacheStats()
	logrus.Infof("Application started successfully. Cache stats: %+v", stats)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logrus.Infof("app - Run - signal: %s", s.String())
	}

	logrus.Println("Shutting down gracefully...")

	if err := consumer.Stop(); err != nil {
		logrus.Errorf("Error stopping consumer: %v", err)
	}

	logrus.Info("Application shutdown completed")
}
