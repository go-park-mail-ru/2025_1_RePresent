package main

import (
	"log"
	"os"
	"os/signal"
	"retarget/configs"
	authApp "retarget/internal/auth-service/app"
	"syscall"
)

func main() {
	cfg, err := configs.LoadConfigs("configs/database.yml", "configs/mail.yml", "configs/auth-redis.yml")
	if err != nil {
		log.Fatal(err)
	}
	go authApp.Run(cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
