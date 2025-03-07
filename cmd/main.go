package main

import (
	"fmt"
	"log"

	"RE/configs"
	"RE/internal/app"
)

func main() {
	cfg, err := configs.LoadConfig("configs/database.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	app.Run(cfg)
}
