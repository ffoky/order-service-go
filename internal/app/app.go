package app

import (
	"WBTECH_L0/internal/usecases/kafka/generator"
	"WBTECH_L0/internal/usecases/service/processor"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appConfig "WBTECH_L0/config"
	pkgHttp "WBTECH_L0/internal/controller/http"
	"WBTECH_L0/internal/repository/cache"
	"WBTECH_L0/internal/repository/kafka"
	"WBTECH_L0/internal/repository/postgres"
	"WBTECH_L0/internal/repository/postgres/delivery"
	"WBTECH_L0/internal/repository/postgres/items"
	"WBTECH_L0/internal/repository/postgres/order"
	"WBTECH_L0/internal/repository/postgres/payment"
	sender "WBTECH_L0/internal/usecases/kafka/sender"
	"WBTECH_L0/internal/usecases/service"
	"WBTECH_L0/internal/usecases/service/worker"
	pkgHttpServer "WBTECH_L0/pkg/http"
	pkgPostgres "WBTECH_L0/pkg/postgres"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	ordersNumToGenerate = 10
	workersCount        = 5
	messageBufferSize   = 100
	shutdownTimeout     = 5 * time.Second
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

	orderProcessor := processor.NewProcessor(orderUseCase)

	messagesChan := make(chan []byte, messageBufferSize)

	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w := worker.NewWorker(workerID, messagesChan, orderProcessor)
			w.Start()
		}(i + 1)
	}
	logrus.Infof("Started %d workers", workersCount)

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
		cfg.KafkaConfig.Brokers,
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.GroupID,
	)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go func() {
		logrus.Info("Starting Kafka message reader...")
		consumer.StartReading(messagesChan)
	}()

	r := chi.NewRouter()
	orderHandlers := pkgHttp.NewOrderHandler(orderUseCase)
	staticHandlers := pkgHttp.NewStaticHandler()

	orderHandlers.WithOrderHandlers(r)
	staticHandlers.WithStaticHandlers(r)

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	httpServer, err := pkgHttpServer.CreateServerWithShutdown(r, cfg.HTTPConfig.Address)
	if err != nil {
		logrus.Fatalf("Failed to create HTTP server: %v", err)
	}

	logrus.Infof("HTTP server starting on %s", cfg.HTTPConfig.Address)
	logrus.Infof("Open http://%s in your browser to see the demo page", "localhost:8081")

	stats := orderUseCase.GetCacheStats()
	logrus.Infof("Application started successfully. Cache stats: %+v", stats)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	logrus.Infof("app - Run - signal: %s", s.String())

	logrus.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := pkgHttpServer.ShutdownServer(shutdownCtx, httpServer); err != nil {
		logrus.Errorf("Error shutting down HTTP server: %v", err)
	} else {
		logrus.Info("HTTP server stopped gracefully")
	}

	if err := consumer.Stop(); err != nil {
		logrus.Errorf("Error stopping consumer: %v", err)
	} else {
		logrus.Info("Kafka consumer stopped gracefully")
	}

	close(messagesChan)
	wg.Wait()
	logrus.Info("All workers stopped")

	logrus.Info("Application shutdown completed")
}
