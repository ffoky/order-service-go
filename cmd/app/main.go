// @title           WBTECH L0 API
// @version         1.0
// @description     API для получения данных о заказе. Сервис обработки заказов с использованием Kafka, PostgreSQL и кэширования.

// @host      localhost:8081
// @BasePath  /

// @tag.name orders
// @tag.description Получение данных о заказе
package main

import (
	appConfig "WBTECH_L0/config"
	_ "WBTECH_L0/docs"
	"WBTECH_L0/internal/app"
	"log"
)

func main() {
	cfg, err := appConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
