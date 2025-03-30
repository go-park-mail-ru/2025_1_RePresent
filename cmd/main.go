package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"retarget/configs"

	// profileApp "retarget/internal/profile-service/app"
	advApp "retarget/internal/adv-service/app"
	mailApp "retarget/internal/mail-service/app"
)

func main() {
	cfg, err := configs.LoadConfigs("configs/database-dev.yml", "configs/mail.yml")
	if err != nil {
		log.Fatal(err)
	}

	// go authApp.Run(cfg)
	go mailApp.Run(cfg)
	go advApp.Run(cfg)
	// go profileApp.Run(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
