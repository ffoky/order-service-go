package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type AppFlags struct {
	LogFormat string
}

func ParseFlags() AppFlags {
	logFormat := flag.String("logf", "", "Logs format")
	flag.Parse()
	return AppFlags{
		LogFormat: *logFormat,
	}
}

func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	// Устанавливаем значения из переменных окружения после загрузки конфига
	if appCfg, ok := cfg.(*AppConfig); ok {
		if dbUser := os.Getenv("POSTGRES_USER"); dbUser != "" {
			appCfg.DatabaseConfig.User = dbUser
		}
		if dbPassword := os.Getenv("POSTGRES_PASSWORD"); dbPassword != "" {
			appCfg.DatabaseConfig.Password = dbPassword
		}
		if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
			appCfg.DatabaseConfig.DBName = dbName
		}
	}
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

type AppConfig struct {
	HTTPConfig     `yaml:"http"`
	KafkaConfig    `yaml:"kafka"`
	DatabaseConfig `yaml:"database"`
}
