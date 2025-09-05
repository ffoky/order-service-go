package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	HTTPConfig     HTTPConfig
	KafkaConfig    KafkaConfig
	DatabaseConfig DatabaseConfig
	CacheTTL       time.Duration
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type HTTPConfig struct {
	Address string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &AppConfig{}

	var err error

	cfg.HTTPConfig, err = loadHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load HTTP config: %w", err)
	}

	cfg.KafkaConfig, err = loadKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Kafka config: %w", err)
	}

	cfg.DatabaseConfig, err = loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Database config: %w", err)
	}

	cacheTTLStr, err := getEnv("CACHE_TTL")
	if err != nil {
		return nil, fmt.Errorf("CACHE_TTL is required: %w", err)
	}

	ttl, err := time.ParseDuration(cacheTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL format: %w", err)
	}
	cfg.CacheTTL = ttl

	return cfg, nil
}

func loadHTTPConfig() (HTTPConfig, error) {
	addr, err := getEnv("SERVER_ADDR")
	if err != nil {
		return HTTPConfig{}, err
	}

	port, err := getEnv("SERVER_PORT")
	if err != nil {
		return HTTPConfig{}, err
	}

	return HTTPConfig{
		Address: addr + ":" + port,
	}, nil
}

func loadKafkaConfig() (KafkaConfig, error) {
	var brokers []string

	if brokersEnv := os.Getenv("KAFKA_BROKERS"); brokersEnv != "" {
		brokers = strings.Split(brokersEnv, ",")
		for i, b := range brokers {
			brokers[i] = strings.TrimSpace(b)
		}
	} else {
		host, err := getEnv("KAFKA_HOST")
		if err != nil {
			return KafkaConfig{}, err
		}

		port, err := getEnv("KAFKA_PORT")
		if err != nil {
			return KafkaConfig{}, err
		}

		brokers = []string{host + ":" + port}
	}

	topic, err := getEnv("KAFKA_TOPIC")
	if err != nil {
		return KafkaConfig{}, err
	}

	groupID, err := getEnv("KAFKA_GROUP_ID")
	if err != nil {
		return KafkaConfig{}, err
	}

	return KafkaConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	host, err := getEnv("POSTGRES_HOST")
	if err != nil {
		return DatabaseConfig{}, err
	}

	portStr, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return DatabaseConfig{}, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	user, err := getEnv("POSTGRES_USER")
	if err != nil {
		return DatabaseConfig{}, err
	}

	password, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return DatabaseConfig{}, err
	}

	dbname, err := getEnv("POSTGRES_DB")
	if err != nil {
		return DatabaseConfig{}, err
	}

	sslmode, err := getEnv("POSTGRES_SSLMODE")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is required but not set", key)
	}
	return value, nil
}
