package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"retarget/configs"
	"retarget/internal/app"
	mailApp "retarget/internal/mail-service/app"
)

func main() {
	cfg, err := configs.LoadConfigs("configs/database-dev.yml", "configs/mail.yml")
	if err != nil {
		log.Fatal(err)
	}

	go app.Run(cfg)
	go mailApp.Run(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
