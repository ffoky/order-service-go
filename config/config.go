package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	DefaultConfigPath = "./config/config.yml"
	DefaultKafkaPort  = "9092"
	DefaultServerPort = "8081"
)

var (
	envPostgresHost     = os.Getenv("POSTGRES_HOST")
	envPostgresPort     = os.Getenv("POSTGRES_PORT")
	envPostgresUser     = os.Getenv("POSTGRES_USER")
	envPostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	envPostgresDB       = os.Getenv("POSTGRES_DB")
	envServerAddr       = os.Getenv("SERVER_ADDR")
	envServerPort       = os.Getenv("SERVER_PORT")
	envKafkaHost        = os.Getenv("KAFKA_HEALTHCHECK_HOST")
	envKafkaPort        = os.Getenv("KAFKA_PORT")
	envKafkaTopic       = os.Getenv("KAFKA_TOPIC")
)

type AppConfig struct {
	HTTPConfig     `yaml:"http"`
	KafkaConfig    `yaml:"kafka"`
	DatabaseConfig `yaml:"database"`
	CacheTTL       time.Duration `yaml:"cache_ttl"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var cfg AppConfig
	if err := cleanenv.ReadConfig(DefaultConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	overrideFromEnv(&cfg)

	return &cfg, nil
}

func overrideFromEnv(cfg *AppConfig) {
	if envPostgresHost != "" {
		cfg.DatabaseConfig.Host = envPostgresHost
	}
	if envPostgresPort != "" {
		cfg.DatabaseConfig.Port = parseInt(envPostgresPort, cfg.DatabaseConfig.Port)
	}
	if envPostgresUser != "" {
		cfg.DatabaseConfig.User = envPostgresUser
	}
	if envPostgresPassword != "" {
		cfg.DatabaseConfig.Password = envPostgresPassword
	}
	if envPostgresDB != "" {
		cfg.DatabaseConfig.DBName = envPostgresDB
	}

	if envServerAddr != "" {
		port := envServerPort
		if port == "" {
			port = DefaultServerPort
		}
		cfg.HTTPConfig.Address = envServerAddr + ":" + port
	}

	if envKafkaHost != "" {
		kafkaPort := envKafkaPort
		if kafkaPort == "" {
			kafkaPort = DefaultKafkaPort
		}
		kafkaAddress := envKafkaHost + ":" + kafkaPort

		if len(cfg.KafkaConfig.Brokers) > 0 {
			cfg.KafkaConfig.Brokers[0] = kafkaAddress
		} else {
			cfg.KafkaConfig.Brokers = []string{kafkaAddress}
		}
	}
	if envKafkaTopic != "" {
		cfg.KafkaConfig.Topic = envKafkaTopic
	}
}

func parseInt(s string, defaultValue int) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}
