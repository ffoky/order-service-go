package main

import (
	appConfig "WBTECH_L0/config"
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
